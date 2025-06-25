package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yuldoshevgg/menu-mono/generated/auth_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/helper"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/guregu/null.v4/zero"
)

type roleService struct {
	cfg  *config.Config
	repo repository.Store
	auth_service.UnimplementedRoleServiceServer
}

func NewRoleService(cfg *config.Config, repo repository.Store) *roleService {
	return &roleService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *roleService) CreateRole(ctx context.Context, req *auth_service.CreateRoleRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().CreateRole(ctx, sqlc.CreateRoleParams{
		Title:       req.Title,
		Description: zero.StringFrom(req.Description),
	}); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *roleService) GetRole(ctx context.Context, req *auth_service.GetRoleRequest) (*auth_service.GetRoleResponse, error) {
	role, err := s.repo.Postgres().GetRole(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &auth_service.GetRoleResponse{
		Id:          role.ID,
		Title:       role.Title,
		Description: role.Description.String,
		CreatedAt:   role.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *roleService) UpdateRole(ctx context.Context, req *auth_service.UpdateRoleRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().UpdateRole(ctx, sqlc.UpdateRoleParams{
		ID:          req.Id,
		Title:       zero.StringFrom(req.Title),
		Description: zero.StringFrom(req.Description),
	}); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *roleService) DeleteRole(ctx context.Context, req *auth_service.DeleteRoleRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().DeleteRole(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete role: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *roleService) ListRoles(ctx context.Context, req *auth_service.ListRolesRequest) (*auth_service.ListRolesResponse, error) {
	var resp []*auth_service.GetRoleResponse

	roles, err := s.repo.Postgres().ListRoles(ctx, sqlc.ListRolesParams{
		Limit:  req.Limit,
		Offset: req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	if err := helper.MarshalUnMarshal(roles, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal roles: %w", err)
	}

	total, err := s.repo.Postgres().CountRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count roles: %w", err)
	}

	return &auth_service.ListRolesResponse{
		Roles: resp,
		Total: int32(total),
	}, nil
}
