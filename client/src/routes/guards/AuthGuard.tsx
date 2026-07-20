import { Navigate } from "react-router";
import { useAuthStore } from "../../features/auth/store/store";
import type { PropsWithChildren } from "react";

const AuthGuard = ({ children }: PropsWithChildren) => {
  const token = useAuthStore((state) => state.token);

  // If no token, redirect to login
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default AuthGuard;
