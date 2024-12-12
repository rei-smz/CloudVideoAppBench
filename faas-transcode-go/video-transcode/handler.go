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
	ReqType string `json:"type"`
	Path    string `json:"path"`
	Object  string `json:"object"`
	Args    map[string]string
}

type ReqInfo struct {
	Status    string `bson:"status" json:"status"`
	StartTime int64  `bson:"start_time" json:"start_time"`
	EndTime   int64  `bson:"end_time" json:"end_time"`
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
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. ")
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

func handleRequest(tmpDir string, objPath string, objName string, reqArgs map[string]string,
	coll *mongo.Collection, infoID primitive.ObjectID) {
	defer os.RemoveAll(tmpDir)

	localInPath := filepath.Join(tmpDir, objName)
	err := downloadVideo(minioClient, filepath.Join(objPath, objName), localInPath)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//statusTable.Store(tmpDir, "error")
		updateReqInfo(coll, infoID, bson.D{{"$set", bson.D{{"status", "error"}}}})
		return
	}

	localOut, err := runFFMpeg(reqArgs, tmpDir, localInPath)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//statusTable.Store(tmpDir, "error")
		updateReqInfo(coll, infoID, bson.D{{"$set", bson.D{{"status", "error"}}}})
		return
	}

	localOutPath := filepath.Join(tmpDir, localOut)
	err = uploadVideo(minioClient, filepath.Join(objPath, localOut), localOutPath)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//statusTable.Store(tmpDir, "error")
		updateReqInfo(coll, infoID, bson.D{{"$set", bson.D{{"status", "error"}}}})
		return
	}

	//statusTable.Store(tmpDir, "completed")
	updateReqInfo(coll, infoID, bson.D{{"$set", bson.D{{"status", "success"},
		{"end_time", time.Now().UnixMilli()}}}})
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

	objPath := reqData.Path
	objName := reqData.Object
	reqType := reqData.ReqType
	if objPath == "" || objName == "" || reqType == "" {
		http.Error(w, "Path, object name and type are required", http.StatusBadRequest)
		return
	}

	coll := mongoClient.Database(infoDBName).Collection(infoCollName)
	if reqType == "request" {
		tmpDir, err := os.MkdirTemp("", filepath.Base(objPath)+"-*")
		if err != nil {
			http.Error(w, "Error creating temp directory", http.StatusInternalServerError)
			return
		}
		newReqInfo := ReqInfo{
			Status:    "running",
			StartTime: time.Now().UnixMilli(),
			EndTime:   -1,
		}
		result, err := coll.InsertOne(context.TODO(), newReqInfo)
		if err != nil {
			http.Error(w, "Error creating req info", http.StatusInternalServerError)
			return
		}
		//statusTable.Store(tmpDir, "running")
		go handleRequest(tmpDir, objPath, objName, reqData.Args, coll, result.InsertedID.(primitive.ObjectID))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.InsertedID.(primitive.ObjectID).Hex()))
	} else if reqType == "query" {
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
		var reqInfo ReqInfo
		err = coll.FindOne(context.TODO(), bson.M{"_id": infoID}).Decode(reqInfo)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find req info: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(reqInfo)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid req type", http.StatusBadRequest)
	}
}
