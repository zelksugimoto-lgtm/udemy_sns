import React, { useState } from 'react';
import {
  Card,
  CardContent,
  CardActions,
  Avatar,
  Typography,
  IconButton,
  Box,
  Menu,
  MenuItem,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  Comment as CommentIcon,
  Bookmark as BookmarkIcon,
  BookmarkBorder as BookmarkBorderIcon,
  Share as ShareIcon,
  MoreVert as MoreVertIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { formatDistanceToNow } from 'date-fns';
import { useAuth } from '../../contexts/AuthContext';
import * as postsApi from '../../api/endpoints/posts';
import type { components } from '../../api/generated/schema';

type PostResponse = components['schemas']['response.PostResponse'];

interface PostCardProps {
  post: PostResponse;
  showActions?: boolean;
  isDetailView?: boolean;
}

const PostCard: React.FC<PostCardProps> = ({ post, showActions = true, isDetailView = false }) => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  // ローカル状態で管理（楽観的更新用）
  const [isLiked, setIsLiked] = useState(post.is_liked || false);
  const [isBookmarked, setIsBookmarked] = useState(post.is_bookmarked || false);
  const [likesCount, setLikesCount] = useState(post.likes_count || 0);

  const isOwnPost = user?.id === post.user?.id;

  // propsが更新されたらローカル状態も更新
  React.useEffect(() => {
    setIsLiked(post.is_liked || false);
    setIsBookmarked(post.is_bookmarked || false);
    setLikesCount(post.likes_count || 0);
  }, [post.id, post.is_liked, post.is_bookmarked, post.likes_count]);

  // いいね機能（楽観的更新）
  const likeMutation = useMutation({
    mutationFn: () => postsApi.likePost(post.id),
    onMutate: async () => {
      // ローカル状態を即座に更新
      setIsLiked(true);
      setLikesCount((prev) => prev + 1);

      await queryClient.cancelQueries({ queryKey: ['timeline'] });
      await queryClient.cancelQueries({ queryKey: ['post', post.id] });
      await queryClient.cancelQueries({ queryKey: ['userPosts'] });

      // タイムライン更新
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old) return old;
        return {
          ...old,
          pages: old.pages.map((page: any) => ({
            ...page,
            posts: page.posts.map((p: any) =>
              p.id === post.id
                ? { ...p, is_liked: true, likes_count: (p.likes_count || 0) + 1 }
                : p
            ),
          })),
        };
      });

      // 投稿詳細更新
      queryClient.setQueryData(['post', post.id], (old: any) => {
        if (!old) return old;
        return { ...old, is_liked: true, likes_count: (old.likes_count || 0) + 1 };
      });

      // ユーザー投稿一覧更新（全てのuserPostsクエリを更新）
      const userPostsQueries = queryClient.getQueriesData({ queryKey: ['userPosts'] });
      userPostsQueries.forEach(([queryKey, oldData]) => {
        if (oldData) {
          queryClient.setQueryData(queryKey, {
            ...oldData,
            posts: (oldData as any).posts?.map((p: any) =>
              p.id === post.id
                ? { ...p, is_liked: true, likes_count: (p.likes_count || 0) + 1 }
                : p
            ),
          });
        }
      });
    },
    onSuccess: () => {
      // 成功時はサーバーから最新データを取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
    onError: () => {
      // エラー時は元に戻す
      setIsLiked(post.is_liked || false);
      setLikesCount(post.likes_count || 0);

      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
  });

  const unlikeMutation = useMutation({
    mutationFn: () => postsApi.unlikePost(post.id),
    onMutate: async () => {
      // ローカル状態を即座に更新
      setIsLiked(false);
      setLikesCount((prev) => Math.max(prev - 1, 0));

      await queryClient.cancelQueries({ queryKey: ['timeline'] });
      await queryClient.cancelQueries({ queryKey: ['post', post.id] });
      await queryClient.cancelQueries({ queryKey: ['userPosts'] });

      // タイムライン更新
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old) return old;
        return {
          ...old,
          pages: old.pages.map((page: any) => ({
            ...page,
            posts: page.posts.map((p: any) =>
              p.id === post.id
                ? { ...p, is_liked: false, likes_count: Math.max((p.likes_count || 0) - 1, 0) }
                : p
            ),
          })),
        };
      });

      // 投稿詳細更新
      queryClient.setQueryData(['post', post.id], (old: any) => {
        if (!old) return old;
        return { ...old, is_liked: false, likes_count: Math.max((old.likes_count || 0) - 1, 0) };
      });

      // ユーザー投稿一覧更新（全てのuserPostsクエリを更新）
      const userPostsQueries2 = queryClient.getQueriesData({ queryKey: ['userPosts'] });
      userPostsQueries2.forEach(([queryKey, oldData]) => {
        if (oldData) {
          queryClient.setQueryData(queryKey, {
            ...oldData,
            posts: (oldData as any).posts?.map((p: any) =>
              p.id === post.id
                ? { ...p, is_liked: false, likes_count: Math.max((p.likes_count || 0) - 1, 0) }
                : p
            ),
          });
        }
      });
    },
    onSuccess: () => {
      // 成功時はサーバーから最新データを取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
    onError: () => {
      // エラー時は元に戻す
      setIsLiked(post.is_liked || false);
      setLikesCount(post.likes_count || 0);

      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
  });

  // ブックマーク機能（楽観的更新）
  const bookmarkMutation = useMutation({
    mutationFn: () => postsApi.bookmarkPost(post.id),
    onMutate: async () => {
      // ローカル状態を即座に更新
      setIsBookmarked(true);

      await queryClient.cancelQueries({ queryKey: ['timeline'] });
      await queryClient.cancelQueries({ queryKey: ['post', post.id] });
      await queryClient.cancelQueries({ queryKey: ['userPosts'] });

      // タイムライン更新
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old) return old;
        return {
          ...old,
          pages: old.pages.map((page: any) => ({
            ...page,
            posts: page.posts.map((p: any) =>
              p.id === post.id ? { ...p, is_bookmarked: true } : p
            ),
          })),
        };
      });

      // 投稿詳細更新
      queryClient.setQueryData(['post', post.id], (old: any) => {
        if (!old) return old;
        return { ...old, is_bookmarked: true };
      });

      // ユーザー投稿一覧更新（全てのuserPostsクエリを更新）
      const userPostsQueries3 = queryClient.getQueriesData({ queryKey: ['userPosts'] });
      userPostsQueries3.forEach(([queryKey, oldData]) => {
        if (oldData) {
          queryClient.setQueryData(queryKey, {
            ...oldData,
            posts: (oldData as any).posts?.map((p: any) =>
              p.id === post.id ? { ...p, is_bookmarked: true } : p
            ),
          });
        }
      });
    },
    onSuccess: () => {
      // 成功時はサーバーから最新データを取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
    onError: () => {
      // エラー時は元に戻す
      setIsBookmarked(post.is_bookmarked || false);

      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
  });

  const unbookmarkMutation = useMutation({
    mutationFn: () => postsApi.unbookmarkPost(post.id),
    onMutate: async () => {
      // ローカル状態を即座に更新
      setIsBookmarked(false);

      await queryClient.cancelQueries({ queryKey: ['timeline'] });
      await queryClient.cancelQueries({ queryKey: ['post', post.id] });
      await queryClient.cancelQueries({ queryKey: ['userPosts'] });

      // タイムライン更新
      queryClient.setQueriesData({ queryKey: ['timeline'] }, (old: any) => {
        if (!old) return old;
        return {
          ...old,
          pages: old.pages.map((page: any) => ({
            ...page,
            posts: page.posts.map((p: any) =>
              p.id === post.id ? { ...p, is_bookmarked: false } : p
            ),
          })),
        };
      });

      // 投稿詳細更新
      queryClient.setQueryData(['post', post.id], (old: any) => {
        if (!old) return old;
        return { ...old, is_bookmarked: false };
      });

      // ユーザー投稿一覧更新（全てのuserPostsクエリを更新）
      const userPostsQueries4 = queryClient.getQueriesData({ queryKey: ['userPosts'] });
      userPostsQueries4.forEach(([queryKey, oldData]) => {
        if (oldData) {
          queryClient.setQueryData(queryKey, {
            ...oldData,
            posts: (oldData as any).posts?.map((p: any) =>
              p.id === post.id ? { ...p, is_bookmarked: false } : p
            ),
          });
        }
      });
    },
    onSuccess: () => {
      // 成功時はサーバーから最新データを取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
    onError: () => {
      // エラー時は元に戻す
      setIsBookmarked(post.is_bookmarked || false);

      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post', post.id] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['userLikedPosts'] });
      queryClient.invalidateQueries({ queryKey: ['bookmarks'] });
    },
  });

  // 削除機能
  const deleteMutation = useMutation({
    mutationFn: () => postsApi.deletePost(post.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
    },
  });

  const handleLikeToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (isLiked) {
      unlikeMutation.mutate();
    } else {
      likeMutation.mutate();
    }
  };

  const handleBookmarkToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (isBookmarked) {
      unbookmarkMutation.mutate();
    } else {
      bookmarkMutation.mutate();
    }
  };

  const handleMenuOpen = (e: React.MouseEvent<HTMLElement>) => {
    e.stopPropagation();
    setAnchorEl(e.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleDelete = async (e: React.MouseEvent) => {
    e.stopPropagation();
    handleMenuClose();
    if (window.confirm('この投稿を削除しますか？')) {
      await deleteMutation.mutateAsync();
    }
  };

  const handlePostClick = () => {
    if (!isDetailView) {
      navigate(`/posts/${post.id}`);
    }
  };

  const handleUserClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (post.user?.username) {
      navigate(`/users/${post.user.username}`);
    }
  };

  const handleShare = async (e: React.MouseEvent) => {
    e.stopPropagation();
    const url = `${window.location.origin}/posts/${post.id}`;
    try {
      await navigator.clipboard.writeText(url);
      alert('投稿のURLをコピーしました');
    } catch (error) {
      console.error('URLのコピーに失敗しました:', error);
    }
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return formatDistanceToNow(date, { addSuffix: true });
    } catch {
      return '';
    }
  };

  return (
    <Card
      elevation={0}
      sx={{
        borderBottom: '1px solid',
        borderColor: 'divider',
        borderRadius: 0,
        cursor: isDetailView ? 'default' : 'pointer',
        '&:hover': {
          backgroundColor: isDetailView ? 'transparent' : 'action.hover',
        },
      }}
      onClick={handlePostClick}
      data-testid="post-card"
    >
      <CardContent sx={{ p: 2, pb: 1 }}>
        <Box sx={{ display: 'flex', gap: 2 }}>
          {/* Avatar */}
          <Avatar
            alt={post.user?.display_name}
            src={post.user?.avatar_url || undefined}
            sx={{ width: 48, height: 48, cursor: 'pointer' }}
            onClick={handleUserClick}
          >
            {post.user?.display_name?.charAt(0).toUpperCase()}
          </Avatar>

          <Box sx={{ flex: 1, minWidth: 0 }}>
            {/* User Info */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
              <Typography
                variant="subtitle2"
                sx={{ fontWeight: 700, cursor: 'pointer' }}
                onClick={handleUserClick}
              >
                {post.user?.display_name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                @{post.user?.username}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                · {formatDate(post.created_at)}
              </Typography>

              <Box sx={{ ml: 'auto' }}>
                <IconButton size="small" onClick={handleMenuOpen} data-testid="post-menu-button">
                  <MoreVertIcon fontSize="small" />
                </IconButton>
                <Menu
                  anchorEl={anchorEl}
                  open={Boolean(anchorEl)}
                  onClose={handleMenuClose}
                >
                  {isOwnPost ? (
                    <MenuItem onClick={handleDelete} data-testid="post-delete-button">削除</MenuItem>
                  ) : (
                    <MenuItem onClick={handleMenuClose}>通報</MenuItem>
                  )}
                </Menu>
              </Box>
            </Box>

            {/* Content */}
            <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap', mb: 1 }} data-testid="post-content">
              {post.content}
            </Typography>
          </Box>
        </Box>
      </CardContent>

      {showActions && (
        <CardActions sx={{ px: 2, pb: 1, justifyContent: 'space-around' }}>
          {/* Comment */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <IconButton size="small" color="default" data-testid="post-comment-button">
              <CommentIcon fontSize="small" />
            </IconButton>
            <Typography variant="caption" color="text.secondary">
              {post.comments_count || 0}
            </Typography>
          </Box>

          {/* Like */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <IconButton
              size="small"
              color={isLiked ? 'error' : 'default'}
              onClick={handleLikeToggle}
              data-testid="post-like-button"
            >
              {isLiked ? <FavoriteIcon fontSize="small" /> : <FavoriteBorderIcon fontSize="small" />}
            </IconButton>
            <Typography variant="caption" color="text.secondary">
              {likesCount}
            </Typography>
          </Box>

          {/* Bookmark */}
          <IconButton
            size="small"
            color={isBookmarked ? 'primary' : 'default'}
            onClick={handleBookmarkToggle}
          >
            {isBookmarked ? <BookmarkIcon fontSize="small" /> : <BookmarkBorderIcon fontSize="small" />}
          </IconButton>

          {/* Share */}
          <IconButton size="small" color="default" onClick={handleShare}>
            <ShareIcon fontSize="small" />
          </IconButton>
        </CardActions>
      )}
    </Card>
  );
};

export default PostCard;
