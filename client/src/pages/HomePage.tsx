import { useEffect } from "react";
import MoviePresenter, {
  type Movie,
} from "../features/movie/presenter/MoviePresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import PaginationPresenter from "../shared/components/layout/pagination/PaginationPresenter";
import useFetch from "../shared/hooks/useFetch";

const movies: Movie[] = [
  {
    title: "Neon Horizon",
    description:
      "A cybernetic detective uncovers a conspiracy that shakes the foundations of a rain-soaked, futuristic metropolis.",
    image: "https://images.unsplash.com/photo-1578301978693-85fa9c0320b9",
    rating: 8.4,
  },
  {
    title: "Whispers of the Oak",
    description:
      "Two siblings discover an ancient, enchanted forest behind their grandmother's house that holds secrets of the past.",
    image: "https://images.unsplash.com/photo-1511497584788-876760111969",
    rating: 7.9,
  },
  {
    title: "The Last Alchemist",
    description:
      "In an alternate 19th century, a rogue scholar races against time to find the elixir of life before an empire claims it.",
    image: "https://images.unsplash.com/photo-1516979187457-637abb4f9353",
    rating: 8.1,
  },
  {
    title: "Neon Horizon",
    description:
      "A cybernetic detective uncovers a conspiracy that shakes the foundations of a rain-soaked, futuristic metropolis.",
    image: "https://images.unsplash.com/photo-1578301978693-85fa9c0320b9",
    rating: 8.4,
  },
  {
    title: "Whispers of the Oak",
    description:
      "Two siblings discover an ancient, enchanted forest behind their grandmother's house that holds secrets of the past.",
    image: "https://images.unsplash.com/photo-1511497584788-876760111969",
    rating: 7.9,
  },
  {
    title: "The Last Alchemist",
    description:
      "In an alternate 19th century, a rogue scholar races against time to find the elixir of life before an empire claims it.",
    image: "https://images.unsplash.com/photo-1516979187457-637abb4f9353",
    rating: 8.1,
  },
  {
    title: "Whispers of the Oak",
    description:
      "Two siblings discover an ancient, enchanted forest behind their grandmother's house that holds secrets of the past.",
    image: "https://images.unsplash.com/photo-1511497584788-876760111969",
    rating: 7.9,
  },
  {
    title: "The Last Alchemist",
    description:
      "In an alternate 19th century, a rogue scholar races against time to find the elixir of life before an empire claims it.",
    image: "https://images.unsplash.com/photo-1516979187457-637abb4f9353",
    rating: 8.1,
  },
];

// type FetchMoviesData = {
//   totalPage: number;
//   currentPage: number;
//   data: Movie[];
// };

const fetchMovies = async (page: number) => {
  await new Promise((resolve) => setTimeout(resolve, 3000));
  const response = {
    totalPage: 20,
    currentPage: page,
    data: movies,
  };
  return response;
};

const HomePage = () => {
  const { handleFetch, status } = useFetch({
    fetchFunction: fetchMovies,
  });

  useEffect(() => {
    const fetchMoviesData = async () => {
      await handleFetch(1);
    };

    fetchMoviesData();
  }, []);

  return (
    <PageContainer customClass="flex-col gap-4">
      <MoviePresenter
        movies={status.fetchedData?.data || []}
        isLoading={status.isFetching}
      />
      <PaginationPresenter
        totalPage={status.fetchedData?.totalPage || 0}
        currentPage={status.fetchedData?.currentPage || 0}
        onPageChange={handleFetch}
        isLoading={status.isFetching}
      />
    </PageContainer>
  );
};

export default HomePage;
