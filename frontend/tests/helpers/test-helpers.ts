import { Page } from '@playwright/test';

/**
 * テスト用ヘルパー関数
 * E2Eテストで共通的に使用される操作をまとめたユーティリティ関数
 */

/**
 * ランダムなユーザー情報を生成
 */
export function generateRandomUser() {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(7);
  return {
    email: `test_${timestamp}_${random}@example.com`,
    username: `test_${timestamp}_${random}`,
    displayName: `Test User ${timestamp}`,
    password: 'testpassword123',
  };
}

/**
 * ユーザー登録
 */
export async function registerUser(page: Page, user?: { email: string; username: string; displayName: string; password: string }) {
  const testUser = user || generateRandomUser();

  await page.goto('/register');
  await page.getByTestId('register-email').fill(testUser.email);
  await page.getByTestId('register-username').fill(testUser.username);
  await page.getByTestId('register-displayname').fill(testUser.displayName);
  await page.getByTestId('register-password').fill(testUser.password);
  await page.getByTestId('register-submit').click();

  // ホーム画面への遷移を待つ
  await page.waitForURL('/');

  return testUser;
}

/**
 * ログイン
 */
export async function loginUser(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.getByTestId('login-email').fill(email);
  await page.getByTestId('login-password').fill(password);
  await page.getByTestId('login-submit').click();

  // ホーム画面への遷移を待つ
  await page.waitForURL('/');
}

/**
 * ログアウト
 */
export async function logout(page: Page) {
  // ヘッダーのユーザーメニューアイコンをクリックしてメニューを開く
  await page.getByTestId('header-user-menu').click();

  // メニューが表示されるまで少し待つ
  await page.waitForTimeout(500);

  // ログアウトボタンをクリック
  await page.getByTestId('header-logout').click();

  // ログインページへの遷移を待つ
  await page.waitForURL('/login');
}

/**
 * 投稿作成
 */
export async function createPost(page: Page, content: string) {
  await page.goto('/');
  await page.getByTestId('post-form-content').fill(content);
  await page.getByTestId('post-form-submit').click();

  // 投稿が完了するまで待つ（投稿フォームがクリアされることを確認）
  await page.waitForTimeout(1000);
}

/**
 * コメント作成
 */
export async function createComment(page: Page, postId: string, content: string) {
  await page.goto(`/posts/${postId}`);
  await page.getByTestId('comment-form-content').fill(content);
  await page.getByTestId('comment-form-submit').click();

  // コメントが完了するまで待つ
  await page.waitForTimeout(1000);
}

/**
 * 投稿が表示されるまで待つ
 */
export async function waitForPostToAppear(page: Page, content: string) {
  await page.waitForSelector(`[data-testid="post-content"]:has-text("${content}")`);
}

/**
 * コメントが表示されるまで待つ
 */
export async function waitForCommentToAppear(page: Page, content: string) {
  await page.waitForSelector(`[data-testid="comment-content"]:has-text("${content}")`);
}

/**
 * ユーザープロフィールページへ移動
 */
export async function goToProfile(page: Page, username: string) {
  await page.goto(`/users/${username}`);
}

/**
 * 設定ページへ移動
 */
export async function goToSettings(page: Page) {
  await page.goto('/settings');
}

/**
 * いいねボタンをクリック
 */
export async function likePost(page: Page, postIndex: number = 0) {
  const likeButtons = await page.getByTestId('post-like-button').all();
  await likeButtons[postIndex].click();
  await page.waitForTimeout(500);
}

/**
 * コメントのいいねボタンをクリック
 */
export async function likeComment(page: Page, commentIndex: number = 0) {
  const likeButtons = await page.getByTestId('comment-like-button').all();
  await likeButtons[commentIndex].click();
  await page.waitForTimeout(500);
}

/**
 * フォローボタンをクリック
 */
export async function followUser(page: Page) {
  await page.getByTestId('profile-follow-button').click();
  await page.waitForTimeout(500);
}

/**
 * フォロー解除ボタンをクリック
 */
export async function unfollowUser(page: Page) {
  await page.getByTestId('profile-unfollow-button').click();
  await page.waitForTimeout(500);
}
