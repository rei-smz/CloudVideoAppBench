package function

import (
	"context"
	"encoding/json"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type ReqArgs struct {
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	AudioCodec string `json:"acodec"`
	VideoCodec string `json:"vcodec"`
}

type ReqData struct {
	Path   string `json:"path"`
	Object string `json:"object"`
	Args   map[string]string
}

var bucketName = os.Getenv("BUCKET_NAME")

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

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	var reqData ReqData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	minioClient := initMinio()
	objPath := reqData.Path
	objName := reqData.Object
	if objPath == "" || objName == "" {
		http.Error(w, "Path and object name are required", http.StatusBadRequest)
		return
	}

	tmpDir, err := os.MkdirTemp("", filepath.Base(objPath)+"-*")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		http.Error(w, "Error creating temp directory", http.StatusInternalServerError)
		return
	}

	localInPath := filepath.Join(tmpDir, objName)
	err = downloadVideo(minioClient, filepath.Join(objPath, objName), localInPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	localOut, err := runFFMpeg(reqData.Args, tmpDir, localInPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	localOutPath := filepath.Join(tmpDir, localOut)
	err = uploadVideo(minioClient, filepath.Join(objPath, localOut), localOutPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Body: success"))
}
