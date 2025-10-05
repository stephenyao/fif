import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "./AuthContext";
import type { ReactNode } from "react";

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const { loggedIn, initializing } = useAuth();
  const location = useLocation();

  if (initializing) return null;

  if (!loggedIn) {
    return <Navigate to="/" replace state={{ from: location }} />;
  }

  return <>{children}</>;
}
