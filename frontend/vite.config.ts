import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // テストモードの場合はポート3001を使用
  const port = mode === 'test' ? 3001 : 3000;

  return {
    plugins: [react()],
    server: {
      port,
      host: true,
    },
  };
});
