package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"test/config"
	"test/internal/service"
	"test/internal/service/archive"
	"test/internal/service/email"
	"testing"
)
func TestGetArchiveInfo(t *testing.T) {
	handler,_ := setupTestServer()
	ts := httptest.NewServer(http.HandlerFunc(handler.GetArchiveInfo))
	

	defer ts.Close()

	mockZIPFile, _ := createMockZIPFile()
	mockFileHeader := createFileHeader("sample.zip", mockZIPFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", mockFileHeader.Filename)
	part.Write(mockZIPFile)
	writer.Close()

	resp, err := http.Post(ts.URL+"/api/archive/information", writer.FormDataContentType(), body)
	if err != nil {
		t.Errorf("Failed to send a POST request: %s", err.Error())
		return
	}
	fmt.Println(resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

}

func TestCreateArchiveHandler(t *testing.T) {
	handler,_ := setupTestServer()
	ts := httptest.NewServer(http.HandlerFunc(handler.HandleCreateArchive))
	

	defer ts.Close()

	mockFileHeaders := []*multipart.FileHeader{
		createFileHeader("file1.txt", []byte("This is the content of file 1")),
		createFileHeader("file2.jpg", []byte("This is the content of file 2")),
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, header := range mockFileHeaders {
		part, _ := writer.CreateFormFile("files[]", header.Filename)
		part.Write([]byte("Sample file content"))
	}
	writer.Close()

	resp, err := http.Post(ts.URL+"/api/archive/files", writer.FormDataContentType(), body)
	if err != nil {
		t.Errorf("Failed to send a POST request: %s", err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

}

func setupTestServer() (*Handler, error) {
	// Set up a test HTTP server.
	config := config.MustLoad()
	//int archive service 
	archive := archive.NewArchiveService()
	//init email service 
	email := email.NewEmailService()
	//init service
	service := service.NewService(archive, email, config)
	//init handler
	handler := NewHandler(service)

	return handler, nil
}
func createMockZIPFile() ([]byte, error) {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	fileContents := map[string]string{
		"file1.txt": "This is the content of file 1",
		"file2.txt": "This is the content of file 2",
	}

	for fileName, content := range fileContents {
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, err
		}

		_, err = fileWriter.Write([]byte(content))
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return zipBuffer.Bytes(), nil
}

func createFileHeader(fileName string, content []byte) *multipart.FileHeader {
	fileHeader := &multipart.FileHeader{
		Filename: fileName,
		Size:     int64(len(content)),
		Header:   make(map[string][]string),
	}

	fileHeader.Header.Set("Content-Type", "application/octet-stream")
	return fileHeader
}
