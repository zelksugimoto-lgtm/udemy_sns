import React, { useRef, useState } from 'react';
import {
  Box,
  IconButton,
  CircularProgress,
  Alert,
} from '@mui/material';
import {
  Image as ImageIcon,
  Close as CloseIcon,
} from '@mui/icons-material';
import {
  validateFiles,
  uploadMultipleImages,
  formatFileSize,
  MAX_UPLOAD_FILES,
  type UploadedFile,
} from '../../utils/fileUpload';

interface MediaUploadProps {
  onUploadComplete: (files: UploadedFile[]) => void;
  disabled?: boolean;
}

interface PreviewFile {
  file: File;
  preview: string;
}

const MediaUpload: React.FC<MediaUploadProps> = ({ onUploadComplete, disabled = false }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewFiles, setPreviewFiles] = useState<PreviewFile[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);

    if (files.length === 0) return;

    // バリデーション
    const validation = validateFiles(files);
    if (!validation.valid) {
      setError(validation.error || '');
      return;
    }

    // プレビュー作成
    const newPreviews: PreviewFile[] = files.map((file) => ({
      file,
      preview: URL.createObjectURL(file),
    }));

    setPreviewFiles((prev) => [...prev, ...newPreviews].slice(0, MAX_UPLOAD_FILES));
    setError(null);

    // input要素をリセット（同じファイルを再度選択できるように）
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleRemoveFile = (index: number) => {
    setPreviewFiles((prev) => {
      const newPreviews = [...prev];
      // プレビューURLをクリーンアップ
      URL.revokeObjectURL(newPreviews[index].preview);
      newPreviews.splice(index, 1);
      return newPreviews;
    });
  };

  const handleUpload = async () => {
    if (previewFiles.length === 0) return;

    setIsUploading(true);
    setError(null);

    try {
      const files = previewFiles.map((pf) => pf.file);
      const uploadedFiles = await uploadMultipleImages(files);

      // アップロード成功後、プレビューをクリア
      previewFiles.forEach((pf) => URL.revokeObjectURL(pf.preview));
      setPreviewFiles([]);

      onUploadComplete(uploadedFiles);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '画像のアップロードに失敗しました';
      setError(errorMessage);
    } finally {
      setIsUploading(false);
    }
  };

  // プレビューが追加されたら自動的にアップロード
  React.useEffect(() => {
    if (previewFiles.length > 0 && !isUploading) {
      handleUpload();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [previewFiles.length]);

  const handleIconClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <Box>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        multiple
        onChange={handleFileSelect}
        style={{ display: 'none' }}
        disabled={disabled || isUploading}
      />

      {error && (
        <Alert severity="error" sx={{ mb: 1 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <IconButton
          color="primary"
          onClick={handleIconClick}
          disabled={disabled || isUploading || previewFiles.length >= MAX_UPLOAD_FILES}
          size="small"
          data-testid="media-upload-button"
        >
          {isUploading ? (
            <CircularProgress size={20} />
          ) : (
            <ImageIcon />
          )}
        </IconButton>

        {previewFiles.length > 0 && (
          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
            {previewFiles.map((pf, index) => (
              <Box
                key={index}
                sx={{
                  position: 'relative',
                  width: 60,
                  height: 60,
                  borderRadius: 1,
                  overflow: 'hidden',
                  border: '1px solid',
                  borderColor: 'divider',
                }}
              >
                <img
                  src={pf.preview}
                  alt={`Preview ${index + 1}`}
                  style={{
                    width: '100%',
                    height: '100%',
                    objectFit: 'cover',
                  }}
                />
                <IconButton
                  size="small"
                  sx={{
                    position: 'absolute',
                    top: 2,
                    right: 2,
                    backgroundColor: 'rgba(0, 0, 0, 0.5)',
                    color: 'white',
                    '&:hover': {
                      backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    },
                    width: 20,
                    height: 20,
                  }}
                  onClick={() => handleRemoveFile(index)}
                  disabled={isUploading}
                >
                  <CloseIcon sx={{ fontSize: 14 }} />
                </IconButton>
                <Box
                  sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    backgroundColor: 'rgba(0, 0, 0, 0.6)',
                    color: 'white',
                    fontSize: '0.6rem',
                    textAlign: 'center',
                    py: 0.25,
                  }}
                >
                  {formatFileSize(pf.file.size)}
                </Box>
              </Box>
            ))}
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default MediaUpload;
