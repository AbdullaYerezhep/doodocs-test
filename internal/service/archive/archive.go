package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"test/internal/models"
)

type ArchiveService interface {
	GetArchiveInfo(fileHeader *multipart.FileHeader) (Archive, error)
	CreateArchive(file *[]multipart.FileHeader) (bool,error)
}

func NewArchiveService() *Archive {
	return &Archive{}
}

type Archive struct {}

func (as *Archive) GetArchiveInfo(fileheader *multipart.FileHeader) (models.Archive, error) {
	var files []models.File
	file := models.File {
		FilePath: "",
		FileSize: 0,
		MimeType: "",
	}
	
	archiveInfo := models.Archive{
		ArchiveName: fileheader.Filename,
		ArchiveSize: fileheader.Size,
		TotalSize: 0,
		TotalFile: 0,
		Files: files,
	}

	zipFile, err := fileheader.Open()
	if err != nil {
		return models.Archive{}, err
	}
	defer zipFile.Close()
	magicNumber := make([]byte, 4)
	_, err = io.ReadFull(zipFile, magicNumber)
	if err != nil || string(magicNumber) != "PK\x03\x04" {
		return models.Archive{}, fmt.Errorf("the provided file is not a valid ZIP archive")
	}
	reader, err := zip.NewReader(zipFile, fileheader.Size)


	for _, zipContent := range reader.File {
		archiveInfo.TotalSize += float64(zipContent.UncompressedSize64)
		archiveInfo.TotalFile++
		file.FilePath = zipContent.Name
		file.FileSize = float64(zipContent.UncompressedSize64)
		file.MimeType = mime.TypeByExtension(filepath.Ext(zipContent.Name))
		files = append(files, file)
	}
	archiveInfo.Files = files

	if err != nil {
		return models.Archive{}, err
	}
	return archiveInfo, nil
}	

func (as *Archive) CreateArchive(files []multipart.FileHeader) ([]byte, error) {
	allowedMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/xml": true,
		"image/jpeg": true,
		"image/png": true,
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	for _, file := range files {
		if !allowedMimeTypes[file.Header.Get("Content-Type")] {
			return nil, fmt.Errorf("file %s has an invalid MIME type", file.Filename)
		}
		uploadedFile, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer uploadedFile.Close()

		zipFile, err := zipWriter.Create(file.Filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(zipFile, uploadedFile)
		if err != nil {
			return nil, err
		}
	}

	zipData := zipBuffer.Bytes()

	return zipData, nil
}


