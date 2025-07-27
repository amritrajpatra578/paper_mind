import { QAItem } from "@/pages/api/pdf";
import { Box, HStack, Button, VStack, Textarea, Text } from "@chakra-ui/react";
import { FunctionComponent } from "react";

export interface ChatInputProps {
  question: string;
  setQuestion: (val: string) => void;
  onAsk: () => void;
  loading: boolean;
}

export const ChatInput: FunctionComponent<ChatInputProps> = ({
  question,
  setQuestion,
  onAsk,
  loading,
}) => {
  return (
    <Box borderTop="1px solid" borderColor="gray.200" p={4} bg="white">
      <VStack gap={2} align="stretch">
        <Textarea
          placeholder="Ask a question about the PDF..."
          resize="none"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          color="black"
        />
        <Button
          colorScheme="green"
          onClick={onAsk}
          loading={loading}
          loadingText="Thinking..."
        >
          Ask
        </Button>
      </VStack>
    </Box>
  );
};

export interface ChatPanelProps {
  chatHistory: QAItem[];
  onCitationClick: (page: number) => void;
}

export const ChatPanel: FunctionComponent<ChatPanelProps> = ({
  chatHistory,
  onCitationClick,
}) => {
  return (
    <Box
      bg="white"
      shadow="md"
      borderRadius="md"
      display="flex"
      flexDirection="column"
      h="80vh"
    >
      <Box flex="1" overflowY="auto" p={4}>
        {chatHistory.length === 0 ? (
          <Text color="black" fontSize="sm">
            Ask a question about the document to get started.
          </Text>
        ) : (
          chatHistory.map((item, idx) => (
            <Box key={idx} mb={6}>
              <Text fontWeight="medium" color="purple" mb={1}>
                Q: {item.question}
              </Text>
              <Box
                bg="gray.100"
                p={3}
                borderRadius="md"
                whiteSpace="pre-wrap"
                fontSize="sm"
                color="black"
              >
                {item.answer}
              </Box>
              {item.citations.length > 0 && (
                <HStack mt={2} gap={2} flexWrap="wrap">
                  {item.citations.map((page) => (
                    <Button
                      key={page}
                      size="xs"
                      colorScheme="blue"
                      variant="outline"
                      onClick={() => onCitationClick(page)}
                    >
                      Page {page}
                    </Button>
                  ))}
                </HStack>
              )}
            </Box>
          ))
        )}
      </Box>
    </Box>
  );
};
