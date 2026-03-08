package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
)

type UploadHandler struct {
	storageService *service.StorageService
}

func NewUploadHandler(storageService *service.StorageService) *UploadHandler {
	return &UploadHandler{
		storageService: storageService,
	}
}

// FileUploadResponse ファイルアップロードレスポンス
type FileUploadResponse struct {
	MediaURL  string `json:"media_url"`  // Firebase Storage公開URL
	MediaType string `json:"media_type"`
}

// UploadImage godoc
// @Summary      画像アップロード
// @Description  画像をFirebase Storageにアップロードし、URLを返す
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "アップロードする画像ファイル"
// @Success      200   {object}  FileUploadResponse
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Failure      413   {object}  errors.ErrorResponse
// @Failure      500   {object}  errors.ErrorResponse
// @Router       /upload/image [post]
func (h *UploadHandler) UploadImage(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// ファイルを取得
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("ファイルの取得に失敗")
		return c.JSON(http.StatusBadRequest, errors.BadRequest("ファイルが見つかりません"))
	}

	// ファイルサイズ制限チェック
	maxSizeMB, _ := strconv.Atoi(os.Getenv("MAX_IMAGE_SIZE_MB"))
	if maxSizeMB == 0 {
		maxSizeMB = 5 // デフォルト5MB
	}
	maxSizeBytes := int64(maxSizeMB * 1024 * 1024)

	if file.Size > maxSizeBytes {
		return c.JSON(http.StatusRequestEntityTooLarge, errors.BadRequest(fmt.Sprintf("ファイルサイズは%dMB以下にしてください", maxSizeMB)))
	}

	// ファイル形式チェック
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("対応していない画像形式です。JPEG、PNG、GIF、WebPのみアップロード可能です。"))
	}

	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("ファイルを開くことができません")
		return c.JSON(http.StatusInternalServerError, errors.InternalError("ファイルの読み込みに失敗しました"))
	}
	defer src.Close()

	// ファイル名生成
	fileExtension := filepath.Ext(file.Filename)
	if fileExtension == "" {
		// Content-Typeから拡張子を推測
		switch contentType {
		case "image/jpeg":
			fileExtension = ".jpg"
		case "image/png":
			fileExtension = ".png"
		case "image/gif":
			fileExtension = ".gif"
		case "image/webp":
			fileExtension = ".webp"
		default:
			fileExtension = ".jpg"
		}
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	storagePath := fmt.Sprintf("posts/%s/%s", userID.String(), fileName)

	// Firebase Storageにアップロード
	downloadURL, err := h.storageService.UploadFile(c.Request().Context(), storagePath, src, contentType)
	if err != nil {
		log.Error().Err(err).Str("path", storagePath).Msg("Firebase Storageへのアップロード失敗")
		return c.JSON(http.StatusInternalServerError, errors.InternalError("画像のアップロードに失敗しました"))
	}

	// メディアタイプを決定
	mediaType := "image"
	if strings.HasPrefix(contentType, "video/") {
		mediaType = "video"
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("file_name", fileName).
		Str("media_type", mediaType).
		Str("download_url", downloadURL).
		Int64("file_size", file.Size).
		Msg("ファイルアップロード成功")

	return c.JSON(http.StatusOK, FileUploadResponse{
		MediaURL:  downloadURL,
		MediaType: mediaType,
	})
}
