import LoginPresenter from "../features/auth/presenter/LoginPresenter";
import PageContainer from "../components/ui/PageContainer";
import type { LoginSchema } from "../features/auth/schema/login";

const LoginPage = () => {
  const handleLogin = (data: LoginSchema) => {
    console.log(data);
  };

  return (
    <PageContainer>
      <LoginPresenter onSubmit={handleLogin} />
    </PageContainer>
  );
};

export default LoginPage;
