import axios from "axios";

export const movieApi = axios.create({
  baseURL: import.meta.env.VITE_API_URL + "/movies",
  method: "GET",
  timeout: 3500,
});
