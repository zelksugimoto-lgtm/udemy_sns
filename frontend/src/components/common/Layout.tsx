import React from 'react';
import { Box, useMediaQuery, useTheme } from '@mui/material';
import Header from './Header';
import Sidebar from './Sidebar';
import RightSidebar from './RightSidebar';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Header />
      <Box sx={{ display: 'flex', flex: 1, pt: { xs: 7, md: 8 } }}>
        {/* Sidebar - PC only */}
        {!isMobile && (
          <Box
            component="aside"
            sx={{
              width: 240,
              flexShrink: 0,
              borderRight: 1,
              borderColor: 'divider',
              position: 'fixed',
              height: 'calc(100vh - 64px)',
              top: 64,
              left: 0,
              overflowY: 'auto',
            }}
          >
            <Sidebar />
          </Box>
        )}

        {/* Main Content */}
        <Box
          component="main"
          sx={{
            flex: 1,
            ml: { xs: 0, md: '240px' },
            maxWidth: { xs: '100%', md: '600px', lg: '600px' },
            borderRight: { xs: 0, md: 1 },
            borderColor: 'divider',
            minHeight: 'calc(100vh - 64px)',
          }}
        >
          {children}
        </Box>

        {/* Right Sidebar - PC only (lg以上) */}
        {!isMobile && (
          <Box
            component="aside"
            sx={{
              display: { xs: 'none', lg: 'block' },
              width: 350,
              flexShrink: 0,
            }}
          >
            <RightSidebar />
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default Layout;
