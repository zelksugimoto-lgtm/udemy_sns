package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/admin/service"
	"github.com/yourusername/sns-app/internal/admin/session"
	"github.com/yourusername/sns-app/internal/repository"
)

// AdminHandler 管理画面ハンドラ
type AdminHandler struct {
	adminService    service.AdminService
	userMgmtService service.UserManagementService
	resetService    service.PasswordResetService
	postRepo        repository.PostRepository
	commentRepo     repository.CommentRepository
	templates       *template.Template
}

// NewAdminHandler 管理画面ハンドラを生成
func NewAdminHandler(
	adminService service.AdminService,
	userMgmtService service.UserManagementService,
	resetService service.PasswordResetService,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
) *AdminHandler {
	// カスタム関数を定義
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	// テンプレートをパース
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/admin/*.html")
	if err != nil {
		log.Fatal().Err(err).Msg("テンプレートのパースに失敗")
	}

	return &AdminHandler{
		adminService:    adminService,
		userMgmtService: userMgmtService,
		resetService:    resetService,
		postRepo:        postRepo,
		commentRepo:     commentRepo,
		templates:       tmpl,
	}
}

// ShowLoginPage ログイン画面を表示
func (h *AdminHandler) ShowLoginPage(c echo.Context) error {
	data := map[string]interface{}{
		"Error": c.QueryParam("error"),
	}
	return h.templates.ExecuteTemplate(c.Response().Writer, "login.html", data)
}

// Login ログイン処理
func (h *AdminHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// セッション取得
	sess, _ := session.GetStore().Get(c.Request(), session.SessionName)

	// まずルート管理者として認証を試みる
	isRoot, err := h.adminService.AuthenticateRoot(username, password)
	if err != nil {
		log.Error().Err(err).Msg("ルート管理者認証エラー")
	}

	if isRoot {
		// ルート管理者としてログイン
		sess.Values[session.KeyIsRoot] = true
		sess.Values[session.KeyUsername] = username
		sess.Values[session.KeyRole] = "root"
		sess.Save(c.Request(), c.Response().Writer)

		log.Info().Str("username", username).Msg("ルート管理者ログイン")
		return c.Redirect(http.StatusFound, "/admin")
	}

	// 通常の管理者として認証
	admin, err := h.adminService.AuthenticateAdmin(username, password)
	if err != nil {
		log.Error().Err(err).Msg("管理者認証エラー")
		return c.Redirect(http.StatusFound, "/admin/login?error=invalid")
	}

	if admin == nil {
		return c.Redirect(http.StatusFound, "/admin/login?error=invalid")
	}

	// セッションに保存
	sess.Values[session.KeyAdminID] = admin.ID
	sess.Values[session.KeyUsername] = admin.Username
	sess.Values[session.KeyRole] = admin.Role
	sess.Values[session.KeyIsRoot] = false
	sess.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusFound, "/admin")
}

// Logout ログアウト処理
func (h *AdminHandler) Logout(c echo.Context) error {
	sess, _ := session.GetStore().Get(c.Request(), session.SessionName)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusFound, "/admin/login")
}

// ShowDashboard ダッシュボードを表示
func (h *AdminHandler) ShowDashboard(c echo.Context) error {
	// ユーザー統計を取得
	userStats, err := h.userMgmtService.GetUserStats()
	if err != nil {
		log.Error().Err(err).Msg("ユーザー統計取得エラー")
		userStats = map[string]int64{}
	}

	// 投稿数を取得
	var postCount int64
	posts, total, err := h.postRepo.GetTimeline(uuid.Nil, []uuid.UUID{}, 1, 0)
	if err == nil {
		postCount = total
	}
	_ = posts // 使用していない警告を回避

	// コメント数を取得（簡易実装）
	var commentCount int64 = 0 // TODO: リポジトリにCountメソッドを追加

	// パスワードリセット申請数を取得
	resetRequests, resetTotal, err := h.resetService.ListRequests("pending", 10, 0)
	if err != nil {
		log.Error().Err(err).Msg("パスワードリセット申請取得エラー")
		resetTotal = 0
	}
	_ = resetRequests // 使用していない警告を回避

	data := map[string]interface{}{
		"Username":           c.Get("admin_username"),
		"IsRoot":             c.Get("is_root"),
		"PendingUsers":       userStats["pending"],
		"ApprovedUsers":      userStats["approved"],
		"RejectedUsers":      userStats["rejected"],
		"TotalPosts":         postCount,
		"TotalComments":      commentCount,
		"PendingResetRequests": resetTotal,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "dashboard", data)
}

// ListUsers ユーザー一覧を表示
func (h *AdminHandler) ListUsers(c echo.Context) error {
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	users, total, err := h.userMgmtService.ListUsers(status, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("ユーザー一覧取得エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザー一覧の取得に失敗しました")
	}

	totalPages := (int(total) + limit - 1) / limit

	data := map[string]interface{}{
		"Username":   c.Get("admin_username"),
		"IsRoot":     c.Get("is_root"),
		"Users":      users,
		"Status":     status,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "users", data)
}

// ShowUserDetail ユーザー詳細を表示
func (h *AdminHandler) ShowUserDetail(c echo.Context) error {
	userID := c.Param("id")

	user, err := h.userMgmtService.GetUserByID(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("ユーザー詳細取得エラー")
		return echo.NewHTTPError(http.StatusNotFound, "ユーザーが見つかりません")
	}

	data := map[string]interface{}{
		"Username": c.Get("admin_username"),
		"IsRoot":   c.Get("is_root"),
		"User":     user,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "user_detail", data)
}

// ApproveUser ユーザーを承認
func (h *AdminHandler) ApproveUser(c echo.Context) error {
	userID := c.Param("id")
	adminID, _ := c.Get("admin_id").(string)
	isRoot, _ := c.Get("is_root").(bool)

	if err := h.userMgmtService.ApproveUser(userID, adminID, isRoot); err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("ユーザー承認エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーの承認に失敗しました")
	}

	return c.Redirect(http.StatusFound, "/admin/users?status=pending")
}

// RejectUser ユーザーを拒否
func (h *AdminHandler) RejectUser(c echo.Context) error {
	userID := c.Param("id")
	adminID, _ := c.Get("admin_id").(string)
	isRoot, _ := c.Get("is_root").(bool)

	if err := h.userMgmtService.RejectUser(userID, adminID, isRoot); err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("ユーザー拒否エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーの拒否に失敗しました")
	}

	return c.Redirect(http.StatusFound, "/admin/users?status=pending")
}

// ListPasswordResets パスワードリセット申請一覧を表示
func (h *AdminHandler) ListPasswordResets(c echo.Context) error {
	// statusパラメータを取得（空文字列の場合は全件表示）
	status := c.QueryParam("status")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	requests, total, err := h.resetService.ListRequests(status, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("パスワードリセット申請一覧取得エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "パスワードリセット申請一覧の取得に失敗しました")
	}

	totalPages := (int(total) + limit - 1) / limit

	data := map[string]interface{}{
		"Username":   c.Get("admin_username"),
		"IsRoot":     c.Get("is_root"),
		"Requests":   requests,
		"Status":     status,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "password_resets", data)
}

// GenerateResetLink パスワードリセットリンクを発行
func (h *AdminHandler) GenerateResetLink(c echo.Context) error {
	requestID := c.Param("id")
	adminID, _ := c.Get("admin_id").(string)
	isRoot, _ := c.Get("is_root").(bool)

	resetLink, err := h.resetService.GenerateResetLink(requestID, adminID, isRoot)
	if err != nil {
		log.Error().Err(err).Str("request_id", requestID).Msg("パスワードリセットリンク発行エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "リセットリンクの発行に失敗しました")
	}

	// リセットリンクを表示する画面にリダイレクト（またはJSON返却）
	data := map[string]interface{}{
		"Username":  c.Get("admin_username"),
		"IsRoot":    c.Get("is_root"),
		"ResetLink": resetLink,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "reset_link_generated", data)
}

// ListAdmins 管理者一覧を表示
func (h *AdminHandler) ListAdmins(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	admins, total, err := h.adminService.ListAdmins(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("管理者一覧取得エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "管理者一覧の取得に失敗しました")
	}

	totalPages := (int(total) + limit - 1) / limit

	data := map[string]interface{}{
		"Username":   c.Get("admin_username"),
		"IsRoot":     c.Get("is_root"),
		"Admins":     admins,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "admins", data)
}

// ShowCreateAdminForm 管理者作成フォームを表示
func (h *AdminHandler) ShowCreateAdminForm(c echo.Context) error {
	data := map[string]interface{}{
		"Username": c.Get("admin_username"),
		"IsRoot":   c.Get("is_root"),
	}

	return h.templates.ExecuteTemplate(c.Response().Writer, "create_admin", data)
}

// CreateAdmin 管理者を作成
func (h *AdminHandler) CreateAdmin(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	role := c.FormValue("role")

	if role != "admin" && role != "root" {
		role = "admin"
	}

	_, err := h.adminService.CreateAdmin(username, email, password, role)
	if err != nil {
		log.Error().Err(err).Msg("管理者作成エラー")
		return echo.NewHTTPError(http.StatusInternalServerError, "管理者の作成に失敗しました")
	}

	return c.Redirect(http.StatusFound, "/admin/admins")
}
