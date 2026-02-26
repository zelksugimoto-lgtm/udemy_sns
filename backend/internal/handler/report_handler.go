package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/middleware"
	"github.com/yourusername/sns-app/internal/service"
	"github.com/yourusername/sns-app/pkg/errors"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// CreateReport godoc
// @Summary      通報作成
// @Description  投稿、コメント、ユーザーを通報
// @Tags         reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  request.CreateReportRequest  true  "通報内容"
// @Success      201
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /reports [post]
func (h *ReportHandler) CreateReport(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Unauthorized("認証が必要です"))
	}

	// リクエストボディをパース
	var req request.CreateReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効なリクエストです"))
	}

	// バリデーション
	if req.TargetType == "" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("通報対象タイプは必須です"))
	}
	if req.TargetType != "Post" && req.TargetType != "Comment" && req.TargetType != "User" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な通報対象タイプです"))
	}
	if req.Reason == "" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("通報理由は必須です"))
	}
	if req.Reason != "spam" && req.Reason != "inappropriate_content" && req.Reason != "harassment" && req.Reason != "other" {
		return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な通報理由です"))
	}

	// サービス層を呼び出し
	err = h.reportService.CreateReport(userID, &req)
	if err != nil {
		if err.Error() == "投稿が見つかりません" || err.Error() == "コメントが見つかりません" || err.Error() == "ユーザーが見つかりません" {
			return c.JSON(http.StatusNotFound, errors.NotFound(err.Error()))
		}
		if err.Error() == "無効な通報対象タイプです" {
			return c.JSON(http.StatusBadRequest, errors.BadRequest("無効な通報対象タイプです"))
		}
		return c.JSON(http.StatusInternalServerError, errors.InternalError("通報の作成に失敗しました"))
	}

	return c.NoContent(http.StatusCreated)
}
