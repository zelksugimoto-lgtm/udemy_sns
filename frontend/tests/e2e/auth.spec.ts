import { test, expect } from '@playwright/test';
import { registerUser, loginUser, logout, generateRandomUser } from '../helpers/test-helpers';

test.describe('認証機能', () => {
  test('ユーザー登録 → ログイン', async ({ page }) => {
    // ユーザー登録
    const user = await registerUser(page);

    // ホーム画面にリダイレクトされることを確認
    await expect(page).toHaveURL('/');

    // ログアウト
    await logout(page);

    // ログインページにリダイレクトされることを確認
    await expect(page).toHaveURL('/login');

    // 再度ログイン
    await loginUser(page, user.email, user.password);

    // ホーム画面にリダイレクトされることを確認
    await expect(page).toHaveURL('/');
  });

  test('登録失敗（メールアドレス重複）', async ({ page }) => {
    // 1人目のユーザーを登録
    const user = generateRandomUser();
    await registerUser(page, user);

    // ログアウト
    await logout(page);

    // 同じメールアドレスで再度登録を試みる
    await page.goto('/register');
    await page.getByTestId('register-email').fill(user.email);
    await page.getByTestId('register-username').fill(`${user.username}_different`);
    await page.getByTestId('register-displayname').fill(`${user.displayName} Different`);
    await page.getByTestId('register-password').fill(user.password);
    await page.getByTestId('register-submit').click();

    // エラーメッセージが表示されることを確認
    await expect(page.locator('.MuiAlert-message')).toBeVisible({ timeout: 5000 });
  });

  test('ログイン失敗（パスワード間違い）', async ({ page }) => {
    // ユーザー登録
    const user = await registerUser(page);

    // ログアウト
    await logout(page);

    // 間違ったパスワードでログイン
    await page.goto('/login');
    await page.getByTestId('login-email').fill(user.email);
    await page.getByTestId('login-password').fill('wrongpassword123');
    await page.getByTestId('login-submit').click();

    // エラーメッセージが表示されることを確認
    await expect(page.locator('.MuiAlert-message')).toBeVisible({ timeout: 5000 });
  });

  test('ログアウト → 認証が必要なページにアクセス → ログインページにリダイレクト', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // ログアウト
    await logout(page);

    // 設定ページにアクセス（認証が必要）
    await page.goto('/settings');

    // ログインページにリダイレクトされることを確認
    await expect(page).toHaveURL('/login');
  });
});
