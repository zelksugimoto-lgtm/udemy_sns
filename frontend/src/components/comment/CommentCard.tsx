import React, { useState } from 'react';
import {
  Box,
  Avatar,
  Typography,
  IconButton,
  Menu,
  MenuItem,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  Reply as ReplyIcon,
  MoreVert as MoreVertIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { formatDistanceToNow } from 'date-fns';
import { useAuth } from '../../contexts/AuthContext';
import * as commentsApi from '../../api/endpoints/comments';
import CommentForm from './CommentForm';
import type { components } from '../../api/generated/schema';

type CommentResponse = components['schemas']['response.CommentResponse'];

interface CommentCardProps {
  comment: CommentResponse;
  postId: string;
  level?: number;
}

const CommentCard: React.FC<CommentCardProps> = ({ comment, postId, level = 0 }) => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [showReplyForm, setShowReplyForm] = useState(false);

  const isLiked = comment.is_liked || false;
  const isOwnComment = user?.id === comment.user?.id;

  // いいね機能（楽観的更新）
  const likeMutation = useMutation({
    mutationFn: () => commentsApi.likeComment(comment.id),
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ['comments', postId] });

      queryClient.setQueryData(['comments', postId], (old: any) => {
        if (!old) return old;
        return {
          ...old,
          comments: old.comments?.map((c: any) =>
            c.id === comment.id
              ? { ...c, is_liked: true, like_count: (c.like_count || 0) + 1 }
              : c
          ),
        };
      });
    },
    onError: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postId] });
    },
  });

  const unlikeMutation = useMutation({
    mutationFn: () => commentsApi.unlikeComment(comment.id),
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ['comments', postId] });

      queryClient.setQueryData(['comments', postId], (old: any) => {
        if (!old) return old;
        return {
          ...old,
          comments: old.comments?.map((c: any) =>
            c.id === comment.id
              ? { ...c, is_liked: false, like_count: Math.max((c.like_count || 0) - 1, 0) }
              : c
          ),
        };
      });
    },
    onError: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postId] });
    },
  });

  // 削除機能
  const deleteMutation = useMutation({
    mutationFn: () => commentsApi.deleteComment(comment.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postId] });
      queryClient.invalidateQueries({ queryKey: ['post', postId] });
    },
  });

  const handleLikeToggle = () => {
    if (isLiked) {
      unlikeMutation.mutate();
    } else {
      likeMutation.mutate();
    }
  };

  const handleMenuOpen = (e: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(e.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleDelete = () => {
    if (window.confirm('このコメントを削除しますか？')) {
      deleteMutation.mutate();
    }
    handleMenuClose();
  };

  const handleUserClick = () => {
    if (comment.user?.username) {
      navigate(`/users/${comment.user.username}`);
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
    <Box
      sx={{
        p: 2,
        pb: 1,
        pl: level > 0 ? 2 + level * 4 : 2,
        borderLeft: level > 0 ? 1 : 0,
        borderColor: 'divider',
        borderBottom: 1,
      }}
    >
      <Box sx={{ display: 'flex', gap: 2 }}>
        {/* Avatar */}
        <Avatar
          alt={comment.user?.display_name}
          src={comment.user?.avatar_url || undefined}
          sx={{ width: 48, height: 48, cursor: 'pointer' }}
          onClick={handleUserClick}
        >
          {comment.user?.display_name?.charAt(0).toUpperCase()}
        </Avatar>

        <Box sx={{ flex: 1, minWidth: 0 }}>
          {/* User Info */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
            <Typography
              variant="subtitle2"
              sx={{ fontWeight: 700, cursor: 'pointer' }}
              onClick={handleUserClick}
            >
              {comment.user?.display_name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              @{comment.user?.username}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              · {formatDate(comment.created_at)}
            </Typography>

            <Box sx={{ ml: 'auto' }}>
              <IconButton size="small" onClick={handleMenuOpen}>
                <MoreVertIcon fontSize="small" />
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
              >
                {isOwnComment ? (
                  <MenuItem onClick={handleDelete}>削除</MenuItem>
                ) : (
                  <MenuItem onClick={handleMenuClose}>通報</MenuItem>
                )}
              </Menu>
            </Box>
          </Box>

          {/* Content */}
          <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mb: 1 }}>
            {comment.content}
          </Typography>

          {/* Actions */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            {/* Like */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <IconButton
                size="small"
                color={isLiked ? 'error' : 'default'}
                onClick={handleLikeToggle}
              >
                {isLiked ? (
                  <FavoriteIcon fontSize="small" />
                ) : (
                  <FavoriteBorderIcon fontSize="small" />
                )}
              </IconButton>
              <Typography variant="caption" color="text.secondary">
                {comment.like_count || 0}
              </Typography>
            </Box>

            {/* Reply */}
            {level < 3 && (
              <IconButton
                size="small"
                onClick={() => setShowReplyForm(!showReplyForm)}
              >
                <ReplyIcon fontSize="small" />
              </IconButton>
            )}
          </Box>

          {/* Reply Form */}
          {showReplyForm && (
            <Box sx={{ mt: 2 }}>
              <CommentForm
                postId={postId}
                parentCommentId={comment.id}
                onSuccess={() => setShowReplyForm(false)}
              />
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default CommentCard;
