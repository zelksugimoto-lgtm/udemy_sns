import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation, useSearchParams } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
  Alert,
} from '@mui/material';
import HourglassEmptyIcon from '@mui/icons-material/HourglassEmpty';
import BlockIcon from '@mui/icons-material/Block';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import { useAuth } from '../contexts/AuthContext';
import { apiClient } from '../api/client';

const PendingApproval: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const { logout } = useAuth();
  const [isApproved, setIsApproved] = useState(false);

  // クエリパラメータまたはlocation.stateからユーザー情報を取得
  const userStatus = searchParams.get('status') || (location.state as { status?: string })?.status || 'pending';
  const userEmail = searchParams.get('email') || (location.state as { email?: string })?.email || '';

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const handlePasswordReset = () => {
    navigate('/password-reset/request', { state: { email: userEmail } });
  };

  // 定期的にユーザーステータスをチェック（承認されたら通知）
  useEffect(() => {
    const checkApprovalStatus = async () => {
      try {
        const response = await apiClient.get('/auth/me');
        const user = response.data; // /auth/me は直接UserResponseを返す

        // ユーザーが承認済みになっていたら状態を更新
        if (user.status === 'approved') {
          setIsApproved(true);
        }
      } catch (error) {
        // エラーは無視（ログアウトされている場合など）
      }
    };

    // 初回ロード時にステータスをチェック
    checkApprovalStatus();

    // 5秒ごとにステータスをチェック
    const intervalId = setInterval(checkApprovalStatus, 5000);

    // コンポーネントのアンマウント時にインターバルをクリア
    return () => clearInterval(intervalId);
  }, []);

  // 承認済みになったら自動的にログアウト→ログインページへリダイレクト
  useEffect(() => {
    if (isApproved) {
      // 3秒後にログアウトしてログインページへ
      const timeoutId = setTimeout(async () => {
        await logout();
        navigate('/login', {
          state: { message: 'アカウントが承認されました。再度ログインしてください。' },
        });
      }, 3000);

      return () => clearTimeout(timeoutId);
    }
  }, [isApproved, logout, navigate]);

  // ステータスに応じた表示内容
  const isRejected = userStatus === 'rejected';

  let Icon = HourglassEmptyIcon;
  let iconColor = 'warning.main';
  let title = '承認待ち';
  let subtitle = 'アカウントは現在承認待ちです';
  let alertSeverity: 'error' | 'info' | 'success' = 'info';
  let alertMessage = '管理者がアカウントを確認中です。承認されるまでしばらくお待ちください。';

  if (isApproved) {
    Icon = CheckCircleIcon;
    iconColor = 'success.main';
    title = 'アカウント承認済み';
    subtitle = 'アカウントが承認されました';
    alertSeverity = 'success';
    alertMessage = 'アカウントが承認されました。ログアウトして再度ログインしてください。';
  } else if (isRejected) {
    Icon = BlockIcon;
    iconColor = 'error.main';
    title = 'アカウント拒否';
    subtitle = 'アカウントは拒否されました';
    alertSeverity = 'error';
    alertMessage = 'アカウントが拒否されました。パスワードリセット申請を行うか、新しいアカウントを作成してください。';
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Box display="flex" flexDirection="column" alignItems="center">
          <Icon sx={{ fontSize: 64, color: iconColor, mb: 2 }} />
          <Typography variant="h1" align="center" gutterBottom>
            {title}
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary" paragraph>
            {subtitle}
          </Typography>
        </Box>

        <Alert severity={alertSeverity} sx={{ mt: 2, mb: 3 }}>
          <Typography variant="body2">
            {alertMessage}
          </Typography>
          {userEmail && (
            <Typography variant="body2" sx={{ mt: 1 }}>
              登録メールアドレス: <strong>{userEmail}</strong>
            </Typography>
          )}
        </Alert>

        <Box sx={{ mt: 3 }}>
          {isApproved ? (
            <>
              <Typography variant="body2" color="text.secondary" paragraph>
                新しいトークンを取得するため、再度ログインしてください。
              </Typography>
              <Button
                fullWidth
                variant="contained"
                onClick={handleLogout}
                data-testid="pending-relogin"
              >
                再ログイン
              </Button>
            </>
          ) : (
            <>
              <Typography variant="body2" color="text.secondary" paragraph>
                承認に時間がかかる場合や、アカウントにアクセスできない場合は、
                パスワードリセット申請をご利用ください。
              </Typography>

              <Button
                fullWidth
                variant="outlined"
                onClick={handlePasswordReset}
                sx={{ mb: 2 }}
                data-testid="pending-password-reset"
              >
                パスワードリセット申請
              </Button>

              <Button
                fullWidth
                variant="contained"
                onClick={handleLogout}
                data-testid="pending-logout"
              >
                ログアウト
              </Button>
            </>
          )}
        </Box>
      </Paper>
    </Container>
  );
};

export default PendingApproval;
