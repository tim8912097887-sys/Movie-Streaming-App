import { useNavigate } from "react-router";
import Button from "../../../components/ui/Button";
import Form from "../ui/Form";
import GenreSelector from "../ui/GenreSelector";
import Input from "../ui/Input";
import InputGroup from "../ui/InputGroup";
import { registerSchema, type RegisterSchema } from "../schema/register";
import { useFormData } from "../hook/useFormData";
import ErrorText from "../../../components/ui/ErrorText";
import useFormSubmit from "../hook/useFormSubmit";
import { Spinner } from "../../../components/ui/Spinner";
import { useEffect } from "react";
import { toast } from "react-toastify";

const genres = [
  { genre_id: 1, name: "Comedy" },
  { genre_id: 2, name: "Drama" },
  { genre_id: 3, name: "Western" },
  { genre_id: 4, name: "Fantasy" },
  { genre_id: 5, name: "Thriller" },
  { genre_id: 6, name: "Sci-Fi" },
  { genre_id: 7, name: "Action" },
  { genre_id: 8, name: "Mystery" },
  { genre_id: 9, name: "Crime" },
];

type RegisterPresenterProps = {
  onSubmit: (data: RegisterSchema) => Promise<void>;
};

const RegisterPresenter = ({ onSubmit }: RegisterPresenterProps) => {
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

  const { handleSubmit, status } = useFormSubmit<RegisterSchema>({
    submitFunction: onSubmit,
    validateSchema: registerSchema,
  });

  const formSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    await handleSubmit(formData);
  };

  useEffect(() => {
    if (status.isSuccess) {
      toast.success("Register successfully", {
        autoClose: 1500,
        position: "top-right",
        onClose: () => navigate("/login"),
      });
    }
  }, [status.isSuccess]);

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
          <GenreSelector
            genres={genres}
            selectedGenres={formData.favorite_genres}
            onChange={setGenres}
          />
          {errors.favorite_genres && (
            <ErrorText>{errors.favorite_genres}</ErrorText>
          )}
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
