package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/delaemon/go-uuidv4"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var FAKEUUID, _ = uuidv4.Generate() //"60cd5e90-d79f-465e-93a9-134f39d6015f"

var fakeShard = setShard(FAKEUUID)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/upload", UploadJpeg).Methods("POST")
	return r
}

func SetupTest(dirname string, ct string) {
	getProcessingPoolSize = func() int {
		return 1
	}

	setContentType = func(m textproto.MIMEHeader) []string {
		r := make([]string, 1)
		r[0] = ct
		return r
	}

	getBasePath = func() string {
		return dirname
	}

	setUUIDv4 = func() (string, error) {
		return FAKEUUID, nil
	}

	setFilename = func(fname string) string {
		return FAKEUUID + JPEGEXT
	}
}

func loadImageFile(path string) (*bytes.Buffer, *multipart.Writer) {
	file, _ := os.Open(path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(path))
	io.Copy(part, file)
	writer.Close()

	return body, writer
}

func mockReqRes(body *bytes.Buffer, writer *multipart.Writer) (*http.Request, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	return req, res
}

func TestUploadJpegOK(t *testing.T) {

	dirname := "./test/fileUpload"

	SetupTest(dirname, "image/jpeg")

	assert := assert.New(t)
	path := "./test/small.jpg"
	body, writer := loadImageFile(path)

	req, res := mockReqRes(body, writer)
	Router().ServeHTTP(res, req)

	assert.Equal(res.Code, 200)

	uploadedImgPath := fmt.Sprintf("%s%s%s%s%s", dirname, SLASH, fakeShard, SLASH, FAKEUUID+JPEGEXT)
	_, err := os.Stat(uploadedImgPath)

	assert.NoError(err)

	orig, _ := os.Open(path)
	defer orig.Close()
	origBytes := []byte{}
	orig.Read(origBytes)

	uploaded, _ := os.Open(uploadedImgPath)
	defer uploaded.Close()
	uploadedBytes := []byte{}
	uploaded.Read(uploadedBytes)

	assert.Equal(string(uploadedBytes), string(origBytes))

	os.RemoveAll(dirname + SLASH + fakeShard)
}

func TestUploadJpegTooLarge(t *testing.T) {

	dirname := "./test/fileUpload"

	SetupTest(dirname, "image/jpeg")

	assert := assert.New(t)
	path := "./test/large.jpg"
	body, writer := loadImageFile(path)

	req, res := mockReqRes(body, writer)
	Router().ServeHTTP(res, req)

	assert.Equal(res.Code, 400)
}

func TestUploadJpegWrongType(t *testing.T) {

	dirname := "./test/fileUpload"

	SetupTest(dirname, "image/png")

	assert := assert.New(t)
	path := "./test/wrongType.png"
	body, writer := loadImageFile(path)

	req, res := mockReqRes(body, writer)
	Router().ServeHTTP(res, req)

	assert.Equal(res.Code, 400)
}
