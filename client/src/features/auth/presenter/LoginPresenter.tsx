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
import { useEffect } from "react";
import { toast } from "react-toastify";

type LoginPresenterProps = {
  onSubmit: (data: LoginSchema) => Promise<void>;
};

const LoginPresenter = ({ onSubmit }: LoginPresenterProps) => {
  const navigate = useNavigate();
  const { formData, handleChange, errors } = useFormData<LoginSchema>({
    initialValues: {
      email: "",
      password: "",
    },
    schemaValidater: loginSchema,
  });

  const { handleFetch, status } = useFetch<LoginSchema, void>({
    fetchFunction: onSubmit,
    validateSchema: loginSchema,
  });

  const formSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    await handleFetch(formData);
  };

  useEffect(() => {
    if (status.isSuccess) {
      toast.success("Login successfully", {
        autoClose: 1500,
        position: "top-right",
        onClose: () => navigate("/recommendations"),
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
              disabled: Object.keys(errors).length > 0 || status.isFetching,
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
