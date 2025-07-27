Here’s your **updated `README.md`** with the **Caddy section included** and the **build section removed**. It assumes the project is at [https://github.com/amritrajpatra578/paperMind](https://github.com/amritrajpatra578/paperMind) and follows the structure we've discussed:

---

# 📄 PaperMind

A lightweight web app that lets you **upload a PDF**, **ask questions**, and get **AI-powered answers with citations** to exact pages. Built with:

- Groq LLaMA-3 for smart responses
- Go backend with PDF parsing
- Simple React frontend
- Caddy for local HTTPS/static serving

---

## Features

- Upload and preview PDF files
- Ask questions in chat format
- AI returns answers with page citations
- Click citations to jump to relevant PDF page
- Clean, responsive UI using Chakra UI v3
- Minimal external dependencies

---

## Getting Started

### Backend (Go)

1. **Clone the repo**

   ```bash
   git clone https://github.com/amritrajpatra578/paperMind.git
   cd paperMind
   ```

2. **Install dependencies**
   Make sure you have Go installed (≥1.18). Then:

   ```bash
   cd backend
   go mod tidy
   ```

3. **Run the backend**

   ```bash
   go run main.go
   ```

   This will start the server at: [http://localhost:8080](http://localhost:8080)

---

## 🌐 Frontend

1. **Install dependencies**

   ```bash
   cd frontend
   npm install
   ```

2. **Start the dev server**

   ```bash
   npm run dev
   ```

   Visit: [http://localhost:3000](http://localhost:3000)

---

## 🌍 Optional: Use Caddy to Serve Frontend + Proxy API

Caddy can serve your frontend build and reverse proxy API/PDF endpoints from the Go backend.

### 🧾 Example `Caddyfile`

```caddyfile
paperMind.localhost {
	root * frontend/out
	file_server

	# Reverse proxy API
	@api path /api/*
	reverse_proxy @api localhost:8080

	# Reverse proxy PDF access
	@uploads path /uploads/*
	reverse_proxy @uploads localhost:8080
}
```

### 🧪 Run with Caddy

1. Add to your `/etc/hosts`:

   ```
   127.0.0.1 papermind.localhost
   ```

2. Start Caddy:

   ```bash
   sudo caddy run --config Caddyfile
   ```

3. Open your browser:
   [http://papermind.localhost](http://papermind.localhost)

---

## 📁 Project Structure

```
paperMind/
├── backend/          # Go backend server
│   └── main.go
├── frontend/         # React frontend
│   ├── pages/
│   ├── components/
│   └── ...
├── uploads/          # Stored PDFs
├── Caddyfile         # For optional Caddy setup
└── README.md
```

---

## 📜 License

MIT © [amritrajpatra578](https://github.com/amritrajpatra578)

---
