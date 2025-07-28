---
---

# PaperMind

A lightweight web app that lets you **upload a PDF**, **ask questions**, and get **AI-powered answers with citations** to exact pages. Built with:

- Groq LLaMA-3 for smart responses
- Go backend with PDF parsing
- Simple React frontend
- Caddy for local HTTPS/static serving

---

## Live Demo

you must need to login vercel to access the url

- **url**: [https://paper-mind-git-deployment-amritrajpatra578s-projects.vercel.app/](https://paper-mind-git-deployment-amritrajpatra578s-projects.vercel.app/)

- **Frontend (Vercel)**: [https://paper-mind.vercel.app](https://paper-mind.vercel.app)
- **Backend (Railway)**: [https://papermind-production.up.railway.app/](https://papermind-production.up.railway.app/)

> Make sure to open the frontend link for usage.

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

3. **Start the caddy server**

   ```bash
   caddy run
   ```

   Visit for installation: [https://caddyserver.com/docs/install](https://caddyserver.com/docs/install)

---

4. **Open your browser at following address**

   [http://127.0.0.1:9000](http://127.0.0.1:9000)

---

## Project Structure

```
paperMind/
├── backend/          # Go backend server
│   └── uploads/      # Stored PDFs
│   └── main.go
├── frontend/         # React frontend
│   ├── Caddyfile     # procy server from backend and frontend
│   ├── pages/
│   ├── components/
│   └── ...
└── README.md
```

---

## License

MIT © [amritrajpatra578](https://github.com/amritrajpatra578)

---
