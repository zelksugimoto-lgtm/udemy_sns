import React from 'react';
import { Box, Typography, CircularProgress, Alert } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import Layout from '../components/common/Layout';
import PostCard from '../components/post/PostCard';
import * as bookmarksApi from '../api/endpoints/bookmarks';

const Bookmarks: React.FC = () => {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['bookmarks'],
    queryFn: () => bookmarksApi.getBookmarks({ limit: 100, offset: 0 }),
  });

  return (
    <Layout>
      {/* Header */}
      <Box
        sx={{
          p: 2,
          borderBottom: '1px solid',
          borderColor: 'divider',
        }}
      >
        <Typography variant="h6" fontWeight={700}>
          ブックマーク
        </Typography>
      </Box>

      {/* Loading */}
      {isLoading && (
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: 200,
          }}
        >
          <CircularProgress />
        </Box>
      )}

      {/* Error */}
      {isError && (
        <Alert severity="error" sx={{ m: 2 }}>
          {error instanceof Error ? error.message : 'ブックマークの取得に失敗しました'}
        </Alert>
      )}

      {/* Posts */}
      {data && data.posts.length === 0 && (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="body1" color="text.secondary">
            まだブックマークがありません
          </Typography>
        </Box>
      )}

      {data?.posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </Layout>
  );
};

export default Bookmarks;
