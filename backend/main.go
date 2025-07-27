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

	"github.com/joho/godotenv"
	"github.com/ledongthuc/pdf"
)

// Global variable to store path of last uploaded PDF
var lastUploadedPath string

// Struct for chat messages (used for Groq API)
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Struct for Groq request payload
type GroqRequest struct {
	Messages            []ChatMessage `json:"messages"`
	Model               string        `json:"model"`
	Temperature         float64       `json:"temperature"`
	TopP                float64       `json:"top_p"`
	Stream              bool          `json:"stream"`
	Stop                interface{}   `json:"stop"`
	MaxCompletionTokens int           `json:"max_completion_tokens"`
}

// Struct for client-side question request
type AskRequest struct {
	Question string `json:"question"`
	PdfPath  string `json:"pdfPath"`
}

// Struct for response sent back to frontend
type AskResponse struct {
	Answer    string `json:"answer"`
	Citations []int  `json:"citations"`
}

func main() {
	// Load environment variables (i.e, GROQ_API_KEY)
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file")
	}

	// Routes
	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/ask", handleAsk)

	// uploaded PDFs
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Server running at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleUpload handles PDF uploads and stores the file locally
func handleUpload(w http.ResponseWriter, r *http.Request) {
	const maxSize = 200 << 20 // 200MB max file size
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure upload directory exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Delete last uploaded file (non-user specific for simplicity)
	if lastUploadedPath != "" {
		_ = os.Remove(lastUploadedPath)
	}

	// Save new file
	savePath := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	lastUploadedPath = savePath // Save current file path

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message":  "Uploaded",
		"pdfPath":  "/" + savePath,
		"filename": handler.Filename,
	});err != nil {
		log.Println("failed to write response for upload method:",err)
	}
}

// handleAsk processes the user question and responds using Groq API
func handleAsk(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" || req.PdfPath == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	pages, err := extractText(strings.TrimPrefix(req.PdfPath, "/"))
	if err != nil {
		http.Error(w, "Failed to parse PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Simple keyword matching
	var hits []string
	var citations []int
	lowerQ := strings.ToLower(req.Question)
	for i, p := range pages {
		if strings.Contains(strings.ToLower(p), lowerQ) {
			hits = append(hits, fmt.Sprintf("[Page %d]\n%s", i+1, p))
			citations = append(citations, i+1)
			if len(hits) == 3 {
				break
			}
		}
	}

	// Fallback to first 3 pages if no match
	if len(hits) == 0 {
		for i := 0; i < len(pages) && i < 3; i++ {
			hits = append(hits, fmt.Sprintf("[Page %d]\n%s", i+1, pages[i]))
			citations = append(citations, i+1)
		}
	}

	// Get answer from Groq
	answer, err := callGroq(hits, req.Question)
	if err != nil {
		http.Error(w, "AI failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(AskResponse{Answer: answer, Citations: citations});err != nil {
		log.Println("failed to write response for ask method:",err)
	}
}

// callGroq sends the question and PDF content to Groq API and returns the answer
func callGroq(pages []string, question string) (string, error) {
	key := os.Getenv("GROQ_API_KEY")
	if key == "" {
		return "", fmt.Errorf("missing GROQ_API_KEY")
	}

	payload := GroqRequest{
		Model:               "llama-3.3-70b-versatile",
		Temperature:         1,
		TopP:                1,
		Stream:              false,
		MaxCompletionTokens: 1024,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Pages:\n%s\n\nQuestion: %s", strings.Join(pages, "\n\n"), question),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Groq error: %s", string(msg))
	}

	var parsed struct {
		Choices []struct {
			Message ChatMessage `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("Empty AI response")
	}

	return parsed.Choices[0].Message.Content, nil
}

// extractText reads each page of the uploaded PDF and extracts plain text
func extractText(path string) ([]string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open PDF: %w", err)
	}
	defer f.Close()

	var pages []string
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		text, err := page.GetPlainText(nil)
		if err != nil {
			pages = append(pages, "")
		} else {
			pages = append(pages, strings.TrimSpace(text))
		}
	}
	return pages, nil
}
