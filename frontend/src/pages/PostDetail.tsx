import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, CircularProgress, Alert, Typography, IconButton } from '@mui/material';
import { ArrowBack as ArrowBackIcon } from '@mui/icons-material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import Layout from '../components/common/Layout';
import PostCard from '../components/post/PostCard';
import CommentForm from '../components/comment/CommentForm';
import CommentCard from '../components/comment/CommentCard';
import * as postsApi from '../api/endpoints/posts';
import * as commentsApi from '../api/endpoints/comments';

const PostDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // 投稿取得
  const {
    data: post,
    isLoading: postLoading,
    isError: postError,
  } = useQuery({
    queryKey: ['post', id],
    queryFn: () => postsApi.getPost(id!),
    enabled: !!id,
  });

  // 投稿が削除された場合はホームに戻る
  useEffect(() => {
    if (postError) {
      navigate('/');
    }
  }, [postError, navigate]);

  // コメント取得
  const {
    data: commentsData,
    isLoading: commentsLoading,
    isError: commentsError,
  } = useQuery({
    queryKey: ['comments', id],
    queryFn: () => commentsApi.getComments(id!),
    enabled: !!id,
  });

  const comments = commentsData?.comments || [];

  if (postLoading || commentsLoading) {
    return (
      <Layout>
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
      </Layout>
    );
  }

  if (postError || !post) {
    return (
      <Layout>
        <Alert severity="error" sx={{ m: 2 }}>
          投稿が見つかりませんでした
        </Alert>
      </Layout>
    );
  }

  return (
    <Layout>
      {/* Header with Back Button */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 2,
          p: 2,
          borderBottom: '1px solid',
          borderColor: 'divider',
        }}
      >
        <IconButton onClick={() => navigate(-1)}>
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="h6" fontWeight={700}>
          投稿
        </Typography>
      </Box>

      {/* Post */}
      <PostCard post={post} showActions={true} isDetailView={true} />

      {/* Comment Form */}
      <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
        <CommentForm postId={id!} />
      </Box>

      {/* Comments */}
      <Box>
        {commentsError && (
          <Alert severity="error" sx={{ m: 2 }}>
            コメントの取得に失敗しました
          </Alert>
        )}

        {comments && comments.length === 0 && (
          <Box sx={{ p: 4, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              まだコメントがありません
            </Typography>
          </Box>
        )}

        {comments?.map((comment) => (
          <CommentCard key={comment.id} comment={comment} postId={id!} />
        ))}
      </Box>
    </Layout>
  );
};

export default PostDetail;
