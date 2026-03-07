package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8081/api/v1"
	runs    = 10
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type CreatePostRequest struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type PostResponse struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	Visibility    string    `json:"visibility"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	IsLiked       bool      `json:"is_liked"`
	IsBookmarked  bool      `json:"is_bookmarked"`
	User          UserSimple `json:"user"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserSimple struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

type TimelineResponse struct {
	Posts      []PostResponse `json:"posts"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("タイムラインAPI レスポンスタイム計測")
	fmt.Println("========================================")
	fmt.Println()

	// 1. テストユーザー登録
	fmt.Println("1. テストユーザーを登録中...")
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("benchmarkuser%d", timestamp)
	email := fmt.Sprintf("benchmark%d@example.com", timestamp)

	registerReq := RegisterRequest{
		Email:       email,
		Password:    "password123",
		Username:    username,
		DisplayName: "Benchmark User",
	}

	token, err := registerUser(registerReq)
	if err != nil {
		fmt.Printf("エラー: ユーザー登録に失敗しました: %v\n", err)
		return
	}
	fmt.Printf("✓ ユーザー登録完了（トークン取得済み）\n")
	fmt.Println()

	// 2. テストデータ作成（20件の投稿）
	fmt.Println("2. テストデータを作成中（20件の投稿）...")
	for i := 1; i <= 20; i++ {
		content := fmt.Sprintf("ベンチマークテスト投稿 #%d - %s", i, time.Now().Format("15:04:05"))
		err := createPost(token, content)
		if err != nil {
			fmt.Printf("エラー: 投稿作成に失敗しました: %v\n", err)
			return
		}
		if i%5 == 0 {
			fmt.Printf("✓ %d件作成完了\n", i)
		}
	}
	fmt.Println("✓ テストデータ作成完了")
	fmt.Println()

	// 3. レスポンスタイム計測（10回実行）
	fmt.Printf("3. タイムラインAPIのレスポンスタイムを計測中（%d回実行）...\n", runs)
	var durations []time.Duration
	var totalDuration time.Duration

	for i := 1; i <= runs; i++ {
		duration, err := measureTimeline(token)
		if err != nil {
			fmt.Printf("エラー: タイムライン取得に失敗しました: %v\n", err)
			return
		}
		durations = append(durations, duration)
		totalDuration += duration
		fmt.Printf("  実行 %2d: %6.2f ms\n", i, float64(duration.Microseconds())/1000)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("計測結果")
	fmt.Println("========================================")

	// 平均値
	avgDuration := totalDuration / time.Duration(runs)
	fmt.Printf("平均レスポンスタイム: %.2f ms\n", float64(avgDuration.Microseconds())/1000)

	// 最小値・最大値
	minDuration := durations[0]
	maxDuration := durations[0]
	for _, d := range durations {
		if d < minDuration {
			minDuration = d
		}
		if d > maxDuration {
			maxDuration = d
		}
	}
	fmt.Printf("最小値: %.2f ms\n", float64(minDuration.Microseconds())/1000)
	fmt.Printf("最大値: %.2f ms\n", float64(maxDuration.Microseconds())/1000)

	fmt.Println()
	fmt.Println("========================================")
}

func registerUser(req RegisterRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ステータスコード: %d, レスポンス: %s", resp.StatusCode, string(bodyBytes))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Token, nil
}

func createPost(token, content string) error {
	postReq := CreatePostRequest{
		Content:    content,
		Visibility: "public",
	}

	body, err := json.Marshal(postReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/posts", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ステータスコード: %d, レスポンス: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func measureTimeline(token string) (time.Duration, error) {
	req, err := http.NewRequest("GET", baseURL+"/timeline?limit=20&offset=0", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// レスポンスボディを完全に読み込む
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	duration := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ステータスコード: %d", resp.StatusCode)
	}

	return duration, nil
}
