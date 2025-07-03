package service

import (
	"context"
	"fmt"

	"github.com/yuldoshevgg/menu-mono/generated/auth_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/helper"
	"github.com/yuldoshevgg/menu-mono/pkg/security"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/guregu/null.v4/zero"
)

type userService struct {
	cfg  *config.Config
	repo repository.Store
	auth_service.UnimplementedUserServiceServer
}

func NewUserService(cfg *config.Config, repo repository.Store) *userService {
	return &userService{cfg: cfg, repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req *auth_service.CreateUserRequest) (*emptypb.Empty, error) {
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Postgres().CreateUser(ctx, sqlc.CreateUserParams{
		Username:    req.Username,
		Password:    hashedPassword,
		RoleID:      req.RoleId,
		VenueID:     req.VenueId,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		MiddleName:  zero.StringFrom(req.MiddleName),
	}); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *userService) GetUser(ctx context.Context, req *auth_service.GetUserRequest) (*auth_service.GetUserResponse, error) {
	var resp *auth_service.GetUserResponse

	user, err := s.repo.Postgres().GetUser(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := helper.MarshalUnMarshal(user, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *auth_service.UpdateUserRequest) (*emptypb.Empty, error) {
	var (
		hashedPassword string
		err            error
	)

	if req.Password != "" {
		hashedPassword, err = security.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
	}

	if err := s.repo.Postgres().UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          req.Id,
		Username:    zero.StringFrom(req.Username),
		Password:    zero.StringFrom(hashedPassword),
		RoleID:      zero.StringFrom(req.RoleId),
		VenueID:     zero.StringFrom(req.VenueId),
		FirstName:   zero.StringFrom(req.FirstName),
		LastName:    zero.StringFrom(req.LastName),
		PhoneNumber: zero.StringFrom(req.PhoneNumber),
		MiddleName:  zero.StringFrom(req.MiddleName),
	}); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().DeleteUser(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *userService) ListUsers(ctx context.Context, req *auth_service.ListUsersRequest) (*auth_service.ListUsersResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	users, err := s.repo.Postgres().ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var resp []*auth_service.GetUserResponse

	if err := helper.MarshalUnMarshal(users, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal users: %w", err)
	}

	total, err := s.repo.Postgres().CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &auth_service.ListUsersResponse{
		Users: resp,
		Total: int32(total),
	}, nil
}
