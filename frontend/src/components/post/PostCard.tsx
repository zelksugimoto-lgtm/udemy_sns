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
  Dialog,
  DialogContent,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  Comment as CommentIcon,
  Bookmark as BookmarkIcon,
  BookmarkBorder as BookmarkBorderIcon,
  Share as ShareIcon,
  MoreVert as MoreVertIcon,
  NavigateBefore as NavigateBeforeIcon,
  NavigateNext as NavigateNextIcon,
  Close as CloseIcon,
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

  // 画像モーダル用の状態
  const [imageModalOpen, setImageModalOpen] = useState(false);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

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

  // 画像モーダルを開く
  const handleImageClick = (e: React.MouseEvent, index: number) => {
    e.stopPropagation();
    setSelectedImageIndex(index);
    setImageModalOpen(true);
  };

  // 画像モーダルを閉じる
  const handleCloseImageModal = () => {
    setImageModalOpen(false);
  };

  // 前の画像へ
  const handlePrevImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedImageIndex((prev) => (prev > 0 ? prev - 1 : (post.media?.length || 1) - 1));
  };

  // 次の画像へ
  const handleNextImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedImageIndex((prev) => (prev < (post.media?.length || 1) - 1 ? prev + 1 : 0));
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

            {/* Media Grid - Twitter風レイアウト */}
            {post.media && post.media.length > 0 && (
              <Box
                sx={{
                  mt: 2,
                  display: 'grid',
                  gap: 0.5,
                  // レイアウト: 1枚=1列、2枚=2列、3枚=2列、4枚=2列
                  gridTemplateColumns: post.media.length === 1
                    ? '1fr'
                    : 'repeat(2, 1fr)',
                  // 高さ: 1枚=auto（アスペクト比維持）、2枚以上=固定
                  gridTemplateRows: post.media.length === 1
                    ? 'auto'
                    : post.media.length === 2
                    ? '280px'
                    : post.media.length === 3
                    ? 'repeat(2, 190px)'
                    : 'repeat(2, 190px)',
                }}
              >
                {post.media.map((media, index) => (
                  <Box
                    key={media.id}
                    sx={{
                      position: 'relative',
                      overflow: 'hidden',
                      borderRadius: 2,
                      backgroundColor: 'grey.100',
                      // 1枚の場合: max-heightで見切れを防ぐ
                      ...(post.media!.length === 1 && {
                        maxHeight: '500px',
                      }),
                      // 3枚の場合: 最初の画像を2行にまたがらせる
                      ...(post.media!.length === 3 && index === 0 && {
                        gridRow: 'span 2',
                      }),
                    }}
                  >
                    <img
                      src={media.media_url}
                      alt={`Media ${index + 1}`}
                      style={{
                        width: '100%',
                        height: '100%',
                        // 1枚の場合: contain（見切れ防止）、複数の場合: cover
                        objectFit: post.media!.length === 1 ? 'contain' : 'cover',
                        cursor: 'pointer',
                        display: 'block',
                      }}
                      onClick={(e) => handleImageClick(e, index)}
                    />
                  </Box>
                ))}
              </Box>
            )}
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

      {/* 画像モーダル */}
      <Dialog
        open={imageModalOpen}
        onClose={handleCloseImageModal}
        maxWidth="lg"
        fullWidth
        PaperProps={{
          sx: {
            backgroundColor: 'rgba(0, 0, 0, 0.9)',
            boxShadow: 'none',
          },
        }}
      >
        <DialogContent
          sx={{
            position: 'relative',
            p: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '80vh',
          }}
        >
          {/* 閉じるボタン */}
          <IconButton
            onClick={handleCloseImageModal}
            sx={{
              position: 'absolute',
              top: 16,
              right: 16,
              color: 'white',
              backgroundColor: 'rgba(255, 255, 255, 0.1)',
              '&:hover': {
                backgroundColor: 'rgba(255, 255, 255, 0.2)',
              },
              zIndex: 2,
            }}
          >
            <CloseIcon />
          </IconButton>

          {/* 前の画像ボタン（複数画像の場合のみ） */}
          {post.media && post.media.length > 1 && (
            <IconButton
              onClick={handlePrevImage}
              sx={{
                position: 'absolute',
                left: 16,
                color: 'white',
                backgroundColor: 'rgba(255, 255, 255, 0.1)',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.2)',
                },
                zIndex: 2,
              }}
            >
              <NavigateBeforeIcon fontSize="large" />
            </IconButton>
          )}

          {/* 画像表示 */}
          {post.media && post.media[selectedImageIndex] && (
            <Box
              component="img"
              src={post.media[selectedImageIndex].media_url}
              alt={`Media ${selectedImageIndex + 1}`}
              sx={{
                maxWidth: '100%',
                maxHeight: '80vh',
                objectFit: 'contain',
              }}
            />
          )}

          {/* 次の画像ボタン（複数画像の場合のみ） */}
          {post.media && post.media.length > 1 && (
            <IconButton
              onClick={handleNextImage}
              sx={{
                position: 'absolute',
                right: 16,
                color: 'white',
                backgroundColor: 'rgba(255, 255, 255, 0.1)',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.2)',
                },
                zIndex: 2,
              }}
            >
              <NavigateNextIcon fontSize="large" />
            </IconButton>
          )}

          {/* 画像カウンター（複数画像の場合のみ） */}
          {post.media && post.media.length > 1 && (
            <Typography
              sx={{
                position: 'absolute',
                bottom: 16,
                left: '50%',
                transform: 'translateX(-50%)',
                color: 'white',
                backgroundColor: 'rgba(0, 0, 0, 0.5)',
                px: 2,
                py: 0.5,
                borderRadius: 1,
                zIndex: 2,
              }}
            >
              {selectedImageIndex + 1} / {post.media.length}
            </Typography>
          )}
        </DialogContent>
      </Dialog>
    </Card>
  );
};

export default PostCard;
