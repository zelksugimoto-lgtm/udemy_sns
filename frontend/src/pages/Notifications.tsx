import React, { useEffect } from 'react';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  CircularProgress,
  Alert,
  IconButton,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
  Comment as CommentIcon,
  PersonAdd as PersonAddIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import Layout from '../components/common/Layout';
import * as notificationsApi from '../api/endpoints/notifications';
import type { components } from '../api/generated/schema';

type NotificationResponse = components['schemas']['response.NotificationResponse'];

const Notifications: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationsApi.getNotifications({ limit: 50, offset: 0 }),
  });

  const markAllAsReadMutation = useMutation({
    mutationFn: notificationsApi.markAllNotificationsAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });

  const markAsReadMutation = useMutation({
    mutationFn: notificationsApi.markNotificationAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });

  // 未読の通知がある場合、すべて既読にする
  useEffect(() => {
    const hasUnread = data?.notifications.some((n) => !n.is_read);
    if (hasUnread) {
      markAllAsReadMutation.mutate();
    }
  }, [data]);

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'like':
        return <FavoriteIcon color="error" />;
      case 'comment':
      case 'reply':
        return <CommentIcon color="primary" />;
      case 'follow':
        return <PersonAddIcon color="primary" />;
      default:
        return null;
    }
  };

  const getNotificationText = (notification: NotificationResponse) => {
    const actorDisplayName = notification.actor?.display_name || 'ユーザー';
    const actorUsername = notification.actor?.username;

    const handleActorClick = (e: React.MouseEvent) => {
      e.stopPropagation(); // 通知全体のクリックイベントを止める
      if (actorUsername) {
        navigate(`/users/${actorUsername}`);
      }
    };

    let actionText = '';
    switch (notification.type) {
      case 'like':
        actionText = 'があなたの投稿にいいねしました';
        break;
      case 'comment':
        actionText = 'があなたの投稿にコメントしました';
        break;
      case 'reply':
        actionText = 'があなたのコメントに返信しました';
        break;
      case 'follow':
        actionText = 'があなたをフォローしました';
        break;
      default:
        return <span>{notification.message || ''}</span>;
    }

    return (
      <span>
        <Typography
          component="span"
          sx={{
            fontWeight: 700,
            cursor: 'pointer',
            '&:hover': {
              textDecoration: 'underline',
            },
          }}
          onClick={handleActorClick}
        >
          {actorDisplayName}
        </Typography>
        {actionText}
      </span>
    );
  };

  const handleNotificationClick = (notification: NotificationResponse) => {
    if (!notification.is_read) {
      markAsReadMutation.mutate(notification.id);
    }

    // 通知タイプに応じて遷移
    // post_idがある場合は投稿詳細ページへ遷移（コメント・返信・いいね通知）
    if (notification.post_id) {
      navigate(`/posts/${notification.post_id}`);
    }
    // フォロー通知などpost_idがない場合はアクターのプロフィールへ遷移
    else if (notification.actor?.username) {
      navigate(`/users/${notification.actor.username}`);
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
          通知
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
          {error instanceof Error ? error.message : '通知の取得に失敗しました'}
        </Alert>
      )}

      {/* Notifications */}
      {data && data.notifications.length === 0 && (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="body1" color="text.secondary">
            通知はありません
          </Typography>
        </Box>
      )}

      {data && data.notifications.length > 0 && (
        <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
          {data.notifications.map((notification) => (
            <ListItem
              key={notification.id}
              sx={{
                cursor: 'pointer',
                backgroundColor: notification.is_read ? 'inherit' : 'action.hover',
                borderBottom: '1px solid',
                borderColor: 'divider',
                '&:hover': {
                  backgroundColor: 'action.selected',
                },
              }}
              onClick={() => handleNotificationClick(notification)}
            >
              <ListItemAvatar>
                <Box sx={{ position: 'relative' }}>
                  <Avatar
                    alt={notification.actor?.display_name}
                    src={notification.actor?.avatar_url || undefined}
                  >
                    {notification.actor?.display_name?.charAt(0).toUpperCase()}
                  </Avatar>
                  <Box
                    sx={{
                      position: 'absolute',
                      bottom: -4,
                      right: -4,
                      backgroundColor: 'background.paper',
                      borderRadius: '50%',
                      p: 0.5,
                    }}
                  >
                    {getNotificationIcon(notification.type)}
                  </Box>
                </Box>
              </ListItemAvatar>
              <ListItemText
                primary={getNotificationText(notification)}
                secondary={formatDate(notification.created_at)}
                primaryTypographyProps={{
                  fontWeight: notification.is_read ? 400 : 700,
                }}
              />
            </ListItem>
          ))}
        </List>
      )}
    </Layout>
  );
};

export default Notifications;
