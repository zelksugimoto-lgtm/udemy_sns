import React, { useState } from 'react';
import { useNavigate, useSearchParams, Link as RouterLink } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Link,
  Box,
  Alert,
  CircularProgress,
} from '@mui/material';
import { apiClient, getErrorMessage } from '../api/client';

interface ResetPasswordFormData {
  password: string;
  confirmPassword: string;
}

const ResetPassword: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const token = searchParams.get('token');

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<ResetPasswordFormData>();

  const password = watch('password');

  const onSubmit = async (data: ResetPasswordFormData) => {
    if (!token) {
      setError('リセットトークンが見つかりません');
      return;
    }

    try {
      setIsLoading(true);
      setError('');

      await apiClient.post('/auth/password-reset', {
        token,
        new_password: data.password,
      });

      setSuccess(true);

      // 3秒後にログイン画面へリダイレクト
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  if (!token) {
    return (
      <Container maxWidth="sm" sx={{ mt: 8 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Alert severity="error">
            無効なリセットリンクです。パスワードリセット申請から再度お試しください。
          </Alert>
          <Box textAlign="center" sx={{ mt: 2 }}>
            <Link component={RouterLink} to="/password-reset/request">
              パスワードリセット申請へ
            </Link>
          </Box>
        </Paper>
      </Container>
    );
  }

  if (success) {
    return (
      <Container maxWidth="sm" sx={{ mt: 8 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Alert severity="success">
            パスワードのリセットが完了しました。新しいパスワードでログインしてください。
          </Alert>
          <Box textAlign="center" sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary">
              3秒後にログイン画面へ自動的に移動します...
            </Typography>
          </Box>
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h1" align="center" gutterBottom>
          パスワードリセット
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" gutterBottom>
          新しいパスワードを入力してください
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 3 }}>
          <TextField
            fullWidth
            label="新しいパスワード"
            type="password"
            autoComplete="new-password"
            {...register('password', {
              required: 'パスワードを入力してください',
              minLength: {
                value: 8,
                message: 'パスワードは8文字以上で入力してください',
              },
            })}
            error={!!errors.password}
            helperText={errors.password?.message || '8文字以上で入力してください'}
            sx={{ mb: 2 }}
            inputProps={{ 'data-testid': 'new-password' }}
          />

          <TextField
            fullWidth
            label="パスワード確認"
            type="password"
            autoComplete="new-password"
            {...register('confirmPassword', {
              required: 'パスワード確認を入力してください',
              validate: (value) =>
                value === password || 'パスワードが一致しません',
            })}
            error={!!errors.confirmPassword}
            helperText={errors.confirmPassword?.message}
            sx={{ mb: 3 }}
            inputProps={{ 'data-testid': 'confirm-password' }}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={isLoading}
            sx={{ mb: 2 }}
            data-testid="reset-password-submit"
          >
            {isLoading ? <CircularProgress size={24} /> : 'パスワードをリセット'}
          </Button>

          <Box textAlign="center">
            <Typography variant="body2">
              <Link component={RouterLink} to="/login">
                ログイン画面に戻る
              </Link>
            </Typography>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default ResetPassword;
