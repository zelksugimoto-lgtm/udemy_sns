import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
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
import { useAuth } from '../contexts/AuthContext';
import { registerSchema } from '../utils/validation';
import type { RegisterFormData } from '../utils/validation';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const { register: registerUser, error: authError } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: yupResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setIsLoading(true);
      const user = await registerUser(data);

      // ユーザーのステータスに応じてリダイレクト先を変更
      if (user.status === 'pending' || user.status === 'rejected') {
        // 承認待ち画面へ
        navigate('/pending-approval', {
          state: { status: user.status, email: user.email },
        });
      } else if (user.status === 'approved') {
        // 承認済みの場合はホームへ
        navigate('/');
      } else {
        // その他のステータスの場合も承認待ち画面へ
        navigate('/pending-approval', {
          state: { status: user.status, email: user.email },
        });
      }
    } catch (err) {
      // エラーはAuthContextで管理されている
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h1" align="center" gutterBottom>
          新規登録
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" gutterBottom>
          SNS Application のアカウントを作成
        </Typography>

        {authError && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {authError}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 3 }}>
          <TextField
            fullWidth
            label="メールアドレス"
            type="email"
            autoComplete="email"
            {...register('email')}
            error={!!errors.email}
            helperText={errors.email?.message}
            sx={{ mb: 2 }}
            inputProps={{ 'data-testid': 'register-email' }}
          />

          <TextField
            fullWidth
            label="ユーザー名"
            autoComplete="username"
            {...register('username')}
            error={!!errors.username}
            helperText={errors.username?.message || '英数字とアンダースコアのみ使用可能'}
            sx={{ mb: 2 }}
            inputProps={{ 'data-testid': 'register-username' }}
          />

          <TextField
            fullWidth
            label="表示名"
            {...register('display_name')}
            error={!!errors.display_name}
            helperText={errors.display_name?.message}
            sx={{ mb: 2 }}
            inputProps={{ 'data-testid': 'register-displayname' }}
          />

          <TextField
            fullWidth
            label="パスワード"
            type="password"
            autoComplete="new-password"
            {...register('password')}
            error={!!errors.password}
            helperText={errors.password?.message || '8文字以上'}
            sx={{ mb: 3 }}
            inputProps={{ 'data-testid': 'register-password' }}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={isLoading}
            sx={{ mb: 2 }}
            data-testid="register-submit"
          >
            {isLoading ? <CircularProgress size={24} /> : '登録'}
          </Button>

          <Box textAlign="center">
            <Typography variant="body2">
              すでにアカウントをお持ちですか？{' '}
              <Link component={RouterLink} to="/login">
                ログイン
              </Link>
            </Typography>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default Register;
