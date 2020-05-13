package domain

import (
	"context"
	"time"
)

// Movie ...
type Movie struct {
	MovieID   string  	`json:"movieId"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// MovieUsecase represent the movie's usecases
type MovieUseCae interface {
	Save(ctx context.Context, movie *Movie) error
	Update(ctx context.Context, id string, movie *Movie) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Movie, error)
	FindAll(ctx context.Context) ([]*Movie, error)
}

// MovieRepository represent the movie's repository contract
type MovieRepository interface {
	Save(ctx context.Context, movie *Movie) error
	Update(ctx context.Context, id string, movie *Movie) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Movie, error)
	FindAll(ctx context.Context) ([]*Movie, error)
}