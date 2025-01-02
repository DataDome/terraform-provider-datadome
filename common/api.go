package common

import "context"

// API interface
type API[T any] interface {
	Create(ctx context.Context, params T) (*int, error)
	Read(ctx context.Context, id int) (*T, error)
	Update(ctx context.Context, params T) (*T, error)
	Delete(ctx context.Context, id int) error
}
