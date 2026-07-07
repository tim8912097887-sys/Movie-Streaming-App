package movies

import (
	"context"

	"github.com/tim8912097887-sys/server/internal/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repository struct {
	collection *mongo.Collection
}

func NewMovieRepository(collection *mongo.Collection) *repository {
	return &repository{
		collection: collection,
	}
}

func (r *repository) GetMovies(ctx context.Context, paginationParams PaginationParams) ([]Movie, error) {
	sortCriteria := bson.D{
		{Key: "created_at", Value: -1},
	}
	option := options.Find().SetSort(sortCriteria).SetSkip(int64(paginationParams.Offset)).SetLimit(int64(paginationParams.Limit))

	cursor, err := r.collection.Find(ctx, bson.D{}, option)
	if err != nil {
		return []Movie{}, err
	}

	defer cursor.Close(ctx)

	var movies []Movie

	err = cursor.All(ctx, &movies)

	if err != nil {
		return []Movie{}, err
	}

	return movies, nil
}

func (r *repository) GetMoviesByGenres(ctx context.Context, genres []types.Genres) ([]Movie, error) {
	generesIds := make([]int, len(genres))

	for i, genre := range genres {
		generesIds[i] = genre.GenreID
	}

	option := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(5)
	filter := bson.M{"genres.genre_id": bson.M{"$in": generesIds}}
 
	cursor, err := r.collection.Find(ctx, filter, option)
	if err != nil {
		return []Movie{}, err
	}

	defer cursor.Close(ctx)

	var movies []Movie

	err = cursor.All(ctx, &movies)

	if err != nil {
		return []Movie{}, err
	}

	return movies, nil

}