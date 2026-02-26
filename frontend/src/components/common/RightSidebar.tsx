import React, { useState } from 'react';
import { Box, Paper, Typography, Link, TextField, InputAdornment, List, ListItem, ListItemAvatar, Avatar, ListItemText } from '@mui/material';
import { Search as SearchIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import * as usersApi from '../../api/endpoints/users';

const RightSidebar: React.FC = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');

  // ユーザー検索
  const { data: searchResults } = useQuery({
    queryKey: ['userSearch', searchQuery],
    queryFn: () => usersApi.searchUsers({ query: searchQuery, limit: 5 }),
    enabled: searchQuery.length > 0,
  });

  return (
    <Box sx={{ position: 'sticky', top: 80, width: 280, p: 2 }}>
      {/* ユーザー検索 */}
      <Paper elevation={0} sx={{ p: 1.5, mb: 2, borderRadius: 2, bgcolor: 'background.default', border: '1px solid', borderColor: 'divider' }}>
        <TextField
          fullWidth
          size="small"
          placeholder="ユーザーを検索"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon fontSize="small" />
              </InputAdornment>
            ),
          }}
          sx={{
            mb: searchResults?.users && searchResults.users.length > 0 ? 1 : 0,
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

        {searchResults?.users && searchResults.users.length > 0 && (
          <List dense disablePadding>
            {searchResults.users.map((user) => (
              <ListItem
                key={user.id}
                button
                onClick={() => {
                  navigate(`/users/${user.username}`);
                  setSearchQuery('');
                }}
                sx={{ borderRadius: 1, mb: 0.5 }}
              >
                <ListItemAvatar>
                  <Avatar src={user.avatar_url || undefined} sx={{ width: 32, height: 32 }}>
                    {user.display_name?.charAt(0).toUpperCase()}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={user.display_name}
                  secondary={`@${user.username}`}
                  primaryTypographyProps={{ variant: 'body2', fontWeight: 600 }}
                  secondaryTypographyProps={{ variant: 'caption' }}
                />
              </ListItem>
            ))}
          </List>
        )}
      </Paper>

      {/* フッター */}
      <Box sx={{ px: 1 }}>
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
