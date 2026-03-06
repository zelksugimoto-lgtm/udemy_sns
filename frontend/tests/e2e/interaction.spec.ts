import { test, expect } from '@playwright/test';
import { registerUser, createPost, waitForPostToAppear, likePost, likeComment } from '../helpers/test-helpers';

test.describe('インタラクション機能', () => {
  test('いいね → いいね解除', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `いいねテスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, postContent);

    // いいねボタンをクリック
    await likePost(page, 0);

    // いいねされたことを確認（アイコンの色が変わる）
    const likeButton = page.getByTestId('post-like-button').first();
    await expect(likeButton).toHaveClass(/MuiIconButton-colorError/);

    // いいねを解除
    await likePost(page, 0);

    // いいねが解除されたことを確認（デフォルトカラーに戻る）
    await expect(likeButton).not.toHaveClass(/MuiIconButton-colorError/);
  });

  test('コメント作成 → コメント削除', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `コメントテスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, postContent);

    // 投稿をクリックして詳細ページへ
    await page.getByTestId('post-content').filter({ hasText: postContent }).first().click();

    // 投稿詳細ページに遷移したことを確認
    await expect(page).toHaveURL(/\/posts\/[a-f0-9-]+/);

    // コメントを作成
    const commentContent = `テストコメント ${Date.now()}`;
    await page.getByTestId('comment-form-content').fill(commentContent);
    await page.getByTestId('comment-form-submit').click();

    // コメントが表示されるまで待つ
    await expect(page.getByTestId('comment-content').filter({ hasText: commentContent })).toBeVisible({ timeout: 5000 });

    // コメントの3点メニューを開く
    const commentCard = page.getByTestId('comment-card').filter({ has: page.getByTestId('comment-content').filter({ hasText: commentContent }) }).first();
    await commentCard.getByTestId('comment-menu-button').click();

    // 確認ダイアログでOKをクリック（イベントハンドラを先に設定）
    page.once('dialog', dialog => dialog.accept());

    // 削除ボタンをクリック
    await page.getByTestId('comment-delete-button').click();

    // コメントが削除されて表示されなくなることを確認
    await expect(page.getByTestId('comment-content').filter({ hasText: commentContent })).not.toBeVisible({ timeout: 5000 });
  });

  test('コメントのいいね → コメントのいいね解除', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `コメントいいねテスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, postContent);

    // 投稿をクリックして詳細ページへ
    await page.getByTestId('post-content').filter({ hasText: postContent }).first().click();

    // コメントを作成
    const commentContent = `テストコメント ${Date.now()}`;
    await page.getByTestId('comment-form-content').fill(commentContent);
    await page.getByTestId('comment-form-submit').click();

    // コメントが表示されるまで待つ
    await expect(page.getByTestId('comment-content').filter({ hasText: commentContent })).toBeVisible({ timeout: 5000 });

    // コメントのいいねボタンをクリック
    await likeComment(page, 0);

    // いいねされたことを確認
    const commentLikeButton = page.getByTestId('comment-like-button').first();
    await expect(commentLikeButton).toHaveClass(/MuiIconButton-colorError/);

    // いいねを解除
    await likeComment(page, 0);

    // いいねが解除されたことを確認（デフォルトカラーに戻る）
    await expect(commentLikeButton).not.toHaveClass(/MuiIconButton-colorError/);
  });

  test('コメントの返信 → 返信として表示される', async ({ page }) => {
    // ユーザー登録
    await registerUser(page);

    // 投稿を作成
    const postContent = `返信テスト投稿 ${Date.now()}`;
    await createPost(page, postContent);

    // 投稿が表示されるまで待つ
    await waitForPostToAppear(page, postContent);

    // 投稿をクリックして詳細ページへ
    await page.getByTestId('post-content').filter({ hasText: postContent }).first().click();

    // 親コメントを作成
    const parentCommentContent = `親コメント ${Date.now()}`;
    await page.getByTestId('comment-form-content').fill(parentCommentContent);
    await page.getByTestId('comment-form-submit').click();

    // 親コメントが表示されるまで待つ
    await expect(page.getByTestId('comment-content').filter({ hasText: parentCommentContent })).toBeVisible({ timeout: 5000 });

    // 返信ボタンをクリック
    await page.getByTestId('comment-reply-button').first().click();

    // 返信フォームが表示されることを確認
    await expect(page.getByTestId('comment-form-content').nth(1)).toBeVisible();

    // 返信を作成
    const replyContent = `返信コメント ${Date.now()}`;
    await page.getByTestId('comment-form-content').nth(1).fill(replyContent);
    await page.getByTestId('comment-form-submit').nth(1).click();

    // 返信が親コメントの下に表示されることを確認
    await expect(page.getByTestId('comment-content').filter({ hasText: replyContent })).toBeVisible({ timeout: 5000 });
  });
});
