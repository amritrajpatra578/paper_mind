package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/upload", handleUpload).Methods("POST")
	r.HandleFunc("/api/ask", handleAsk).Methods("POST")

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // setting the size limit to 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	filePath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	textByPage, err := extractTextByPage(filePath)
	if err != nil {
		http.Error(w, "Failed to parse PDF", http.StatusInternalServerError)
		return
	}

	for i, pageText := range textByPage {
		log.Printf("Page %d:\n%s\n", i+1, pageText)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded and parsed successfully: %s\n", handler.Filename)
}

func extractTextByPage(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	n := r.NumPage()
	var pages []string
	for i := 1; i <= n; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			pages = append(pages, "")
		} else {
			pages = append(pages, strings.TrimSpace(content))
		}
	}
	return pages, nil
}

func handleAsk(w http.ResponseWriter, r *http.Request) {
	type AskRequest struct {
		Question string `json:"question"`
		PdfPath  string `json:"pdfPath"`
	}
	type AskResponse struct {
		Answer    string `json:"answer"`
		Citations []int  `json:"citations"`
	}

	var req AskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Question == "" || req.PdfPath == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	pages, err := extractTextByPage(req.PdfPath)
	if err != nil {
		http.Error(w, "Unable to read PDF content", http.StatusInternalServerError)
		return
	}

	// Naive search for relevant pages (to be replaced with vector search)
	var relevantPages []int
	var relevantText bytes.Buffer
	questionLower := strings.ToLower(req.Question)
	for i, text := range pages {
		if strings.Contains(strings.ToLower(text), questionLower) {
			relevantPages = append(relevantPages, i+1)
			relevantText.WriteString(fmt.Sprintf("[Page %d]\n%s\n\n", i+1, text))
			if len(relevantPages) >= 3 {
				break
			}
		}
	}

	answer := "I'm an early version of the assistant. Here's what I found:\n" + relevantText.String()

	resp := AskResponse{
		Answer:    answer,
		Citations: relevantPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
