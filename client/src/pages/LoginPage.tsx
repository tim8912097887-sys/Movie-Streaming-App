import LoginPresenter from "../features/auth/presenter/LoginPresenter";
import PageContainer from "../shared/components/ui/PageContainer";
import type { LoginSchema } from "../features/auth/schema/login";

const LoginPage = () => {
  const handleLogin = async (data: LoginSchema) => {
    await new Promise((resolve) => setTimeout(resolve, 3000));
    console.log(data);
  };

  return (
    <PageContainer>
      <LoginPresenter onSubmit={handleLogin} />
    </PageContainer>
  );
};

export default LoginPage;
