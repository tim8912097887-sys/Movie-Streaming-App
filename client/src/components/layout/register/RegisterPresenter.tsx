import Button from "../../ui/Button";
import Form from "../../ui/Form";
import Input from "../../ui/Input";
import InputGroup from "../../ui/InputGroup";

const RegisterPresenter = () => {
  return (
    <div className="w-70 md:w-100 rounded-xl border border-gray-200 bg-slate-300 p-8 shadow-lg md:p-12">
      <Form onSubmit={() => {}}>
        <Form.Head>Register Movie Streaming</Form.Head>
        <Form.Content>
          <InputGroup name="name" label="Name">
            <Input
              inputProps={{
                type: "text",
                required: true,
                placeholder: "name",
                name: "name",
                id: "name",
              }}
            />
          </InputGroup>
          <InputGroup name="email" label="Email">
            <Input
              inputProps={{
                type: "email",
                required: true,
                placeholder: "example@ex.com",
                name: "email",
                id: "email",
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
            Register
          </Button>
        </Form.Footer>
      </Form>
    </div>
  );
};

export default RegisterPresenter;
