import React from 'react';
import { Box, Paper, Typography, Link } from '@mui/material';

const RightSidebar: React.FC = () => {
  return (
    <Box sx={{ position: 'sticky', top: 80, width: 320, p: 2 }}>
      {/* トレンド（Phase 2で実装） */}
      <Paper elevation={0} sx={{ p: 2, mb: 2, borderRadius: 2, bgcolor: 'background.default', border: 1, borderColor: 'divider' }}>
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 700 }}>
          トレンド
        </Typography>
        <Typography variant="body2" color="text.secondary">
          準備中...
        </Typography>
      </Paper>

      {/* おすすめユーザー（Phase 2で実装） */}
      <Paper elevation={0} sx={{ p: 2, mb: 2, borderRadius: 2, bgcolor: 'background.default', border: 1, borderColor: 'divider' }}>
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 700 }}>
          おすすめユーザー
        </Typography>
        <Typography variant="body2" color="text.secondary">
          準備中...
        </Typography>
      </Paper>

      {/* フッター */}
      <Box sx={{ px: 2 }}>
        <Typography variant="caption" color="text.secondary" sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          <Link href="#" underline="hover" color="inherit">利用規約</Link>
          <span>·</span>
          <Link href="#" underline="hover" color="inherit">プライバシーポリシー</Link>
          <span>·</span>
          <Link href="#" underline="hover" color="inherit">ヘルプ</Link>
        </Typography>
        <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
          © 2025 SNS App
        </Typography>
      </Box>
    </Box>
  );
};

export default RightSidebar;
