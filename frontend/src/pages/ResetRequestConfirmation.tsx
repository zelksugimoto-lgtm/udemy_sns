import React from 'react';
import { useNavigate, useLocation, Link as RouterLink } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Button,
  Link,
  Box,
  Alert,
} from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';

const ResetRequestConfirmation: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  // location.stateからメールアドレスを取得
  const userEmail = (location.state as { email?: string })?.email || '';

  const handleBackToLogin = () => {
    navigate('/login');
  };

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Box display="flex" flexDirection="column" alignItems="center">
          <CheckCircleOutlineIcon sx={{ fontSize: 64, color: 'success.main', mb: 2 }} />
          <Typography variant="h1" align="center" gutterBottom>
            申請完了
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary" paragraph>
            パスワードリセット申請を受け付けました
          </Typography>
        </Box>

        <Alert severity="success" sx={{ mt: 2, mb: 3 }}>
          <Typography variant="body2">
            管理者がパスワードリセット申請を確認次第、ご連絡いたします。
          </Typography>
          {userEmail && (
            <Typography variant="body2" sx={{ mt: 1 }}>
              申請メールアドレス: <strong>{userEmail}</strong>
            </Typography>
          )}
        </Alert>

        <Box sx={{ mt: 3 }}>
          <Typography variant="body2" color="text.secondary" paragraph>
            管理者からの連絡をお待ちください。承認されると、パスワードリセット用のリンクが発行されます。
          </Typography>

          <Button
            fullWidth
            variant="contained"
            onClick={handleBackToLogin}
            data-testid="confirmation-back-to-login"
          >
            ログイン画面に戻る
          </Button>

          <Box textAlign="center" sx={{ mt: 2 }}>
            <Typography variant="body2">
              <Link component={RouterLink} to="/register">
                新規アカウント登録
              </Link>
            </Typography>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default ResetRequestConfirmation;
