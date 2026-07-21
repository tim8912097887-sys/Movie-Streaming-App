import RegisterPresenter from "../features/auth/presenter/RegisterPresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import { getGenres, registerUser } from "../features/auth/api/request";
import useFetch from "../shared/hooks/useFetch";
import { useEffect } from "react";

const RegisterPage = () => {
  const { handleFetch, status } = useFetch({
    fetchFunction: getGenres,
  });

  useEffect(() => {
    const fetchGenres = async () => {
      await handleFetch();
    };

    fetchGenres();
  }, []);

  return (
    <PageContainer>
      <RegisterPresenter
        onSubmit={registerUser}
        genres={status.fetchedData || []}
        isFetching={status.isFetching}
        fetchingError={status.error}
      />
    </PageContainer>
  );
};

export default RegisterPage;
