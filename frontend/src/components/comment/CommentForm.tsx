import React, { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createCommentSchema } from '../../utils/validation';
import type { CreateCommentFormData } from '../../utils/validation';
import * as commentsApi from '../../api/endpoints/comments';
import { CHAR_LIMITS } from '../../utils/constants';

interface CommentFormProps {
  postId: string;
  parentCommentId?: string;
  onSuccess?: () => void;
}

const CommentForm: React.FC<CommentFormProps> = ({ postId, parentCommentId, onSuccess }) => {
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateCommentFormData>({
    resolver: yupResolver(createCommentSchema),
    defaultValues: {
      content: '',
      parent_comment_id: parentCommentId,
    },
  });

  const content = watch('content');
  const charCount = content?.length || 0;
  const remainingChars = CHAR_LIMITS.COMMENT_CONTENT - charCount;

  const createCommentMutation = useMutation({
    mutationFn: (data: CreateCommentFormData) =>
      commentsApi.createComment(postId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postId] });
      queryClient.invalidateQueries({ queryKey: ['post', postId] });
      reset();
      setError(null);
      if (onSuccess) {
        onSuccess();
      }
    },
    onError: (err: Error) => {
      setError(err.message || 'コメントの投稿に失敗しました');
    },
  });

  const onSubmit = async (data: CreateCommentFormData) => {
    try {
      setError(null);
      await createCommentMutation.mutateAsync(data);
    } catch (err) {
      // エラーは mutation の onError で処理
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mb: 2 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <TextField
        fullWidth
        multiline
        minRows={2}
        placeholder={parentCommentId ? '返信を入力...' : 'コメントを入力...'}
        {...register('content')}
        error={!!errors.content}
        helperText={errors.content?.message || `${remainingChars} / ${CHAR_LIMITS.COMMENT_CONTENT}`}
        inputProps={{ 'data-testid': 'comment-form-content' }}
        sx={{
          mb: 1,
          '& .MuiOutlinedInput-root': {
            '& fieldset': {
              borderColor: 'divider',
            },
            '&:hover fieldset': {
              borderColor: 'divider',
            },
            '&.Mui-focused fieldset': {
              borderColor: 'primary.main',
            },
          },
        }}
      />

      <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button
          type="submit"
          variant="contained"
          disabled={createCommentMutation.isPending || charCount === 0 || remainingChars < 0}
          sx={{
            borderRadius: 25,
            textTransform: 'none',
            fontWeight: 700,
            px: 3,
          }}
          data-testid="comment-form-submit"
        >
          {createCommentMutation.isPending ? (
            <CircularProgress size={20} color="inherit" />
          ) : parentCommentId ? (
            '返信'
          ) : (
            'コメント'
          )}
        </Button>
      </Box>
    </Box>
  );
};

export default CommentForm;
