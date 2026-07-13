import { createBrowserRouter } from "react-router";
import RootLayout from "../layouts/RootLayout";

const router = createBrowserRouter([
  {
    element: <RootLayout />,
    children: [
      {
        path: "/",
        element: <div>Home</div>,
      },
      {
        path: "/recommendations",
        element: <div>Recommendations</div>,
      },
      {
        path: "/login",
        element: <div>Login</div>,
      },
      {
        path: "/register",
        element: <div>Register</div>,
      },
    ],
  },
]);

export default router;
