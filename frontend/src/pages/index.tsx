import { ChatInput, ChatPanel } from "@/components/pdf/ChatPanel";
import { PDFViewer, UploadDialog } from "@/components/pdf/PDFViewer";
import { Box, Heading, HStack } from "@chakra-ui/react";
import { useState } from "react";
import {
  askQuestion,
  AskResponse,
  QAItem,
  uploadPDF,
  UploadResponse,
} from "./api/pdf";

export default function App() {
  const [pdfState, setPdfState] = useState({
    file: null as File | null,
    path: "",
    progress: 0,
    currentPage: 1,
  });

  const [chatState, setChatState] = useState({
    question: "",
    history: [] as QAItem[],
    loading: false,
  });

  const handleUpload = async () => {
    if (!pdfState.file) return;

    try {
      const data: UploadResponse = await uploadPDF(
        pdfState.file,
        (progress) => {
          setPdfState((prev) => ({ ...prev, progress }));
        }
      );

      setPdfState((prev) => ({
        ...prev,
        path: data.pdfPath,
        progress: 0,
        currentPage: 1,
      }));

      setChatState((prev) => ({ ...prev, history: [] }));
      alert("Upload successful!!");
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      alert(err.response?.data?.error || "Upload failed.");
    }
  };

  const handleAsk = async () => {
    if (!chatState.question || !pdfState.path) {
      alert("Please upload a PDF and type a question.");
      return;
    }

    try {
      setChatState((prev) => ({ ...prev, loading: true }));

      const data: AskResponse = await askQuestion(
        chatState.question,
        pdfState.path
      );

      const newItem: QAItem = {
        question: chatState.question,
        answer: data.answer,
        citations: data.citations,
      };

      setChatState((prev) => ({
        ...prev,
        history: [...prev.history, newItem],
        question: "",
        loading: false,
      }));
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      alert(err.response?.data?.error || "Question failed.");
      setChatState((prev) => ({ ...prev, loading: false }));
    }
  };

  const scrollToPDFPage = (page: number) => {
    setPdfState((prev) => ({ ...prev, currentPage: page }));
  };

  const backendBaseURL = "https://papermind-production.up.railway.app/";

  return (
    <Box minH="100vh" bg="gray.100" p={4}>
      <HStack justify="space-between" mb={4}>
        <Heading size="2xl" color="black">
          PaperMind
        </Heading>

        <UploadDialog
          onUpload={handleUpload}
          uploadProgress={pdfState.progress}
          setFile={(file) => setPdfState((prev) => ({ ...prev, file }))}
        />
      </HStack>

      <Box
        display={{ base: "block", md: "grid" }}
        gridTemplateColumns="1fr 1fr"
        gap={4}
      >
        <PDFViewer
          pdfPath={
            pdfState.path
              ? `${backendBaseURL.replace(/\/$/, "")}${pdfState.path}`
              : ""
          }
          currentPage={pdfState.currentPage}
        />
        <Box
          bg="white"
          shadow="md"
          borderRadius="md"
          display="flex"
          flexDirection="column"
          h="80vh"
        >
          <ChatPanel
            chatHistory={chatState.history}
            onCitationClick={scrollToPDFPage}
          />
          <ChatInput
            question={chatState.question}
            setQuestion={(q) =>
              setChatState((prev) => ({ ...prev, question: q }))
            }
            onAsk={handleAsk}
            loading={chatState.loading}
          />
        </Box>
      </Box>
    </Box>
  );
}
