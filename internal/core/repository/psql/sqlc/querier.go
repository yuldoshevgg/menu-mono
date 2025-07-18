// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
)

type Querier interface {
	CountCategories(ctx context.Context, venueID string) (int64, error)
	CountMenuItems(ctx context.Context, venueID string) (int64, error)
	CountRoles(ctx context.Context) (int64, error)
	CountUsers(ctx context.Context) (int64, error)
	CountVenues(ctx context.Context) (int64, error)
	CreateCategory(ctx context.Context, arg CreateCategoryParams) error
	CreateMenuItem(ctx context.Context, arg CreateMenuItemParams) error
	CreateQrLink(ctx context.Context, arg CreateQrLinkParams) (string, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	CreateVenue(ctx context.Context, arg CreateVenueParams) error
	DeleteCategory(ctx context.Context, id string) error
	DeleteMenuItem(ctx context.Context, id string) error
	DeleteRole(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
	DeleteVenue(ctx context.Context, id string) error
	ExistVenue(ctx context.Context, id string) (bool, error)
	GetCategory(ctx context.Context, id string) (Category, error)
	GetMenuItem(ctx context.Context, id string) (MenuItem, error)
	GetRole(ctx context.Context, id string) (Role, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetVenue(ctx context.Context, id string) (Venue, error)
	ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error)
	ListMenuItems(ctx context.Context, arg ListMenuItemsParams) ([]MenuItem, error)
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListVenues(ctx context.Context, arg ListVenuesParams) ([]Venue, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error
	UpdateMenuItem(ctx context.Context, arg UpdateMenuItemParams) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateVenue(ctx context.Context, arg UpdateVenueParams) error
}

var _ Querier = (*Queries)(nil)
