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
import { loginSchema } from '../utils/validation';
import type { LoginFormData } from '../utils/validation';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login, error: authError } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: yupResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    try {
      setIsLoading(true);
      await login(data);
      navigate('/'); // ログイン成功後、ホームへ
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
          ログイン
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" gutterBottom>
          SNS Application にログイン
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
            inputProps={{ 'data-testid': 'login-email' }}
          />

          <TextField
            fullWidth
            label="パスワード"
            type="password"
            autoComplete="current-password"
            {...register('password')}
            error={!!errors.password}
            helperText={errors.password?.message}
            sx={{ mb: 3 }}
            inputProps={{ 'data-testid': 'login-password' }}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={isLoading}
            sx={{ mb: 2 }}
            data-testid="login-submit"
          >
            {isLoading ? <CircularProgress size={24} /> : 'ログイン'}
          </Button>

          <Box textAlign="center">
            <Typography variant="body2">
              アカウントをお持ちでないですか？{' '}
              <Link component={RouterLink} to="/register">
                新規登録
              </Link>
            </Typography>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default Login;
