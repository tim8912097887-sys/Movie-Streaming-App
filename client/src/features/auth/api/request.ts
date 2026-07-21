import { axiosWrapper } from "../../../shared/api/axios-wrapper";
import type { GenresData } from "../../../shared/schema/data/genres";
import type { LoginData } from "../../../shared/schema/data/login";
import type { RegisterData } from "../../../shared/schema/data/register";
import type { LoginSchema } from "../schema/login";
import type { RegisterSchema } from "../schema/register";
import { authApi } from "./axios";

export const registerUser = axiosWrapper<RegisterSchema, RegisterData>(
  async (data: RegisterSchema) => {
    const response = await authApi.post("/register", data);
    return response.data.data;
  },
);

export const loginUser = axiosWrapper<LoginSchema, LoginData>(
  async (data: LoginSchema) => {
    const response = await authApi.post("/login", data, {
      withCredentials: true,
    });
    return response.data.data;
  },
);

export const logoutUser = axiosWrapper<void, void>(async () => {
  await authApi.post(
    "/logout",
    {},
    {
      withCredentials: true,
    },
  );
});

export const getGenres = axiosWrapper<void, GenresData>(async () => {
  const response = await authApi.get("/genres");
  return response.data.data;
});
