package handler

import (
	"net/http"
	"strings"
)
	
func (h *Handler) HandleSendFileByEmail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	emailsString := r.PostFormValue("emails")

    // Split the string into individual email addresses.
    emails := strings.Split(emailsString, ",")
	fileHeaders := r.MultipartForm.File["file"]
    if len(fileHeaders) == 0 {
        http.Error(w, "No file was uploaded", http.StatusBadRequest)
        return
    }
    file := fileHeaders[0]
    

	ok, err := h.service.Email.SendFileByEmail(emails, file, h.service.Config)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email sent successfully"))
	} else {
		w.Write([]byte(err.Error()))
	}
}