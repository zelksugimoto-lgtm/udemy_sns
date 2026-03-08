import axios from 'axios';

// 画像形式の定義
export const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];

// ファイルサイズ制限（バイト単位）
const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB
const MAX_VIDEO_SIZE = 20 * 1024 * 1024; // 20MB

// 最大アップロード数
export const MAX_UPLOAD_FILES = 4;

export interface UploadedFile {
  media_url: string;
  media_type: 'image' | 'video';
  display_order: number;
}

/**
 * ファイルの形式を検証
 */
export function validateFileType(file: File): { valid: boolean; error?: string } {
  if (ALLOWED_IMAGE_TYPES.includes(file.type)) {
    return { valid: true };
  }

  return { valid: false, error: '対応していない画像形式です。JPEG、PNG、GIF、WebPのみアップロード可能です。' };
}

/**
 * ファイルのサイズを検証
 */
export function validateFileSize(file: File): { valid: boolean; error?: string } {
  const isImage = ALLOWED_IMAGE_TYPES.includes(file.type);
  const maxSize = isImage ? MAX_IMAGE_SIZE : MAX_VIDEO_SIZE;

  if (file.size > maxSize) {
    const maxSizeMB = maxSize / (1024 * 1024);
    return { valid: false, error: `ファイルサイズは${maxSizeMB}MB以下にしてください。` };
  }

  return { valid: true };
}

/**
 * 複数ファイルのバリデーション
 */
export function validateFiles(files: File[]): { valid: boolean; error?: string } {
  if (files.length === 0) {
    return { valid: true };
  }

  if (files.length > MAX_UPLOAD_FILES) {
    return { valid: false, error: `画像は最大${MAX_UPLOAD_FILES}枚までアップロードできます。` };
  }

  for (const file of files) {
    const typeValidation = validateFileType(file);
    if (!typeValidation.valid) {
      return typeValidation;
    }

    const sizeValidation = validateFileSize(file);
    if (!sizeValidation.valid) {
      return sizeValidation;
    }
  }

  return { valid: true };
}

/**
 * バックエンドに画像をアップロード
 * @param file アップロードするファイル
 * @param displayOrder 表示順序
 * @returns アップロードされたファイル情報
 */
export async function uploadImageToBackend(
  file: File,
  displayOrder: number
): Promise<UploadedFile> {
  // ファイル形式検証
  const typeValidation = validateFileType(file);
  if (!typeValidation.valid) {
    throw new Error(typeValidation.error);
  }

  // ファイルサイズ検証
  const sizeValidation = validateFileSize(file);
  if (!sizeValidation.valid) {
    throw new Error(sizeValidation.error);
  }

  // FormDataを作成
  const formData = new FormData();
  formData.append('file', file);

  try {
    // バックエンドのアップロードエンドポイントにPOST
    const response = await axios.post<{ media_url: string; media_type: string }>(
      `${import.meta.env.VITE_API_BASE_URL}/upload/image`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        withCredentials: true,
      }
    );

    return {
      media_url: response.data.media_url,
      media_type: response.data.media_type as 'image' | 'video',
      display_order: displayOrder,
    };
  } catch (error) {
    console.error('画像のアップロードに失敗しました:', error);
    if (axios.isAxiosError(error) && error.response) {
      const errorData = error.response.data as { error?: { message?: string } };
      throw new Error(errorData.error?.message || '画像のアップロードに失敗しました。');
    }
    throw new Error('画像のアップロードに失敗しました。もう一度お試しください。');
  }
}

/**
 * 複数の画像をアップロード
 * @param files アップロードするファイルの配列
 * @returns アップロードされたファイル情報の配列
 */
export async function uploadMultipleImages(
  files: File[]
): Promise<UploadedFile[]> {
  // ファイル全体のバリデーション
  const validation = validateFiles(files);
  if (!validation.valid) {
    throw new Error(validation.error);
  }

  // 並列アップロード
  const uploadPromises = files.map((file, index) =>
    uploadImageToBackend(file, index)
  );

  try {
    return await Promise.all(uploadPromises);
  } catch (error) {
    console.error('複数画像のアップロードに失敗しました:', error);
    throw error;
  }
}

/**
 * ファイルサイズを人間が読める形式に変換
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}
