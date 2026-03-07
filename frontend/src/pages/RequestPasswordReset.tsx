import React, { useState } from 'react';
import { useNavigate, useLocation, Link as RouterLink } from 'react-router-dom';
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

interface PasswordResetFormData {
  email: string;
  reason: string;
}

const RequestPasswordReset: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  // location.stateからメールアドレスを取得（任意）
  const prefilledEmail = (location.state as { email?: string })?.email || '';

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PasswordResetFormData>({
    defaultValues: {
      email: prefilledEmail,
      reason: '',
    },
  });

  const onSubmit = async (data: PasswordResetFormData) => {
    try {
      setIsLoading(true);
      setError('');

      await apiClient.post('/auth/password-reset/request', data);

      // 成功時、確認画面へ遷移
      navigate('/password-reset/confirmation', { state: { email: data.email } });
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h1" align="center" gutterBottom>
          パスワードリセット申請
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" gutterBottom>
          アカウントにアクセスできない場合は、パスワードリセット申請を行ってください
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 3 }}>
          <TextField
            fullWidth
            label="メールアドレス"
            type="email"
            autoComplete="email"
            {...register('email', {
              required: 'メールアドレスを入力してください',
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: '有効なメールアドレスを入力してください',
              },
            })}
            error={!!errors.email}
            helperText={errors.email?.message}
            sx={{ mb: 2 }}
            inputProps={{ 'data-testid': 'reset-email' }}
          />

          <TextField
            fullWidth
            label="申請理由"
            multiline
            rows={4}
            {...register('reason', {
              required: '申請理由を入力してください',
              minLength: {
                value: 10,
                message: '申請理由は10文字以上で入力してください',
              },
              maxLength: {
                value: 500,
                message: '申請理由は500文字以内で入力してください',
              },
            })}
            error={!!errors.reason}
            helperText={errors.reason?.message || '10文字以上500文字以内で入力してください'}
            sx={{ mb: 3 }}
            inputProps={{ 'data-testid': 'reset-reason' }}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={isLoading}
            sx={{ mb: 2 }}
            data-testid="reset-submit"
          >
            {isLoading ? <CircularProgress size={24} /> : '申請する'}
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

export default RequestPasswordReset;
