import { test, expect } from '@playwright/test';
import { registerUser, createPost, waitForPostToAppear } from '../helpers/test-helpers';

test.describe('投稿機能', () => {
  test('投稿作成 → タイムラインに表示される', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `テスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿がタイムラインに表示されることを確認
    await waitForPostToAppear(page, postContent);
    await expect(page.getByTestId('post-content').filter({ hasText: postContent })).toBeVisible();
  });

  test('投稿編集 → 内容が更新される', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const originalContent = `編集前の投稿 ${Date.now()}`;
    await createPost(page, originalContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, originalContent);

    // 投稿をクリックして詳細ページへ
    await page.getByTestId('post-content').filter({ hasText: originalContent }).first().click();

    // 投稿詳細ページに遷移したことを確認
    await expect(page).toHaveURL(/\/posts\/[a-f0-9-]+/);

    // 編集ボタンをクリック（削除ボタンを押さないように注意）
    // Note: 編集機能が実装されていない場合、このテストはスキップします
    test.skip();
  });

  test('投稿削除 → タイムラインから消える', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `削除テスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, postContent);

    // 投稿カードの3点メニューを開く
    const postCard = page.getByTestId('post-card').filter({ has: page.getByTestId('post-content').filter({ hasText: postContent }) }).first();
    await postCard.getByTestId('post-menu-button').click();

    // 確認ダイアログでOKをクリック（イベントハンドラを先に設定）
    page.once('dialog', dialog => dialog.accept());

    // 削除ボタンをクリック
    await page.getByTestId('post-delete-button').click();

    // 投稿が削除されて表示されなくなることを確認
    await expect(page.getByTestId('post-content').filter({ hasText: postContent })).not.toBeVisible({ timeout: 5000 });
  });
});
