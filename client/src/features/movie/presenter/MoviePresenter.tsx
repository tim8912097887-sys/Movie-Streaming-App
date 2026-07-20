import Card from "../ui/Card";
import CardSkeleton from "../ui/CardSkeleton";

export type Movie = {
  title: string;
  description: string;
  image: string;
  rating: number;
};

type MoviePresenterProps = {
  movies: Movie[];
  isLoading: boolean;
};

const MoviePresenter = ({ movies, isLoading }: MoviePresenterProps) => {
  if (isLoading && movies.length === 0) {
    return (
      <div className="h-full w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-4 justify-items-center">
        {Array.from({ length: 8 }).map((_, index) => (
          <CardSkeleton key={index} />
        ))}
      </div>
    );
  }

  if (movies.length === 0) {
    return (
      <div className="h-1/2 w-full flex items-center justify-center">
        <p className="text-2xl font-bold text-black tracking-tight">
          No movies found
        </p>
      </div>
    );
  }

  return (
    <div className="h-full w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-4 justify-items-center">
      {movies.map((movie, key) => (
        <Card key={key}>
          <Card.Head>
            <img
              src={movie.image}
              alt={movie.title}
              className="w-full h-full object-cover"
            />
          </Card.Head>
          <Card.Body>
            <h3 className="font-bold text-2xl text-white tracking-tight">
              {movie.title}
            </h3>
            <p className="text-xs text-gray-300 line-clamp-3 leading-relaxed">
              {movie.description}
            </p>
          </Card.Body>
          <Card.Footer>
            <p className="text-xs font-semibold text-rose-400">
              Rating: {movie.rating}
            </p>
          </Card.Footer>
        </Card>
      ))}
    </div>
  );
};

export default MoviePresenter;
