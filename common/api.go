package common

import "context"

// API interface
type API[T any, I comparable] interface {
	Create(ctx context.Context, params T) (*I, error)
	Read(ctx context.Context, id I) (*T, error)
	Update(ctx context.Context, params T) (*T, error)
	Delete(ctx context.Context, id I) error
}
