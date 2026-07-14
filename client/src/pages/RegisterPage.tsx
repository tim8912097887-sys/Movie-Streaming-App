import RegisterPresenter from "../features/auth/presenter/RegisterPresenter";
import PageContainer from "../components/ui/PageContainer";
import type { RegisterSchema } from "../features/auth/schema/register";

const RegisterPage = () => {
  const handleRegister = async (data: RegisterSchema) => {
    await new Promise((resolve) => setTimeout(resolve, 3000));
    console.log(data);
  };

  return (
    <PageContainer>
      <RegisterPresenter onSubmit={handleRegister} />
    </PageContainer>
  );
};

export default RegisterPage;
