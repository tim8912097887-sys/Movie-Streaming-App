import Button from "../../../components/ui/Button";
import { loginSchema, type LoginSchema } from "../schema/login";
import Form from "../ui/Form";
import Input from "../ui/Input";
import InputGroup from "../ui/InputGroup";
import { useFormData } from "../hook/useFormData";
import ErrorText from "../../../components/ui/ErrorText";
import useFormSubmit from "../hook/useFormSubmit";
import { Spinner } from "../../../components/ui/Spinner";
import { useNavigate } from "react-router";
import { useEffect } from "react";

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

  const { handleSubmit, status } = useFormSubmit<LoginSchema>({
    submitFunction: onSubmit,
    validateSchema: loginSchema,
  });

  const formSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    await handleSubmit(formData);
  };

  useEffect(() => {
    if (status.isSuccess) navigate("/recommendations");
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
              disabled: Object.keys(errors).length > 0 || status.isSubmitting,
            }}
            size="md"
            color="primary"
            btnType="normal"
          >
            {status.isSubmitting ? (
              <Spinner size="sm" color="white" />
            ) : (
              "Login"
            )}
          </Button>
          {status.error && <ErrorText>{status.error}</ErrorText>}
        </Form.Footer>
      </Form>
    </div>
  );
};

export default LoginPresenter;
