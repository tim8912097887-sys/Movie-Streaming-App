import RegisterPresenter from "../features/auth/presenter/RegisterPresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import { registerUser } from "../features/auth/api/request";

const RegisterPage = () => {
  return (
    <PageContainer>
      <RegisterPresenter onSubmit={registerUser} />
    </PageContainer>
  );
};

export default RegisterPage;
