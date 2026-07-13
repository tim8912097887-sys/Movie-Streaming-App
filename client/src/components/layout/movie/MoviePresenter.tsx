import Card from "../../ui/Card";

export type Movie = {
  title: string;
  description: string;
  image: string;
  rating: number;
};

type MoviePresenterProps = {
  movies: Movie[];
};

const MoviePresenter = ({ movies }: MoviePresenterProps) => {
  return (
    <div className="h-full w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-4 justify-items-center">
      {movies.map((movie) => (
        <Card key={movie.title}>
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
