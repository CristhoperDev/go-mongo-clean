package usecase

import (
	"context"
	"github.com/cristhoperdev/go-mongo-clean/domain"
	"time"
)

type movieUseCase struct {
	movieRepo	   domain.MovieRepository
	contextTimeout time.Duration
}

// NewMovieUsecase will create new an movieUsecase object representation of domain.movieUsecase interface
func NewMovieUsecase(movieRepo domain.MovieRepository, contextTimeout time.Duration) domain.MovieUseCae {
	return &movieUseCase{
		movieRepo,
		contextTimeout,
	}
}

func (m movieUseCase) Save(ctx context.Context, movie *domain.Movie) error {
	context, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.movieRepo.Save(context, movie)
	if err != nil {
		return  err
	}

	return err
}

func (m movieUseCase) Update(ctx context.Context, id string, movie *domain.Movie) error {
	context, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	movie.UpdatedAt = time.Now()
	return m.movieRepo.Update(context, id, movie)
}

func (m movieUseCase) Delete(ctx context.Context, id string) error {
	context, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	existedMovie, err := m.movieRepo.FindByID(context, id)
	if err != nil {
		return err
	}

	if existedMovie == (&domain.Movie{}) {
		return domain.ErrNotFound
	}

	return m.movieRepo.Delete(context, id)
}

func (m movieUseCase) FindByID(ctx context.Context, id string) (*domain.Movie, error) {
	context, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	movie, err := m.movieRepo.FindByID(context, id)
	return movie, err
}

func (m movieUseCase) FindAll(ctx context.Context) ([]*domain.Movie, error) {
	context, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	movie, err := m.movieRepo.FindAll(context)
	return movie, err
}