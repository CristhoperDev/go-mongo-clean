package http

import (
	"github.com/cristhoperdev/go-mongo-clean/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// MovieHandler  represent the httphandler for movie
type MovieHandler struct {
	MUseCase domain.MovieUseCae
}

// NewMovieHandler will initialize the movies/ resources endpoint
func NewMovieHandler(e *echo.Echo, MUseCase domain.MovieUseCae) {
	handler := &MovieHandler{
		MUseCase,
	}
	e.GET("/movies", handler.FetchMovie)
	e.POST("/movies", handler.Store)
	e.GET("/movies/:id", handler.GetByID)
	e.DELETE("/movies/:id", handler.Delete)
}

// FetchMovie will fetch all movies
func (m *MovieHandler) FetchMovie(c echo.Context) error {
	ctx := c.Request().Context()
	movie, err := m.MUseCase.FindAll(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, movie)
}

// GetByID will get movie by given id
func (m *MovieHandler) GetByID(c echo.Context) error {
	movieId := c.Param("id")
	if movieId == "" {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	movies, err := m.MUseCase.FindByID(ctx, movieId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, movies)
}

func isRequestValid(m *domain.Movie) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the movie by given request body
func (m *MovieHandler) Store(c echo.Context) (err error) {
	var movie domain.Movie
	err = c.Bind(&movie)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&movie); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = m.MUseCase.Save(ctx, &movie)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, movie)
}

// Delete will delete movie by given param
func (m *MovieHandler) Delete(c echo.Context) error {
	movieId := c.Param("id")
	if movieId == "" {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	err := m.MUseCase.Delete(ctx, movieId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}