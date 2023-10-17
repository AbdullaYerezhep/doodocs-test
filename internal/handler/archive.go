package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
)

func (h *Handler) GetArchiveInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return
	}
	fileheader := r.MultipartForm.File["file"]
	archive, err := h.service.Archive.GetArchiveInfo(fileheader[0])
	if err != nil {
		http.Error(w, "Failed to get archive info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(archive)
}


func (h *Handler) HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return
	}

	fileHeaders, exists := r.MultipartForm.File["files[]"]
	if !exists {
		http.Error(w, "No files were uploaded", http.StatusBadRequest)
		return
	}

	var fileHeadersNew []multipart.FileHeader // Use a string slice to store filenames
	for _, fileHeader := range fileHeaders {
		fileHeadersNew = append(fileHeadersNew, *fileHeader)
		if err != nil {
			http.Error(w, "Failed to get archive info", http.StatusInternalServerError)
			return
		}
	}

	zipData, err := h.service.Archive.CreateArchive(fileHeadersNew)
	if err != nil {
		http.Error(w, "Failed to create archive", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")

	_, err = w.Write(zipData)
	if err != nil {
		http.Error(w, "Failed to write archive to response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}