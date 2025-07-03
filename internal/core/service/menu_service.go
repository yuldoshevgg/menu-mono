package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yuldoshevgg/menu-mono/generated/menu_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/helper"
	"github.com/yuldoshevgg/menu-mono/pkg/security"
	"google.golang.org/protobuf/types/known/emptypb"
)

type menuService struct {
	cfg  *config.Config
	repo repository.Store
	menu_service.UnimplementedMenuServiceServer
}

func NewMenuService(cfg *config.Config, repo repository.Store) *menuService {
	return &menuService{cfg: cfg, repo: repo}
}

func (s *menuService) CreateMenuItem(ctx context.Context, req *menu_service.CreateMenuItemRequest) (*emptypb.Empty, error) {
	var menuItem sqlc.CreateMenuItemParams

	if err := helper.MarshalUnMarshal(req, &menuItem); err != nil {
		return nil, fmt.Errorf("failed to unmarshal menu item: %w", err)
	}

	if err := s.repo.Postgres().CreateMenuItem(ctx, menuItem); err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *menuService) GetMenuItem(ctx context.Context, req *menu_service.GetMenuItemRequest) (*menu_service.GetMenuItemResponse, error) {
	var resp *menu_service.GetMenuItemResponse

	menuItem, err := s.repo.Postgres().GetMenuItem(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}

	if err := helper.MarshalUnMarshal(menuItem, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal menu item: %w", err)
	}

	return resp, nil
}

func (s *menuService) ListMenuItems(ctx context.Context, req *menu_service.ListMenuItemsRequest) (*menu_service.ListMenuItemsResponse, error) {
	var resp *menu_service.ListMenuItemsResponse

	menuItems, err := s.repo.Postgres().ListMenuItems(ctx, sqlc.ListMenuItemsParams{
		VenueID:    req.VenueId,
		Offset:     req.Page,
		Limit:      req.Limit,
		CategoryID: req.CategoryId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}

	if err := helper.MarshalUnMarshal(menuItems, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal menu items: %w", err)
	}

	total, err := s.repo.Postgres().CountMenuItems(ctx, req.VenueId)
	if err != nil {
		return nil, fmt.Errorf("failed to count menu items: %w", err)
	}

	resp.Total = int32(total)

	return resp, nil
}

func (s *menuService) UpdateMenuItem(ctx context.Context, req *menu_service.UpdateMenuItemRequest) (*emptypb.Empty, error) {
	var menuItem sqlc.UpdateMenuItemParams

	if err := helper.MarshalUnMarshal(req, &menuItem); err != nil {
		return nil, fmt.Errorf("failed to unmarshal menu item: %w", err)
	}

	if err := s.repo.Postgres().UpdateMenuItem(ctx, menuItem); err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *menuService) DeleteMenuItem(ctx context.Context, req *menu_service.DeleteMenuItemRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().DeleteMenuItem(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete menu item: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *menuService) AccessMenu(ctx context.Context, _ *emptypb.Empty) (*menu_service.AccessMenuResponse, error) {
	venueID, ok := ctx.Value(config.VenueID).(string)
	if !ok {
		return nil, fmt.Errorf("venue_id not found in context")
	}

	exist, err := s.repo.Postgres().ExistVenue(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if venue exists: %w", err)
	}

	if !exist {
		return nil, fmt.Errorf("venue not found")
	}

	tableNumber, ok := ctx.Value(config.TableNumber).(string)
	if !ok {
		return nil, fmt.Errorf("table_number not found in context")
	}

	tableNumberInt, err := strconv.ParseInt(tableNumber, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table number: %w", err)
	}

	id, err := s.repo.Postgres().CreateQrLink(ctx, sqlc.CreateQrLinkParams{
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
		VenueID:     venueID,
		TableNumber: int32(tableNumberInt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create qr link: %w", err)
	}

	valueToken := map[string]interface{}{
		"venue_id":     venueID,
		"table_number": tableNumber,
		"link_id":      id,
	}

	token := security.GenerateToken(valueToken, []byte(s.cfg.Token.SecretKey))

	return &menu_service.AccessMenuResponse{
		Link: fmt.Sprintf("%s/menu?token=%s", s.cfg.App.BaseURL, token),
	}, nil
}
