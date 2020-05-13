package mongodb

import (
	"context"
	"github.com/cristhoperdev/go-mongo-clean/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mongoMovieRepository struct {
	db 			*mongo.Database
	collection 	string
}

func (m mongoMovieRepository) Save(ctx context.Context, movie *domain.Movie) error {
	_, err := m.db.Collection(m.collection).InsertOne(ctx, movie)
	return err
}

func (m mongoMovieRepository) Update(ctx context.Context, id string, movie *domain.Movie) error {
	movie.UpdatedAt = time.Now()
	_, err := m.db.Collection(m.collection).UpdateOne(ctx, bson.M{"movieId": id}, movie)
	return err
}

func (m mongoMovieRepository) Delete(ctx context.Context, id string) error {
	_, err := m.db.Collection(m.collection).DeleteOne(ctx, bson.M{"movieid": id})
	return err
}

func (m mongoMovieRepository) FindByID(ctx context.Context, id string) (*domain.Movie, error) {
	var movie domain.Movie
	err := m.db.Collection(m.collection).FindOne(ctx, bson.M{"movieid": id}).Decode(&movie)
	return &movie, err
}

func (m mongoMovieRepository) FindAll(ctx context.Context) ([]*domain.Movie, error) {
	var movies []*domain.Movie
	cursor, err := m.db.Collection(m.collection).Find(ctx, bson.M{})
	err = cursor.All(ctx, &movies)
	if err != nil {
		return nil, err
	}
	return movies, err
}

// NewMongoMovieRepository will create an object that represent the movie.Repository interface
func NewMongoMovieRepository(Conn *mongo.Database, collection string) domain.MovieRepository {
	return &mongoMovieRepository{
		Conn,
		collection,
	}
}