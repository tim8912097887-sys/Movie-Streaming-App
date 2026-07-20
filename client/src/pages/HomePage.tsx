import { useEffect, useState } from "react";
import MoviePresenter from "../features/movie/presenter/MoviePresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import PaginationPresenter from "../shared/components/layout/pagination/PaginationPresenter";
import useFetch from "../shared/hooks/useFetch";
import ErrorFallback from "../shared/components/ui/ErrorFallback";
import { getMovies } from "../features/movie/api/request";

const HomePage = () => {
  const { handleFetch, status } = useFetch({
    fetchFunction: getMovies,
  });

  const [attemptedPage, setAttemptedPage] = useState(1);

  useEffect(() => {
    const fetchMoviesData = async () => {
      await handleFetch(1);
    };

    fetchMoviesData();
  }, []);

  const handleRetry = () => {
    handleFetch(attemptedPage);
  };

  if (status.error && !status.fetchedData && !status.isFetching) {
    return (
      <PageContainer customClass="justify-center items-center">
        <ErrorFallback
          message={status.error || "Failed to retrieve movies."}
          onRetry={handleRetry}
        />
      </PageContainer>
    );
  }

  return (
    <PageContainer customClass="flex-col gap-4">
      <MoviePresenter
        movies={status.fetchedData?.movies || []}
        isLoading={status.isFetching}
      />
      <PaginationPresenter
        setAttemptedPage={setAttemptedPage}
        totalPage={status.fetchedData?.total_page || 0}
        currentPage={status.fetchedData?.current_page || 0}
        onPageChange={handleFetch}
        isLoading={status.isFetching}
      />
    </PageContainer>
  );
};

export default HomePage;
