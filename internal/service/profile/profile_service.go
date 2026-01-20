package profileService

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/model/request"
	"auth_service/internal/model/user"
	tokenrepo "auth_service/internal/repository/token"
	userrepo "auth_service/internal/repository/user"
	db "auth_service/internal/storage/minio"
	"auth_service/pkg/validation"

	//"auth_service/internal/utils"

	"github.com/minio/minio-go/v7"
)

type Profile_Service interface {
	GetProfile(ctx context.Context, userID int64) *ProfileService
	UpdateProfile(ctx context.Context, userID int64, req request.UpdateProfileRequest) (*user.User, error)
	DeleteProfile(ctx context.Context, userID int64) error
	UploadPhoto(ctx context.Context, userID int64, file io.Reader, fileName string, fileSize int64) (string, error)
}

type ProfileService struct {
	userRepo  *userrepo.UserRepository
	tokenRepo *tokenrepo.TokenRepository
}

func NewProfileService(userRepo *userrepo.UserRepository, tokenRepo *tokenrepo.TokenRepository) *ProfileService {
	return &ProfileService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID int64) (*user.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID int64, req request.UpdateProfileRequest) (*user.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && !validation.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if req.Email != "" && req.Email != user.Email.String {
		exists, err := s.userRepo.CheckEmailExists(ctx, req.Email, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already in use")
		}
	}

	user.Name = validation.SanitizeInput(req.Name)
	if req.Email != "" {
		user.Email = sql.NullString{String: req.Email, Valid: true}
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *ProfileService) DeleteProfile(ctx context.Context, userID int64) error {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	err = s.tokenRepo.DeleteRefreshToken(ctx, userID)

	return err
}

func (s *ProfileService) UploadPhoto(ctx context.Context, userID int64, file io.Reader, fileName string, fileSize int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	cfg := config.App.Minio

	ext := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("users/%d/profile_%d%s", userID, time.Now().UnixNano(), ext)

	_, err = db.MinioClient.PutObject(ctx, cfg.Bucket, newFileName, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/" + strings.TrimPrefix(ext, "."),
	})
	if err != nil {
		return "", err
	}

	photoURL := fmt.Sprintf("http://%s/%s/%s", cfg.Domain, cfg.Bucket, newFileName)

	user.PhotoURL = sql.NullString{String: photoURL, Valid: true}
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return "", err
	}
	return photoURL, nil
}
