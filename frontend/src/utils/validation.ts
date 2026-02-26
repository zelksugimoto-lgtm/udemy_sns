import * as yup from 'yup';
import { CHAR_LIMITS, PASSWORD_POLICY, USERNAME_POLICY } from './constants';

// ユーザー登録バリデーションスキーマ
export const registerSchema = yup.object({
  email: yup
    .string()
    .required('Email is required')
    .email('Invalid email format'),
  password: yup
    .string()
    .required('Password is required')
    .min(PASSWORD_POLICY.MIN_LENGTH, `Password must be at least ${PASSWORD_POLICY.MIN_LENGTH} characters`),
  username: yup
    .string()
    .required('Username is required')
    .min(USERNAME_POLICY.MIN_LENGTH, `Username must be at least ${USERNAME_POLICY.MIN_LENGTH} characters`)
    .max(USERNAME_POLICY.MAX_LENGTH, `Username must be at most ${USERNAME_POLICY.MAX_LENGTH} characters`)
    .matches(USERNAME_POLICY.PATTERN, 'Username can only contain letters, numbers, and underscores'),
  display_name: yup
    .string()
    .required('Display name is required')
    .min(1, 'Display name is required')
    .max(CHAR_LIMITS.DISPLAY_NAME, `Display name must be at most ${CHAR_LIMITS.DISPLAY_NAME} characters`),
});

// ログインバリデーションスキーマ
export const loginSchema = yup.object({
  email: yup
    .string()
    .required('Email is required')
    .email('Invalid email format'),
  password: yup
    .string()
    .required('Password is required'),
});

// 投稿作成バリデーションスキーマ
export const createPostSchema = yup.object({
  content: yup
    .string()
    .required('Post content is required')
    .min(1, 'Post content is required')
    .max(CHAR_LIMITS.POST_CONTENT, `Post must be at most ${CHAR_LIMITS.POST_CONTENT} characters`),
  visibility: yup
    .string()
    .oneOf(['public', 'followers', 'private'], 'Invalid visibility')
    .default('public'),
});

// 投稿更新バリデーションスキーマ
export const updatePostSchema = yup.object({
  content: yup
    .string()
    .min(1, 'Post content is required')
    .max(CHAR_LIMITS.POST_CONTENT, `Post must be at most ${CHAR_LIMITS.POST_CONTENT} characters`),
  visibility: yup
    .string()
    .oneOf(['public', 'followers', 'private'], 'Invalid visibility'),
});

// コメント作成バリデーションスキーマ
export const createCommentSchema = yup.object({
  content: yup
    .string()
    .required('Comment content is required')
    .min(1, 'Comment content is required')
    .max(CHAR_LIMITS.COMMENT_CONTENT, `Comment must be at most ${CHAR_LIMITS.COMMENT_CONTENT} characters`),
  parent_comment_id: yup.string().optional(),
});

// プロフィール更新バリデーションスキーマ
export const updateProfileSchema = yup.object({
  display_name: yup
    .string()
    .min(1, 'Display name is required')
    .max(CHAR_LIMITS.DISPLAY_NAME, `Display name must be at most ${CHAR_LIMITS.DISPLAY_NAME} characters`),
  bio: yup
    .string()
    .max(CHAR_LIMITS.BIO, `Bio must be at most ${CHAR_LIMITS.BIO} characters`),
  avatar_url: yup.string().url('Invalid URL format'),
  header_url: yup.string().url('Invalid URL format'),
});

// 通報バリデーションスキーマ
export const createReportSchema = yup.object({
  target_type: yup
    .string()
    .required('Target type is required')
    .oneOf(['Post', 'Comment', 'User'], 'Invalid target type'),
  target_id: yup
    .string()
    .required('Target ID is required'),
  reason: yup
    .string()
    .required('Reason is required')
    .oneOf(['spam', 'inappropriate_content', 'harassment', 'other'], 'Invalid reason'),
  comment: yup
    .string()
    .max(1000, 'Comment must be at most 1000 characters'),
});

export type RegisterFormData = yup.InferType<typeof registerSchema>;
export type LoginFormData = yup.InferType<typeof loginSchema>;
export type CreatePostFormData = yup.InferType<typeof createPostSchema>;
export type UpdatePostFormData = yup.InferType<typeof updatePostSchema>;
export type CreateCommentFormData = yup.InferType<typeof createCommentSchema>;
export type UpdateProfileFormData = yup.InferType<typeof updateProfileSchema>;
export type CreateReportFormData = yup.InferType<typeof createReportSchema>;
