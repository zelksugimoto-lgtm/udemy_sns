import React, { useState } from 'react';
import {
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  Divider,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import Layout from '../components/common/Layout';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import { updateProfileSchema } from '../utils/validation';
import type { UpdateProfileFormData } from '../utils/validation';
import * as usersApi from '../api/endpoints/users';
import { THEMES } from '../utils/constants';

const Settings: React.FC = () => {
  const { user } = useAuth();
  const { currentTheme, setTheme } = useTheme();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<UpdateProfileFormData>({
    resolver: yupResolver(updateProfileSchema),
    defaultValues: {
      display_name: user?.display_name || '',
      bio: user?.bio || '',
    },
  });

  const updateProfileMutation = useMutation({
    mutationFn: usersApi.updateProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] });
      setSuccess('プロフィールを更新しました');
      setError(null);
    },
    onError: (err: Error) => {
      setError(err.message || 'プロフィールの更新に失敗しました');
      setSuccess(null);
    },
  });

  const onSubmit = async (data: UpdateProfileFormData) => {
    try {
      setError(null);
      setSuccess(null);
      await updateProfileMutation.mutateAsync(data);
    } catch (err) {
      // エラーは mutation の onError で処理
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
          設定
        </Typography>
      </Box>

      <Box sx={{ p: 2 }}>
        {/* Profile Settings */}
        <Paper elevation={0} sx={{ p: 3, mb: 3, border: '1px solid', borderColor: 'divider' }}>
          <Typography variant="h6" gutterBottom>
            プロフィール設定
          </Typography>
          <Divider sx={{ mb: 2 }} />

          {error && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
              {success}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit(onSubmit)}>
            <TextField
              fullWidth
              label="表示名"
              {...register('display_name')}
              error={!!errors.display_name}
              helperText={errors.display_name?.message}
              sx={{ mb: 2 }}
              inputProps={{ 'data-testid': 'settings-displayname' }}
            />

            <TextField
              fullWidth
              label="自己紹介"
              multiline
              minRows={4}
              {...register('bio')}
              error={!!errors.bio}
              helperText={errors.bio?.message}
              sx={{ mb: 2 }}
              inputProps={{ 'data-testid': 'settings-bio' }}
            />

            <Button
              type="submit"
              variant="contained"
              disabled={updateProfileMutation.isPending}
              sx={{
                borderRadius: 25,
                textTransform: 'none',
                fontWeight: 700,
                px: 3,
              }}
              data-testid="settings-submit"
            >
              {updateProfileMutation.isPending ? (
                <CircularProgress size={20} color="inherit" />
              ) : (
                '保存'
              )}
            </Button>
          </Box>
        </Paper>

        {/* Theme Settings */}
        <Paper elevation={0} sx={{ p: 3, border: '1px solid', borderColor: 'divider' }}>
          <Typography variant="h6" gutterBottom>
            テーマ設定
          </Typography>
          <Divider sx={{ mb: 2 }} />

          <FormControl fullWidth>
            <InputLabel>テーマ</InputLabel>
            <Select
              value={currentTheme}
              label="テーマ"
              onChange={(e) => setTheme(e.target.value)}
            >
              <MenuItem value={THEMES.LIGHT}>ライト</MenuItem>
              <MenuItem value={THEMES.DARK}>ダーク</MenuItem>
            </Select>
          </FormControl>
        </Paper>
      </Box>
    </Layout>
  );
};

export default Settings;
