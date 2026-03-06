import { test, expect } from '@playwright/test';
import { registerUser, createPost, waitForPostToAppear, followUser, unfollowUser, goToProfile, goToSettings } from '../helpers/test-helpers';

test.describe('フォロー・プロフィール機能', () => {
  test('フォロー → フォロー解除 → フォロー一覧に反映', async ({ page, context }) => {
    // 1人目のユーザーを登録
    const user1 = await registerUser(page);

    // 投稿を作成（フォローするユーザーを見つけやすくするため）
    const postContent = `ユーザー1の投稿 ${Date.now()}`;
    await createPost(page, postContent);
    await waitForPostToAppear(page, postContent);

    // ログアウト
    await page.getByRole('button', { name: user1.displayName.charAt(0).toUpperCase() }).click();
    await page.getByTestId('header-logout').click();

    // 2人目のユーザーを登録
    const user2 = await registerUser(page);

    // user1のプロフィールページへ移動
    await goToProfile(page, user1.username);

    // フォローボタンをクリック
    await followUser(page);

    // user1のフォロワー数が1になることを確認（user2がフォローした）
    await expect(page.getByTestId('profile-followers-count')).toContainText('1');

    // フォロー解除ボタンをクリック
    await unfollowUser(page);

    // user1のフォロワー数が0になることを確認
    await expect(page.getByTestId('profile-followers-count')).toContainText('0');
  });

  test('プロフィール編集', async ({ page }) => {
    // ユーザー登録
    const user = await registerUser(page);

    // 設定ページへ移動
    await goToSettings(page);

    // 表示名と自己紹介を変更
    const newDisplayName = `更新後の表示名 ${Date.now()}`;
    const newBio = `更新後の自己紹介 ${Date.now()}`;

    await page.getByTestId('settings-displayname').fill(newDisplayName);
    await page.getByTestId('settings-bio').fill(newBio);
    await page.getByTestId('settings-submit').click();

    // 成功メッセージが表示されることを確認
    await expect(page.locator('.MuiAlert-message').filter({ hasText: 'プロフィールを更新しました' })).toBeVisible({ timeout: 5000 });

    // プロフィールページへ移動
    await goToProfile(page, user.username);

    // 更新された表示名が表示されることを確認
    await expect(page.getByText(newDisplayName)).toBeVisible();

    // 更新された自己紹介が表示されることを確認
    await expect(page.getByText(newBio)).toBeVisible();
  });
});
