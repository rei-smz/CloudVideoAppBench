package function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ReqArgs struct {
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	AudioCodec string `json:"acodec"`
	VideoCodec string `json:"vcodec"`
}

type ReqData struct {
	ReqType string            `json:"type" bson:"type"`
	Path    string            `json:"path" bson:"path"`
	Object  string            `json:"object" bson:"object"`
	Args    map[string]string `json: "args" bson:"args"`
}

type ReqInfo struct {
	Status    string  `bson:"status" json:"status"`
	StartTime int64   `bson:"start_time" json:"start_time"`
	EndTime   int64   `bson:"end_time" json:"end_time"`
	Retry     int     `bson:"retry" json:"retry"`
	Data      ReqData `bson:"req_data" json: "data"`
}

var bucketName = os.Getenv("BUCKET_NAME")

var infoDBName = os.Getenv("MONGO_DB")

var infoCollName = os.Getenv("MONGO_COLLECTION")

var minioClient = initMinio()

var mongoClient = initMongo()

func getSecret(key string) string {
	path := filepath.Join("/var/openfaas/secrets", key)
	data, err := os.ReadFile(path)

	if err != nil {
		return ""
	}
	return string(data)
}

func initMinio() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := getSecret("minio-access-key")
	secretAccessKey := getSecret("minio-secret-key")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return client
}

func initMongo() *mongo.Client {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		log.Fatal("Set your 'MONGO_URI' environment variable. ")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return client
}

func updateReqInfo(coll *mongo.Collection, infoID primitive.ObjectID, update bson.D) error {
	_, err := coll.UpdateByID(context.TODO(), infoID, update)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func downloadVideo(client *minio.Client, objPath string, localPath string) error {
	err := client.FGetObject(context.Background(),
		bucketName,
		objPath,
		localPath,
		minio.GetObjectOptions{})
	return err
}

func uploadVideo(client *minio.Client, objPath string, localPath string) error {
	_, err := client.FPutObject(context.Background(), bucketName, objPath, localPath, minio.PutObjectOptions{})
	return err
}

func runFFMpeg(args map[string]string, tmpPath string, tmpIn string) (string, error) {
	outFormat, ok := args["format"]
	if !ok {
		outFormat = "mp4"
	}
	localOutTmp := "out." + outFormat
	resolution, ok := args["resolution"]
	if !ok {
		resolution = "1280x720"
	}
	aCodec, ok := args["acodec"]
	if !ok {
		aCodec = "copy"
	}
	vCodec := args["vcodec"]

	cmdArgs := []string{"-hide_banner", "-loglevel", "warning", "-y", "-i", tmpIn}
	if resolution != "no" {
		cmdArgs = append(cmdArgs, "-s", resolution)
	}
	if aCodec != "" {
		cmdArgs = append(cmdArgs, "-acodec", aCodec)
	}
	if vCodec != "" {
		cmdArgs = append(cmdArgs, "-vcodec", vCodec)
	}
	cmdArgs = append(cmdArgs, filepath.Join(tmpPath, localOutTmp))
	cmd := exec.Command("ffmpeg", cmdArgs...)
	_, err := cmd.CombinedOutput()
	return localOutTmp, err
}

func handleRequest(w http.ResponseWriter, reqData ReqData, coll *mongo.Collection, infoID primitive.ObjectID) {
	objPath := reqData.Path
	objName := reqData.Object
	reqArgs := reqData.Args
	if objPath == "" || objName == "" {
		http.Error(w, "Path and object name are required", http.StatusBadRequest)
		return
	}
	tmpDir, err := os.MkdirTemp("", filepath.Base(objPath)+"-*")
	if err != nil {
		http.Error(w, "Error creating temp directory", http.StatusInternalServerError)
		return
	}

	id := infoID
	if id == primitive.NilObjectID {
		newReqInfo := ReqInfo{
			Status:    "running",
			StartTime: time.Now().UnixMilli(),
			EndTime:   -1,
			Retry:     0,
			Data:      reqData,
		}
		result, err := coll.InsertOne(context.TODO(), newReqInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id = result.InsertedID.(primitive.ObjectID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"key\": \"%v\"}", id.Hex())))
	}

	go func() {
		defer os.RemoveAll(tmpDir)

		localInPath := filepath.Join(tmpDir, objName)
		err = downloadVideo(minioClient, filepath.Join(objPath, objName), localInPath)
		if err != nil {
			updateReqInfo(coll, id, bson.D{{"$set", bson.D{{"status", "error"}}}})
			return
		}

		localOut, err := runFFMpeg(reqArgs, tmpDir, localInPath)
		if err != nil {
			updateReqInfo(coll, id, bson.D{{"$set", bson.D{{"status", "error"}}}})
			return
		}

		localOutPath := filepath.Join(tmpDir, localOut)
		err = uploadVideo(minioClient, filepath.Join(objPath, localOut), localOutPath)
		if err != nil {
			updateReqInfo(coll, id, bson.D{{"$set", bson.D{{"status", "error"}}}})
			return
		}

		updateReqInfo(coll, id, bson.D{{"$set", bson.D{{"status", "success"},
			{"end_time", time.Now().UnixMilli()}}}})
	}()
}

func handleQuery(w http.ResponseWriter, reqData ReqData, coll *mongo.Collection) {
	key := reqData.Args["key"]
	if key == "" {
		http.Error(w, "Query key is required", http.StatusBadRequest)
		return
	}

	infoID, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		http.Error(w, "Invalid object id", http.StatusBadRequest)
		return
	}
	reqInfo := &ReqInfo{}
	err = coll.FindOne(context.TODO(), bson.M{"_id": infoID}).Decode(reqInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find req info: %v", err), http.StatusInternalServerError)
		return
	}

	if reqInfo.Status == "running" && time.Since(time.UnixMilli(reqInfo.StartTime)) > time.Minute*10 {
		if reqInfo.Retry < 3 {
			updateReqInfo(coll, infoID, bson.D{{"$inc", bson.D{{"retry", 1}}}})
			handleRequest(w, reqInfo.Data, coll, infoID)
		} else {
			reqInfo.Status = "error"
			updateReqInfo(coll, infoID, bson.D{{"$set", bson.D{{"status", "error"}}}})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reqInfo)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if minioClient == nil || mongoClient == nil {
		http.Error(w, "Unable to connect to Minio or MongoDB", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	var reqData ReqData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqType := reqData.ReqType
	if reqType == "" {
		http.Error(w, "Type is required", http.StatusBadRequest)
		return
	}

	coll := mongoClient.Database(infoDBName).Collection(infoCollName)
	if reqType == "request" {
		handleRequest(w, reqData, coll, primitive.NilObjectID)
	} else if reqType == "query" {
		handleQuery(w, reqData, coll)
	} else {
		http.Error(w, "Invalid req type", http.StatusBadRequest)
	}
}
