import React, { useEffect, useRef, useCallback } from 'react';
import { Box, CircularProgress, Typography, Alert } from '@mui/material';
import { useInfiniteQuery } from '@tanstack/react-query';
import * as postsApi from '../../api/endpoints/posts';
import PostCard from '../post/PostCard';
import { PAGINATION } from '../../utils/constants';

const Timeline: React.FC = () => {
  const observerTarget = useRef<HTMLDivElement>(null);

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError,
    error,
  } = useInfiniteQuery({
    queryKey: ['timeline'],
    queryFn: ({ pageParam = 0 }) =>
      postsApi.getTimeline({ limit: PAGINATION.DEFAULT_LIMIT, offset: pageParam }),
    getNextPageParam: (lastPage, allPages) => {
      if (!lastPage.pagination.has_more) {
        return undefined;
      }
      return allPages.length * PAGINATION.DEFAULT_LIMIT;
    },
    initialPageParam: 0,
  });

  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const [target] = entries;
      if (target.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    },
    [fetchNextPage, hasNextPage, isFetchingNextPage]
  );

  useEffect(() => {
    const element = observerTarget.current;
    if (!element) return;

    const option = { threshold: 0 };
    const observer = new IntersectionObserver(handleObserver, option);
    observer.observe(element);

    return () => observer.unobserve(element);
  }, [handleObserver]);

  if (isLoading) {
    return (
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
    );
  }

  if (isError) {
    return (
      <Alert severity="error" sx={{ m: 2 }}>
        {error instanceof Error ? error.message : 'タイムラインの取得に失敗しました'}
      </Alert>
    );
  }

  const posts = data?.pages.flatMap((page) => page.posts) || [];

  if (posts.length === 0) {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="body1" color="text.secondary">
          まだ投稿がありません。
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          フォローしているユーザーの投稿がここに表示されます。
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}

      {/* Intersection Observer Target */}
      <div ref={observerTarget} style={{ height: '20px' }} />

      {/* Loading More */}
      {isFetchingNextPage && (
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            py: 3,
            gap: 1,
          }}
        >
          <CircularProgress size={32} />
          <Typography variant="body2" color="text.secondary">
            読み込み中...
          </Typography>
        </Box>
      )}

      {/* No More Posts */}
      {!hasNextPage && posts.length > 0 && (
        <Box sx={{ p: 2, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            すべての投稿を表示しました
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default Timeline;
