package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"os"

	"github.com/delaemon/go-uuidv4"
	"github.com/gorilla/mux"
)

// JPEGEXT extension const
const JPEGEXT = ".jpeg"

// SLASH const
const SLASH = "/"

var getBasePath = func() string {
	return os.Getenv("BASE_PATH")
}

var setContentType = func(mimeHeader textproto.MIMEHeader) []string {
	return mimeHeader["Content-Type"]
}

var setUUIDv4 = func() (string, error) {
	return uuidv4.Generate()
}

var setShard = func(uuid string) string {
	return uuid[len(uuid)-2:]
}

var setFilename = func(fname string) string {
	return fname
}

// respondWithError
func respondWithError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, code int, key string, value string) {
	response, _ := json.Marshal(map[string]string{key: value})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func saveImage(fullPath string) {
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

// UploadJpeg recives JPEG images only with a maximum size of 8192Kb
func UploadJpeg(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(8 * 1024); err != nil {
		log.Fatal(err)
	}

	basePath := getBasePath()
	if basePath == "" {
		log.Fatal("BASE_PATH env variable not set!")
	}

	uuid, uErr := setUUIDv4()
	shard := setShard(uuid)

	if uErr != nil {
		log.Fatal(uErr)
	}

	file, handler, err := r.FormFile("file")
	handler.Filename = setFilename(uuid + JPEGEXT)
	contentType := setContentType(handler.Header)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if handler.Size > 8192 || (len(contentType) > 0 && contentType[0] != "image/jpeg") {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, "image_id", uuid)

	path := fmt.Sprintf("%s%s%s%s", basePath, SLASH, shard, SLASH)
	if err := os.MkdirAll(path, os.ModePerm); err == nil {
		saveImage(path + handler.Filename)
		return
	}
}

func main() {
	fmt.Println("UPLOAD THUMBNAILS SERVICE STARTED")
	fmt.Println("listening on port: 8080")

	router := mux.NewRouter()
	router.HandleFunc("/upload", UploadJpeg).Name("/upload").Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
