// main.go (Updated to fix struct literal issue)
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

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Messages             []ChatMessage `json:"messages"`
	Model                string        `json:"model"`
	Temperature          float64       `json:"temperature"`
	TopP                 float64       `json:"top_p"`
	Stream               bool          `json:"stream"`
	Stop                 interface{}   `json:"stop"`
	MaxCompletionTokens  int           `json:"max_completion_tokens"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Could not load .env file (continuing anyway)")
	}

	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/ask", handleAsk)

	// Serve frontend
	http.Handle("/", http.FileServer(http.Dir("frontend/out")))

	// Serve uploaded PDFs
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := filepath.Join("uploads", handler.Filename)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Uploaded",
		"pdfPath":  "/" + path,
		"filename": handler.Filename,
	})
}

func extractText(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pages []string
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		text, err := p.GetPlainText(nil)
		if err != nil {
			pages = append(pages, "")
		} else {
			pages = append(pages, strings.TrimSpace(text))
		}
	}
	return pages, nil
}

func handleAsk(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Question string `json:"question"`
		PdfPath  string `json:"pdfPath"`
	}
	type Res struct {
		Answer    string `json:"answer"`
		Citations []int  `json:"citations"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" || req.PdfPath == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	pages, err := extractText(strings.TrimPrefix(req.PdfPath, "/"))
	if err != nil {
		http.Error(w, "PDF parse failed", http.StatusInternalServerError)
		return
	}

	var hits []string
	var cites []int
	q := strings.ToLower(req.Question)
	for i, p := range pages {
		if strings.Contains(strings.ToLower(p), q) {
			hits = append(hits, fmt.Sprintf("[Page %d]\n%s", i+1, p))
			cites = append(cites, i+1)
			if len(hits) == 3 {
				break
			}
		}
	}

	if len(hits) == 0 {
		for i := 0; i < len(pages) && i < 3; i++ {
			hits = append(hits, fmt.Sprintf("[Page %d]\n%s", i+1, pages[i]))
			cites = append(cites, i+1)
		}
	}

	groqPayload := GroqRequest{
		Model:               "llama-3.3-70b-versatile",
		Temperature:         1,
		TopP:                1,
		Stream:              false,
		Stop:                nil,
		MaxCompletionTokens: 1024,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Pages:\n%s\n\nQuestion: %s", strings.Join(hits, "\n\n"), req.Question),
			},
		},
	}

	body, err := json.Marshal(groqPayload)
	if err != nil {
		http.Error(w, "Request serialization error", http.StatusInternalServerError)
		return
	}

	groqKey := os.Getenv("GROQ_API_KEY")
	if groqKey == "" {
		http.Error(w, "Missing GROQ_API_KEY", http.StatusInternalServerError)
		return
	}

	reqAPI, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	reqAPI.Header.Set("Content-Type", "application/json")
	reqAPI.Header.Set("Authorization", "Bearer "+groqKey)

	resp, err := http.DefaultClient.Do(reqAPI)
	if err != nil {
		http.Error(w, "Groq API error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Groq Error Response: %s\n", string(msg))
		http.Error(w, "AI response error: "+string(msg), http.StatusInternalServerError)
		return
	}

	var parsed struct {
		Choices []struct {
			Message ChatMessage `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		http.Error(w, "AI response decode error", http.StatusInternalServerError)
		return
	}

	if len(parsed.Choices) == 0 {
		http.Error(w, "Empty AI response", http.StatusInternalServerError)
		return
	}

	answer := parsed.Choices[0].Message.Content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Res{Answer: answer, Citations: cites})
}
