import { movieApi } from "./axios";
import { axiosWrapper } from "../../../shared/api/axios-wrapper";
import type { FetchMoviesData } from "../../../shared/schema/data/movies";

export const getMovies = axiosWrapper<number, FetchMoviesData>(async (page) => {
  const response = await movieApi.get("", { params: { page, limit: 8 } });
  return response.data.data;
});

export const getRecommendedMovies = axiosWrapper<string, FetchMoviesData>(
  async (token) => {
    const response = await movieApi.get("/user", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.data;
  },
);
