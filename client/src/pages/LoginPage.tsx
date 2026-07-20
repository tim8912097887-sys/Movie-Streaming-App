import LoginPresenter from "../features/auth/presenter/LoginPresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import { loginUser } from "../features/auth/api/request";

const LoginPage = () => {
  return (
    <PageContainer>
      <LoginPresenter onSubmit={loginUser} />
    </PageContainer>
  );
};

export default LoginPage;
