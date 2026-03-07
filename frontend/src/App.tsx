import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/common/ProtectedRoute';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import PostDetail from './pages/PostDetail';
import Bookmarks from './pages/Bookmarks';
import Notifications from './pages/Notifications';
import Settings from './pages/Settings';
import Profile from './pages/Profile';
import Follow from './pages/Follow';
import PendingApproval from './pages/PendingApproval';
import RequestPasswordReset from './pages/RequestPasswordReset';
import ResetRequestConfirmation from './pages/ResetRequestConfirmation';
import ResetPassword from './pages/ResetPassword';

// React Query クライアント作成
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5分
    },
  },
});

const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <Router>
            <Routes>
              {/* 認証不要なルート */}
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/pending-approval" element={<PendingApproval />} />
              <Route path="/password-reset/request" element={<RequestPasswordReset />} />
              <Route path="/password-reset/confirmation" element={<ResetRequestConfirmation />} />
              <Route path="/password-reset/reset" element={<ResetPassword />} />

              {/* 認証が必要なルート */}
              <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
              <Route path="/posts/:id" element={<ProtectedRoute><PostDetail /></ProtectedRoute>} />
              <Route path="/users/:username" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
              <Route path="/users/:username/:tab" element={<ProtectedRoute><Follow /></ProtectedRoute>} />
              <Route path="/bookmarks" element={<ProtectedRoute><Bookmarks /></ProtectedRoute>} />
              <Route path="/notifications" element={<ProtectedRoute><Notifications /></ProtectedRoute>} />
              <Route path="/settings" element={<ProtectedRoute><Settings /></ProtectedRoute>} />
            </Routes>
          </Router>
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
};

export default App;
