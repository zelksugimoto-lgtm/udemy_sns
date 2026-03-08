import React, { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Paper,
  Avatar,
  Typography,
  CircularProgress,
  Alert,
  IconButton,
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../../contexts/AuthContext';
import { createPostSchema } from '../../utils/validation';
import type { CreatePostFormData } from '../../utils/validation';
import * as postsApi from '../../api/endpoints/posts';
import { CHAR_LIMITS } from '../../utils/constants';
import MediaUpload from './MediaUpload';
import type { UploadedFile } from '../../utils/fileUpload';

const PostForm: React.FC = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);
  const [uploadedMedia, setUploadedMedia] = useState<UploadedFile[]>([]);

  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreatePostFormData>({
    resolver: yupResolver(createPostSchema),
    defaultValues: {
      content: '',
      visibility: 'public',
    },
  });

  const content = watch('content');
  const charCount = content?.length || 0;
  const remainingChars = CHAR_LIMITS.POST_CONTENT - charCount;

  const createPostMutation = useMutation({
    mutationFn: postsApi.createPost,
    onSuccess: () => {
      // タイムラインを再取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      reset();
      setUploadedMedia([]);
      setError(null);
    },
    onError: (err: Error) => {
      setError(err.message || '投稿に失敗しました');
    },
  });

  const onSubmit = async (data: CreatePostFormData) => {
    try {
      setError(null);

      // メディアがある場合は追加
      const postData = {
        ...data,
        ...(uploadedMedia.length > 0 && {
          media: uploadedMedia.map((media) => ({
            media_type: media.media_type,
            media_url: media.media_url,
            display_order: media.display_order,
          })),
        }),
      };

      await createPostMutation.mutateAsync(postData);
    } catch (err) {
      // エラーは mutation の onError で処理
    }
  };

  const handleMediaUploadComplete = (files: UploadedFile[]) => {
    setUploadedMedia((prev) => [...prev, ...files]);
  };

  const handleRemoveMedia = (index: number) => {
    setUploadedMedia((prev) => prev.filter((_, i) => i !== index));
  };

  if (!user) {
    return null;
  }

  return (
    <Paper
      elevation={0}
      sx={{
        p: 2,
        borderBottom: '1px solid',
        borderColor: 'divider',
      }}
    >
      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Box sx={{ display: 'flex', gap: 2 }}>
        <Avatar
          alt={user.display_name}
          src={user.avatar_url || undefined}
          sx={{ width: 48, height: 48 }}
        >
          {user.display_name.charAt(0).toUpperCase()}
        </Avatar>

        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ flex: 1 }}>
          <TextField
            fullWidth
            multiline
            minRows={3}
            placeholder="いまどうしてる？"
            {...register('content')}
            error={!!errors.content}
            helperText={errors.content?.message}
            variant="standard"
            InputProps={{
              disableUnderline: true,
            }}
            inputProps={{ 'data-testid': 'post-form-content' }}
            sx={{
              '& .MuiInputBase-input': {
                fontSize: '1.1rem',
              },
            }}
          />

          {/* Media Upload */}
          <Box sx={{ mt: 2, mb: 1 }}>
            <MediaUpload
              onUploadComplete={handleMediaUploadComplete}
              disabled={createPostMutation.isPending || uploadedMedia.length >= 4}
            />
          </Box>

          {/* Uploaded Media Preview */}
          {uploadedMedia.length > 0 && (
            <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
              {uploadedMedia.map((media, index) => (
                <Box
                  key={index}
                  sx={{
                    position: 'relative',
                    width: 80,
                    height: 80,
                    borderRadius: 1,
                    overflow: 'hidden',
                    border: '1px solid',
                    borderColor: 'divider',
                  }}
                >
                  <img
                    src={media.media_url}
                    alt={`Uploaded ${index + 1}`}
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
                    onClick={() => handleRemoveMedia(index)}
                    disabled={createPostMutation.isPending}
                  >
                    <CloseIcon sx={{ fontSize: 14 }} />
                  </IconButton>
                </Box>
              ))}
            </Box>
          )}

          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              mt: 2,
            }}
          >
            <Typography
              variant="caption"
              color={remainingChars < 0 ? 'error' : 'text.secondary'}
            >
              {remainingChars < 0 && '-'}
              {Math.abs(remainingChars)} / {CHAR_LIMITS.POST_CONTENT}
            </Typography>

            <Button
              type="submit"
              variant="contained"
              disabled={createPostMutation.isPending || charCount === 0 || remainingChars < 0}
              sx={{
                borderRadius: 25,
                textTransform: 'none',
                fontWeight: 700,
                px: 3,
              }}
              data-testid="post-form-submit"
            >
              {createPostMutation.isPending ? (
                <CircularProgress size={20} color="inherit" />
              ) : (
                '投稿する'
              )}
            </Button>
          </Box>
        </Box>
      </Box>
    </Paper>
  );
};

export default PostForm;
