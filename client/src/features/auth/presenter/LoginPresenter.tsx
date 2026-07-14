import { useState } from "react";
import Button from "../../../components/ui/Button";
import type { LoginSchema } from "../schema/login";
import Form from "../ui/Form";
import Input from "../ui/Input";
import InputGroup from "../ui/InputGroup";

type LoginPresenterProps = {
  onSubmit: (data: LoginSchema) => void;
};

const LoginPresenter = ({ onSubmit }: LoginPresenterProps) => {
  const [formData, setFormData] = useState<LoginSchema>({
    email: "",
    password: "",
  });

  const handleSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className="w-70 md:w-100 rounded-xl border border-gray-200 bg-slate-300 p-8 shadow-lg md:p-12">
      <Form onSubmit={handleSubmit}>
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
                onChange: (e) => {
                  setFormData({ ...formData, email: e.target.value });
                },
              }}
            />
          </InputGroup>
          <InputGroup name="password" label="Password">
            <Input
              inputProps={{
                type: "password",
                required: true,
                placeholder: "password",
                name: "password",
                id: "password",
                onChange: (e) => {
                  setFormData({ ...formData, password: e.target.value });
                },
              }}
            />
          </InputGroup>
        </Form.Content>
        <Form.Footer>
          <Button
            buttonProps={{ type: "submit" }}
            size="md"
            color="primary"
            btnType="normal"
          >
            Login
          </Button>
        </Form.Footer>
      </Form>
    </div>
  );
};

export default LoginPresenter;
