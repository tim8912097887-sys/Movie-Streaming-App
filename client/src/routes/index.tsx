import { createBrowserRouter } from "react-router";
import RootLayout from "../layouts/RootLayout";
import LoginPage from "../pages/LoginPage";
import RegisterPage from "../pages/RegisterPage";
import RecommandationPage from "../pages/RecommandationPage";

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
        element: <RecommandationPage />,
      },
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/register",
        element: <RegisterPage />,
      },
    ],
  },
]);

export default router;
