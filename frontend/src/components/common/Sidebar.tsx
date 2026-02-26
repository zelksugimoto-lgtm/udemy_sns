import React from 'react';
import {
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
} from '@mui/material';
import {
  Home as HomeIcon,
  Notifications as NotificationsIcon,
  Bookmark as BookmarkIcon,
  Person as PersonIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

interface SidebarProps {
  onNavigate?: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ onNavigate }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();

  const handleNavigation = (path: string) => {
    navigate(path);
    if (onNavigate) {
      onNavigate();
    }
  };

  const menuItems = [
    { text: 'ホーム', icon: <HomeIcon />, path: '/' },
    { text: '通知', icon: <NotificationsIcon />, path: '/notifications' },
    { text: 'ブックマーク', icon: <BookmarkIcon />, path: '/bookmarks' },
    {
      text: 'プロフィール',
      icon: <PersonIcon />,
      path: user?.username ? `/users/${user.username}` : '/profile',
    },
    { text: '設定', icon: <SettingsIcon />, path: '/settings' },
  ];

  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column', py: 2 }}>
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton
              onClick={() => handleNavigation(item.path)}
              selected={location.pathname === item.path}
              sx={{
                borderRadius: 2,
                mx: 1,
                '&.Mui-selected': {
                  backgroundColor: 'primary.light',
                  '&:hover': {
                    backgroundColor: 'primary.light',
                  },
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: location.pathname === item.path ? 'primary.main' : 'inherit',
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText
                primary={item.text}
                primaryTypographyProps={{
                  fontWeight: location.pathname === item.path ? 700 : 400,
                }}
              />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );
};

export default Sidebar;
