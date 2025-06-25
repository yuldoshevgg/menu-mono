package service

import (
	"context"
	"fmt"

	"github.com/yuldoshevgg/menu-mono/generated/general_service"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository/psql/sqlc"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/helper"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/guregu/null.v4/zero"
)

type venueService struct {
	cfg  *config.Config
	repo repository.Store
	general_service.UnimplementedVenueServiceServer
}

func NewVenueService(cfg *config.Config, repo repository.Store) *venueService {
	return &venueService{cfg: cfg, repo: repo}
}

func (s *venueService) CreateVenue(ctx context.Context, req *general_service.CreateVenueRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().CreateVenue(ctx, sqlc.CreateVenueParams{
		Title: req.Title,
	}); err != nil {
		return nil, fmt.Errorf("failed to create venue: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *venueService) GetVenue(ctx context.Context, req *general_service.GetVenueRequest) (*general_service.GetVenueResponse, error) {
	var resp *general_service.GetVenueResponse

	venue, err := s.repo.Postgres().GetVenue(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get venue: %w", err)
	}

	if err := helper.MarshalUnMarshal(venue, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal venue: %w", err)
	}

	return resp, nil
}

func (s *venueService) UpdateVenue(ctx context.Context, req *general_service.UpdateVenueRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().UpdateVenue(ctx, sqlc.UpdateVenueParams{
		ID:          req.Id,
		Title:       zero.StringFrom(req.Title),
		Slug:        zero.StringFrom(req.Slug),
		VenueID:     zero.StringFrom(req.VenueId),
		LogoUrl:     zero.StringFrom(req.LogoUrl),
		Status:      zero.BoolFrom(req.Status),
		Description: zero.StringFrom(req.Description),
	}); err != nil {
		return nil, fmt.Errorf("failed to update venue: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *venueService) DeleteVenue(ctx context.Context, req *general_service.DeleteVenueRequest) (*emptypb.Empty, error) {
	if err := s.repo.Postgres().DeleteVenue(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete venue: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *venueService) ListVenues(ctx context.Context, req *general_service.ListVenuesRequest) (*general_service.ListVenuesResponse, error) {
	venues, err := s.repo.Postgres().ListVenues(ctx, sqlc.ListVenuesParams{
		Limit:  req.Limit,
		Offset: req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list venues: %w", err)
	}

	var resp []*general_service.GetVenueResponse

	if err := helper.MarshalUnMarshal(venues, &resp); err != nil {
		return nil, fmt.Errorf("failed to marshal venues: %w", err)
	}

	total, err := s.repo.Postgres().CountVenues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count venues: %w", err)
	}

	return &general_service.ListVenuesResponse{
		Venues: resp,
		Total:  int32(total),
	}, nil
}
