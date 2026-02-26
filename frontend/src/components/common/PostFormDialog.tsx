import React, { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Box,
  TextField,
  Button,
  Avatar,
  Typography,
  CircularProgress,
  Alert,
  Slide,
} from '@mui/material';
import type { SlideProps } from '@mui/material/Slide';
import { Close as CloseIcon } from '@mui/icons-material';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../../contexts/AuthContext';
import { createPostSchema } from '../../utils/validation';
import type { CreatePostFormData } from '../../utils/validation';
import * as postsApi from '../../api/endpoints/posts';
import { CHAR_LIMITS } from '../../utils/constants';

interface PostFormDialogProps {
  open: boolean;
  onClose: () => void;
}

const Transition = React.forwardRef(function Transition(
  props: SlideProps,
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

const PostFormDialog: React.FC<PostFormDialogProps> = ({ open, onClose }) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);

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
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      reset();
      setError(null);
      onClose();
    },
    onError: (err: Error) => {
      setError(err.message || '投稿に失敗しました');
    },
  });

  const onSubmit = async (data: CreatePostFormData) => {
    try {
      setError(null);
      await createPostMutation.mutateAsync(data);
    } catch (err) {
      // エラーは mutation の onError で処理
    }
  };

  const handleClose = () => {
    if (!createPostMutation.isPending) {
      reset();
      setError(null);
      onClose();
    }
  };

  if (!user) {
    return null;
  }

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      TransitionComponent={Transition}
      fullWidth
      maxWidth="sm"
      PaperProps={{
        sx: {
          position: 'fixed',
          bottom: 0,
          m: 0,
          width: '100%',
          maxWidth: '100%',
          borderRadius: '16px 16px 0 0',
          maxHeight: '80vh',
        },
      }}
    >
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1, pb: 1 }}>
        <IconButton
          edge="start"
          color="inherit"
          onClick={handleClose}
          aria-label="閉じる"
          disabled={createPostMutation.isPending}
        >
          <CloseIcon />
        </IconButton>
        <Typography variant="h6" component="span" sx={{ flex: 1, textAlign: 'center', mr: 5 }}>
          新規投稿
        </Typography>
      </DialogTitle>

      <DialogContent>
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
              minRows={4}
              maxRows={10}
              placeholder="いまどうしてる？"
              {...register('content')}
              error={!!errors.content}
              helperText={errors.content?.message}
              variant="standard"
              autoFocus
              InputProps={{
                disableUnderline: true,
              }}
              sx={{
                '& .MuiInputBase-input': {
                  fontSize: '1.1rem',
                },
              }}
            />

            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mt: 2,
                pb: 2,
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
      </DialogContent>
    </Dialog>
  );
};

export default PostFormDialog;
