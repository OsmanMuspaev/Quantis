import { Routes, Route } from "react-router-dom";
import MainLayout from "./layouts/MainLayout";

import LoginPage from "./pages/LoginPage";
import HomePage from "./pages/HomePage";
import VerifyPage from "./pages/VerifyPage";
import ProfilePage from "./pages/ProfilePage";
import MyTestsPage from "./pages/MyTestsPage";
import TestResultPage from "./pages/TestResultPage";
import ForbiddenPage from "./pages/ForbiddenPage";
import NotFoundPage from "./pages/NotFoundPage";

import ProtectedRoute from "./router/ProtectedRoute";
import PublicRoute from "./router/PublicRoute";

function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <PublicRoute>
            <HomePage />
          </PublicRoute>
        }
      />

      <Route
        path="/login"
        element={
          <PublicRoute>
            <LoginPage />
          </PublicRoute>
        }
      />

      <Route
        path="/verify"
        element={
          <PublicRoute>
            <VerifyPage />
          </PublicRoute>
        }
      />

      <Route path="/forbidden" element={<ForbiddenPage />} />

      <Route element={<MainLayout />}>
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/tests/:id"
          element={
            <ProtectedRoute>
              <MyTestsPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/results/:id"
          element={
            <ProtectedRoute>
              <TestResultPage />
            </ProtectedRoute>
          }
        />
      </Route>

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;
