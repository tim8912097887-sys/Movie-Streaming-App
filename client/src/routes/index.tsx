import { createBrowserRouter } from "react-router";
import RootLayout from "../shared/layouts/RootLayout";
import LoginPage from "../pages/LoginPage";
import RegisterPage from "../pages/RegisterPage";
import RecommandationPage from "../pages/RecommandationPage";
import HomePage from "../pages/HomePage";
import AuthGuard from "./guards/AuthGuard";

const router = createBrowserRouter([
  {
    element: <RootLayout />,
    children: [
      {
        path: "/",
        element: <HomePage />,
      },
      {
        path: "/recommendations",
        element: (
          <AuthGuard>
            <RecommandationPage />
          </AuthGuard>
        ),
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
