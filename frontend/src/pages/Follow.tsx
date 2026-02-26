import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Tabs,
  Tab,
  List,
  ListItem,
  ListItemAvatar,
  Avatar,
  ListItemText,
  Button,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/common/Layout';
import * as usersApi from '../api/endpoints/users';

const Follow: React.FC = () => {
  const { username, tab } = useParams<{ username: string; tab: 'followers' | 'following' }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user: currentUser, isAuthenticated, isLoading: authLoading } = useAuth();
  const [tabValue, setTabValue] = useState(tab === 'following' ? 1 : 0);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, authLoading, navigate]);

  // タブ変更時にURLも更新
  useEffect(() => {
    if (username) {
      const newTab = tabValue === 0 ? 'followers' : 'following';
      navigate(`/users/${username}/${newTab}`, { replace: true });
    }
  }, [tabValue, username, navigate]);

  // フォロワー一覧取得
  const { data: followersData, isLoading: followersLoading, isError: followersError } = useQuery({
    queryKey: ['followers', username],
    queryFn: () => usersApi.getFollowers(username!, { limit: 50, offset: 0 }),
    enabled: !!username && tabValue === 0,
  });

  // フォロー中一覧取得
  const { data: followingData, isLoading: followingLoading, isError: followingError } = useQuery({
    queryKey: ['following', username],
    queryFn: () => usersApi.getFollowing(username!, { limit: 50, offset: 0 }),
    enabled: !!username && tabValue === 1,
  });

  // フォロー/フォロー解除ミューテーション
  const followMutation = useMutation({
    mutationFn: (targetUsername: string) => usersApi.followUser(targetUsername),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['followers', username] });
      queryClient.invalidateQueries({ queryKey: ['following', username] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });

  const unfollowMutation = useMutation({
    mutationFn: (targetUsername: string) => usersApi.unfollowUser(targetUsername),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['followers', username] });
      queryClient.invalidateQueries({ queryKey: ['following', username] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });

  const handleFollowToggle = (targetUsername: string, isFollowing: boolean) => {
    if (isFollowing) {
      unfollowMutation.mutate(targetUsername);
    } else {
      followMutation.mutate(targetUsername);
    }
  };

  const handleUserClick = (clickedUsername: string) => {
    navigate(`/users/${clickedUsername}`);
  };

  if (authLoading) {
    return (
      <Layout>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 200 }}>
          <CircularProgress />
        </Box>
      </Layout>
    );
  }

  const followers = followersData?.users || [];
  const following = followingData?.users || [];
  const isLoading = tabValue === 0 ? followersLoading : followingLoading;
  const isError = tabValue === 0 ? followersError : followingError;
  const users = tabValue === 0 ? followers : following;

  return (
    <Layout>
      <Box sx={{ borderBottom: '1px solid', borderColor: 'divider' }}>
        <Box sx={{ p: 2 }}>
          <Typography variant="h6" sx={{ fontWeight: 700 }}>
            @{username}
          </Typography>
        </Box>
        <Tabs value={tabValue} onChange={(_, newValue) => setTabValue(newValue)}>
          <Tab label="フォロワー" />
          <Tab label="フォロー中" />
        </Tabs>
      </Box>

      {isError && (
        <Alert severity="error" sx={{ m: 2 }}>
          データの取得に失敗しました
        </Alert>
      )}

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : users.length === 0 ? (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="body1" color="text.secondary">
            {tabValue === 0 ? 'フォロワーがいません' : 'フォロー中のユーザーがいません'}
          </Typography>
        </Box>
      ) : (
        <List>
          {users.map((user) => {
            const isOwnProfile = currentUser?.username === user.username;
            return (
              <ListItem
                key={user.id}
                sx={{
                  borderBottom: '1px solid',
                  borderColor: 'divider',
                  py: 2,
                  px: 2,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 2,
                }}
              >
                <ListItemAvatar onClick={() => handleUserClick(user.username)} sx={{ cursor: 'pointer' }}>
                  <Avatar src={user.avatar_url || undefined} sx={{ width: 48, height: 48 }}>
                    {user.display_name?.charAt(0).toUpperCase()}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={
                    <Typography
                      variant="subtitle2"
                      sx={{ fontWeight: 700, cursor: 'pointer' }}
                      onClick={() => handleUserClick(user.username)}
                    >
                      {user.display_name}
                    </Typography>
                  }
                  secondary={
                    <>
                      <Typography variant="body2" color="text.secondary">
                        @{user.username}
                      </Typography>
                      {user.bio && (
                        <Typography variant="body2" sx={{ mt: 0.5 }}>
                          {user.bio}
                        </Typography>
                      )}
                    </>
                  }
                  sx={{ flex: 1 }}
                />
                {!isOwnProfile && (
                  <Button
                    variant={user.is_following ? 'outlined' : 'contained'}
                    size="small"
                    onClick={() => handleFollowToggle(user.username, user.is_following || false)}
                    disabled={followMutation.isPending || unfollowMutation.isPending}
                  >
                    {user.is_following ? 'フォロー解除' : 'フォロー'}
                  </Button>
                )}
              </ListItem>
            );
          })}
        </List>
      )}
    </Layout>
  );
};

export default Follow;
