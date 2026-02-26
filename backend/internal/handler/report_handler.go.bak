package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/pkg/errors"
)

type ReportHandler struct {
	// サービスは後で実装
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
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
	// TODO: 実装
	return c.NoContent(201)
}
