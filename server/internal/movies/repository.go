package movies

import (
	"context"

	"github.com/tim8912097887-sys/server/internal/shared"
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

func (r *repository) GetMovies(
	ctx context.Context,
	paginationParams types.PaginationParams,
) ([]types.Movie, int, error) {

	// Gaurd clause
	if paginationParams.Limit <= 0 {
        paginationParams.Limit = 10 // Default limit fallback
    }
	page := paginationParams.Page
    if page < 1 {
        page = 1
    }
	filter := bson.D{}

	totalCount64, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

    total := int(totalCount64)
	
	if total == 0 {
        return []types.Movie{}, 0, nil
    }

	limit := paginationParams.Limit
    totalPage := (total + limit - 1) / limit

	sortCriteria := bson.D{
		{Key: "created_at", Value: -1},
	}

	opts := options.Find().
		SetSort(sortCriteria).
		SetSkip(int64((paginationParams.Page - 1) * paginationParams.Limit)).
		SetLimit(int64(paginationParams.Limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var movies []types.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, 0, err
	}

	return movies, totalPage, nil
}

func (r *repository) GetMoviesByGenres(ctx context.Context, genres []types.Genres) ([]types.Movie, error) {
	generesIds := make([]int, len(genres))

	for i, genre := range genres {
		generesIds[i] = genre.GenreID
	}

	option := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(8)
	filter := bson.M{"genres.genre_id": bson.M{"$in": generesIds}}
 
	cursor, err := r.collection.Find(ctx, filter, option)
	if err != nil {
		return []types.Movie{}, err
	}

	defer cursor.Close(ctx)

	var movies []types.Movie

	err = cursor.All(ctx, &movies)

	if err != nil {
		return []types.Movie{}, err
	}

	return movies, nil

}

func (r *repository) GetMovieById(ctx context.Context, id string) (types.Movie, error) {
	var movie types.Movie

	idObjectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return types.Movie{}, shared.ErrMovieNotFound
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": idObjectID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Movie{}, shared.ErrMovieNotFound
		}
		return types.Movie{}, err
	}
	return movie, nil
}

func (r *repository) UpdateMovie (ctx context.Context, movie types.UpdateMovieSchema) error {
	update := bson.M{
        "$set": bson.M{
            "rating": movie.Rating,
        },
    }
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": movie.ID}, update)
	if err != nil {
		return err
	}
	return nil
}