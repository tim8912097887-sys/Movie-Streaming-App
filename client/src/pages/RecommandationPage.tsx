import { useEffect } from "react";
import PageContainer from "../shared/components/ui/PageContainer";
import useFetch from "../shared/hooks/useFetch";
import MoviePresenter from "../features/movie/presenter/MoviePresenter";
import { getRecommendedMovies } from "../features/movie/api/request";
import { useAuthStore } from "../features/auth/store/store";
import ErrorFallback from "../shared/components/ui/ErrorFallback";

const RecommandationPage = () => {
  const { handleFetch, status } = useFetch({
    fetchFunction: getRecommendedMovies,
  });
  const token = useAuthStore((state) => state.token);

  useEffect(() => {
    const fetchMoviesData = async () => {
      await handleFetch(token!);
    };

    fetchMoviesData();
  }, []);

  // Handle error retry
  const handleRetry = () => {
    handleFetch(token!);
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
    <PageContainer>
      <MoviePresenter
        movies={status.fetchedData?.movies || []}
        isLoading={status.isFetching}
      />
    </PageContainer>
  );
};

export default RecommandationPage;
