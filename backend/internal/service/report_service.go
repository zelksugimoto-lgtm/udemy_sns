package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/sns-app/internal/dto/request"
	"github.com/yourusername/sns-app/internal/model"
	"github.com/yourusername/sns-app/internal/repository"
)

// ReportService 通報サービスのインターフェース
type ReportService interface {
	CreateReport(userID uuid.UUID, req *request.CreateReportRequest) error
}

type reportService struct {
	reportRepo  repository.ReportRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

// NewReportService 通報サービスのコンストラクタ
func NewReportService(
	reportRepo repository.ReportRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
) ReportService {
	return &reportService{
		reportRepo:  reportRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

// CreateReport 通報作成
func (s *reportService) CreateReport(userID uuid.UUID, req *request.CreateReportRequest) error {
	// 通報対象の存在確認
	switch req.TargetType {
	case "Post":
		post, err := s.postRepo.FindByID(req.TargetID)
		if err != nil {
			return err
		}
		if post == nil {
			return errors.New("投稿が見つかりません")
		}
	case "Comment":
		comment, err := s.commentRepo.FindByID(req.TargetID)
		if err != nil {
			return err
		}
		if comment == nil {
			return errors.New("コメントが見つかりません")
		}
	case "User":
		user, err := s.userRepo.FindByID(req.TargetID)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("ユーザーが見つかりません")
		}
	default:
		return errors.New("無効な通報対象タイプです")
	}

	report := &model.Report{
		ReporterID: userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
		Comment:    req.Comment,
		Status:     "pending",
	}

	return s.reportRepo.Create(report)
}
