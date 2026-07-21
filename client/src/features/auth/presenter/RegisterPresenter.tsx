import { useNavigate } from "react-router";
import Button from "../../../shared/components/ui/Button";
import Form from "../ui/Form";
import GenreSelector from "../ui/GenreSelector";
import Input from "../ui/Input";
import InputGroup from "../ui/InputGroup";
import { registerSchema, type RegisterSchema } from "../schema/register";
import { useFormData } from "../hook/useFormData";
import ErrorText from "../../../shared/components/ui/ErrorText";
import useFetch from "../../../shared/hooks/useFetch";
import { Spinner } from "../../../shared/components/ui/Spinner";
import { toast } from "react-toastify";
import type { RegisterData } from "../../../shared/schema/data/register";
import type { GenresData } from "../../../shared/schema/data/genres";

type RegisterPresenterProps = {
  onSubmit: (data: RegisterSchema) => Promise<RegisterData>;
  genres: GenresData;
  isFetching: boolean;
  fetchingError: string | null;
};

const RegisterPresenter = ({
  onSubmit,
  genres,
  isFetching,
  fetchingError,
}: RegisterPresenterProps) => {
  const navigate = useNavigate();
  const { formData, handleChange, errors } = useFormData<RegisterSchema>({
    initialValues: {
      name: "",
      email: "",
      password: "",
      favorite_genres: [],
    },
    schemaValidater: registerSchema,
  });

  const setGenres = (genres: RegisterSchema["favorite_genres"]) => {
    handleChange({ target: { name: "favorite_genres", value: genres } } as any);
  };

  const { handleFetch, status } = useFetch<RegisterSchema, RegisterData>({
    fetchFunction: onSubmit,
    validateSchema: registerSchema,
  });

  const formSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    await handleFetch(formData);
  };

  if (status.isSuccess) {
    toast.success("Register successfully", {
      autoClose: 1500,
      position: "top-right",
      onClose: () => navigate("/login"),
    });
  }

  return (
    <div className="w-70 md:w-100 rounded-xl border border-gray-200 bg-slate-300 p-8 shadow-lg md:p-12">
      <Form onSubmit={formSubmit}>
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
                onChange: handleChange,
              }}
            />
            {errors.name && <ErrorText>{errors.name}</ErrorText>}
          </InputGroup>
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
          {isFetching ? (
            <Spinner />
          ) : (
            <GenreSelector
              genres={genres}
              selectedGenres={formData.favorite_genres}
              onChange={setGenres}
            />
          )}
          {fetchingError && <ErrorText>{fetchingError}</ErrorText>}
          {errors.favorite_genres && (
            <ErrorText>{errors.favorite_genres}</ErrorText>
          )}
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
            {status.isFetching ? (
              <Spinner size="sm" color="primary" />
            ) : (
              "Register"
            )}
          </Button>
          {status.error && <ErrorText>{status.error}</ErrorText>}
        </Form.Footer>
      </Form>
    </div>
  );
};

export default RegisterPresenter;
