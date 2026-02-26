import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Avatar,
  Typography,
  Button,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/common/Layout';
import PostCard from '../components/post/PostCard';
import * as usersApi from '../api/endpoints/users';
import * as postsApi from '../api/endpoints/posts';
import * as bookmarksApi from '../api/endpoints/bookmarks';

const Profile: React.FC = () => {
  const { username } = useParams<{ username: string }>();
  const navigate = useNavigate();
  const { user: currentUser, isAuthenticated, isLoading: authLoading } = useAuth();
  const [tabValue, setTabValue] = useState(0);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, authLoading, navigate]);

  const { data: profile, isLoading: profileLoading, isError: profileError } = useQuery({
    queryKey: ['profile', username],
    queryFn: () => usersApi.getUserProfile(username!),
    enabled: !!username,
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['userPosts', username],
    queryFn: () => postsApi.getUserPosts(username!, { limit: 20, offset: 0 }),
    enabled: !!username,
  });

  // プロフィール情報から自分のプロフィールかどうかを判定
  const isOwnProfile = currentUser?.username === username;

  // ブックマーク一覧取得（自分のプロフィールの場合のみ）
  const { data: bookmarksData, isLoading: bookmarksLoading } = useQuery({
    queryKey: ['bookmarks'],
    queryFn: () => bookmarksApi.getBookmarks({ limit: 20, offset: 0 }),
    enabled: isOwnProfile,
  });

  if (authLoading || profileLoading) {
    return (
      <Layout>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 200 }}>
          <CircularProgress />
        </Box>
      </Layout>
    );
  }

  if (profileError || !profile) {
    return (
      <Layout>
        <Alert severity="error" sx={{ m: 2 }}>
          ユーザーが見つかりませんでした
        </Alert>
      </Layout>
    );
  }

  const posts = postsData?.posts || [];
  const bookmarkedPosts = bookmarksData?.posts || [];

  return (
    <Layout>
      <Box sx={{ width: '100%', height: 200, bgcolor: 'grey.300', backgroundImage: profile.header_url ? `url(${profile.header_url})` : 'none', backgroundSize: 'cover', backgroundPosition: 'center' }} />
      <Box sx={{ px: 2, pt: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Avatar alt={profile.display_name} src={profile.avatar_url || undefined} sx={{ width: 120, height: 120, border: '4px solid white', mt: -8 }}>
            {profile.display_name?.charAt(0).toUpperCase()}
          </Avatar>
          {isOwnProfile ? (
            <Button variant="outlined" onClick={() => navigate('/settings')}>プロフィール編集</Button>
          ) : (
            <Button variant="contained">フォロー</Button>
          )}
        </Box>
        <Typography variant="h5" sx={{ fontWeight: 700, mb: 0.5 }}>{profile.display_name}</Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>@{profile.username}</Typography>
        {profile.bio && (<Typography variant="body1" sx={{ mb: 2, whiteSpace: 'pre-wrap' }}>{profile.bio}</Typography>)}
        <Box sx={{ display: 'flex', gap: 3, mb: 2 }}>
          <Typography variant="body2"><strong>{profile.posts_count || 0}</strong> 投稿</Typography>
          <Typography variant="body2"><strong>{profile.following_count || 0}</strong> フォロー中</Typography>
          <Typography variant="body2"><strong>{profile.followers_count || 0}</strong> フォロワー</Typography>
        </Box>
        <Tabs value={tabValue} onChange={(_, newValue) => setTabValue(newValue)}>
          <Tab label="投稿" />
          <Tab label="いいね" disabled />
          <Tab label="ブックマーク" disabled={!isOwnProfile} />
        </Tabs>
      </Box>
      <Box>
        {tabValue === 0 && (
          <>
            {postsLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 4 }}>
                <CircularProgress />
              </Box>
            ) : posts.length === 0 ? (
              <Box sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="body1" color="text.secondary">まだ投稿がありません</Typography>
              </Box>
            ) : (
              posts.map((post) => <PostCard key={post.id} post={post} />)
            )}
          </>
        )}
        {tabValue === 2 && isOwnProfile && (
          <>
            {bookmarksLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 4 }}>
                <CircularProgress />
              </Box>
            ) : bookmarkedPosts.length === 0 ? (
              <Box sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="body1" color="text.secondary">ブックマークがありません</Typography>
              </Box>
            ) : (
              bookmarkedPosts.map((post) => <PostCard key={post.id} post={post} />)
            )}
          </>
        )}
      </Box>
    </Layout>
  );
};

export default Profile;
