package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestEmailHandler(t *testing.T) {
	handler, _ := setupTestServer()
	ts := httptest.NewServer(http.HandlerFunc(handler.HandleSendFileByEmail))
	defer ts.Close()

	payload := "emails=yerezhepabdulla@gmail.com,kopoilamadim@gmail.com"

	tempFile, err := os.CreateTemp("", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString("Sample file content")
	tempFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	payloadPart, err := writer.CreateFormFile("payload", "payload.txt")
	if err != nil {
		t.Fatal(err)
	}
	payloadReader := strings.NewReader(payload)
	io.Copy(payloadPart, payloadReader)

	filePart, err := writer.CreateFormFile("file", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(filePart, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	resp, err := http.Post(ts.URL+"/api/mail/file", writer.FormDataContentType(), body)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
