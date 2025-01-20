package routes

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadAssesment(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form.
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB.
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from the form.
	student_uuid := r.URL.Query().Get("uuid")
	cource_name := r.URL.Query().Get("cource_name")

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the file has a .pdf extension.
	if filepath.Ext(handler.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Create the destination file.

	parts := []string{"assesmentFiles", student_uuid, cource_name}
	result := strings.Join(parts, "/")

	savePath := filepath.Join(result, handler.Filename)
	err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	destFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the uploaded file to the destination file.
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	// Respond to the client.

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "fileuploaded", nil)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
}
