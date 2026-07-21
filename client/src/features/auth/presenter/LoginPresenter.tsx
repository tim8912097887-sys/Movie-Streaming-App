import Button from "../../../shared/components/ui/Button";
import { loginSchema, type LoginSchema } from "../schema/login";
import Form from "../ui/Form";
import Input from "../ui/Input";
import InputGroup from "../ui/InputGroup";
import { useFormData } from "../hook/useFormData";
import ErrorText from "../../../shared/components/ui/ErrorText";
import useFetch from "../../../shared/hooks/useFetch";
import { Spinner } from "../../../shared/components/ui/Spinner";
import { useNavigate } from "react-router";
import { toast } from "react-toastify";
import type { LoginData } from "../../../shared/schema/data/login";
import { useAuthStore } from "../store/store";
import { useEffect } from "react";

type LoginPresenterProps = {
  onSubmit: (data: LoginSchema) => Promise<LoginData>;
};

const LoginPresenter = ({ onSubmit }: LoginPresenterProps) => {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);
  const { formData, handleChange, errors } = useFormData<LoginSchema>({
    initialValues: {
      email: "",
      password: "",
    },
    schemaValidater: loginSchema,
  });

  const { handleFetch, status } = useFetch<LoginSchema, LoginData>({
    fetchFunction: onSubmit,
    validateSchema: loginSchema,
  });

  const formSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    await handleFetch(formData);
  };

  // Notify user login successfully and redirect to recommandation page
  useEffect(() => {
    if (status.isSuccess) {
      toast.success("Login successfully", {
        autoClose: 1500,
        position: "top-right",
        onClose: () => {
          navigate("/recommendations");
          login(status.fetchedData?.access_token || "");
        },
      });
    }
  }, [status.isSuccess]);

  return (
    <div className="w-70 md:w-100 rounded-xl border border-gray-200 bg-slate-300 p-8 shadow-lg md:p-12">
      <Form onSubmit={formSubmit}>
        <Form.Head>Login Movie Streaming</Form.Head>
        <Form.Content>
          <InputGroup name="email" label="Email">
            <Input
              inputProps={{
                type: "email",
                required: true,
                placeholder: "example@ex.com",
                name: "email",
                id: "email",
                onChange: handleChange,
              }}
            />
            {errors.email && <ErrorText>{errors.email}</ErrorText>}
          </InputGroup>
          <InputGroup name="password" label="Password">
            <Input
              inputProps={{
                type: "password",
                required: true,
                placeholder: "password",
                name: "password",
                id: "password",
                onChange: handleChange,
              }}
            />
            {errors.password && <ErrorText>{errors.password}</ErrorText>}
          </InputGroup>
        </Form.Content>
        <Form.Footer>
          <Button
            buttonProps={{
              type: "submit",
              disabled:
                Object.keys(errors).length > 0 ||
                status.isFetching ||
                status.isSuccess,
            }}
            size="md"
            color="primary"
            btnType="normal"
          >
            {status.isFetching ? <Spinner size="sm" color="white" /> : "Login"}
          </Button>
          {status.error && <ErrorText>{status.error}</ErrorText>}
        </Form.Footer>
      </Form>
    </div>
  );
};

export default LoginPresenter;
