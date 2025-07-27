import {
  useDisclosure,
  Button,
  Dialog,
  DialogPositioner,
  DialogContent,
  DialogBody,
  VStack,
  Heading,
  Box,
  Progress,
  HStack,
  Input,
  Text,
} from "@chakra-ui/react";
import { Fragment, FunctionComponent, useRef, useState } from "react";

export interface UploadDialogProps {
  onUpload: () => void;
  uploadProgress: number;
  setFile: (file: File | null) => void;
}

export const UploadDialog: FunctionComponent<UploadDialogProps> = ({
  onUpload,
  uploadProgress,
  setFile,
}) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const { open, onOpen, onClose } = useDisclosure();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setSelectedFile(file);
    setFile(file);
  };

  const handleDialogClose = () => {
    setSelectedFile(null);
    setFile(null);
    onClose();
  };

  const handleUpload = () => {
    onUpload();
    handleDialogClose();
  };

  return (
    <Fragment>
      <Button colorScheme="blue" onClick={onOpen}>
        Upload PDF
      </Button>

      <Dialog.Root open={open} closeOnEscape onExitComplete={handleDialogClose}>
        <DialogPositioner
          style={{
            backdropFilter: "blur(4px)",
            backgroundColor: "rgba(0, 0, 0, 0.3)",
          }}
        >
          <DialogContent borderRadius="md" p={6}>
            <DialogBody>
              <VStack gap={4} align="stretch">
                <Heading size="md" color="black">
                  Upload a PDF
                </Heading>

                <Input
                  type="file"
                  accept="application/pdf"
                  ref={inputRef}
                  onChange={handleFileChange}
                  color="black"
                />

                {uploadProgress > 0 && (
                  <Box w="100%">
                    <Progress.Root value={uploadProgress} size="sm" />
                    <Text fontSize="sm" mt={1} color="black">
                      {uploadProgress}%
                    </Text>
                  </Box>
                )}

                <HStack justify="flex-end">
                  <Button variant="outline" onClick={handleDialogClose}>
                    Cancel
                  </Button>
                  <Button
                    colorScheme="blue"
                    onClick={handleUpload}
                    disabled={!selectedFile}
                  >
                    Upload
                  </Button>
                </HStack>
              </VStack>
            </DialogBody>
          </DialogContent>
        </DialogPositioner>
      </Dialog.Root>
    </Fragment>
  );
};

export interface PDFViewerProps {
  pdfPath: string;
  currentPage: number;
}

export const PDFViewer: FunctionComponent<PDFViewerProps> = ({
  pdfPath,
  currentPage,
}) => {
  return (
    <Box bg="white" shadow="md" borderRadius="md" overflow="hidden" h="80vh">
      {pdfPath ? (
        <iframe
          key={`${pdfPath}#page=${currentPage}`}
          src={`${pdfPath}#page=${currentPage}`}
          title="PDF Viewer"
          style={{ width: "100%", height: "100%", border: "none" }}
        />
      ) : (
        <Box p={6} textAlign="center" color="gray.500">
          No PDF uploaded yet.
        </Box>
      )}
    </Box>
  );
};
