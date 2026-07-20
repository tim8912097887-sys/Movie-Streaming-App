export type MovieData = {
  title: string;
  description: string;
  poster_url: string;
  rating: number;
};

export type FetchMoviesData = {
  total_page: number;
  current_page: number;
  movies: MovieData[];
};
