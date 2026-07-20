import axios from "axios";

export const authApi = axios.create({
  baseURL: import.meta.env.VITE_API_URL + "/users",
  method: "POST",
  timeout: 3500,
  headers: {
    "Content-Type": "application/json",
  },
});
