import axios from "axios";

// Types
export interface QAItem {
  question: string;
  answer: string;
  citations: number[];
}

export interface UploadResponse {
  pdfPath: string;
}

export interface AskResponse {
  answer: string;
  citations: number[];
}

export const uploadPDF = async (
  file: File,
  onProgress: (percent: number) => void
): Promise<UploadResponse> => {
  const formData = new FormData();
  formData.append("file", file); // backend expects 'file'

  const response = await axios.post<UploadResponse>("/api/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (event) => {
      const percent = Math.round((event.loaded * 100) / (event.total || 1));
      onProgress(percent);
    },
  });

  return response.data;
};

export const askQuestion = async (
  question: string,
  pdfPath: string
): Promise<AskResponse> => {
  const response = await axios.post<AskResponse>("/api/ask", {
    question,
    pdfPath,
  });

  return response.data;
};
