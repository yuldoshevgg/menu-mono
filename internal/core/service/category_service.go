package service

import (
	"context"
	"fmt"

	"github.com/yuldoshevgg/menu-mono/generated/menu_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/helper"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/guregu/null.v4/zero"
)

type categoryService struct {
	cfg  *config.Config
	repo repository.Store
	menu_service.UnimplementedCategoryServiceServer
}

func NewCategoryService(cfg *config.Config, repo repository.Store) *categoryService {
	return &categoryService{cfg: cfg, repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *menu_service.CreateCategoryRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().CreateCategory(ctx, sqlc.CreateCategoryParams{
		Name:    req.Name,
		VenueID: req.VenueId,
	}); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *categoryService) GetCategory(ctx context.Context, req *menu_service.GetCategoryRequest) (*menu_service.GetCategoryResponse, error) {
	var resp *menu_service.GetCategoryResponse

	category, err := s.repo.Postgres().GetCategory(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if err := helper.MarshalUnMarshal(category, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal category: %w", err)
	}

	return resp, nil
}

func (s *categoryService) ListCategories(ctx context.Context, req *menu_service.ListCategoriesRequest) (*menu_service.ListCategoriesResponse, error) {
	var resp *menu_service.ListCategoriesResponse

	categories, err := s.repo.Postgres().ListCategories(ctx, sqlc.ListCategoriesParams{
		VenueID: req.VenueId,
		Limit:   req.Limit,
		Offset:  req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	if err := helper.MarshalUnMarshal(categories, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal categories: %w", err)
	}

	total, err := s.repo.Postgres().CountCategories(ctx, req.VenueId)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	resp.Total = int32(total)

	return resp, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, req *menu_service.UpdateCategoryRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:      req.Id,
		Name:    zero.StringFrom(req.Name),
		VenueID: zero.StringFrom(req.VenueId),
	}); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, req *menu_service.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().DeleteCategory(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return &emptypb.Empty{}, nil
}
