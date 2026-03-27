import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { HomePage } from "./pages/HomePage";
import { ChatPage } from "./pages/ChatPage";
import { AdminPage } from "./pages/AdminPage";
import { ProfilePage } from "./pages/ProfilePage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { Navbar } from "./components/layout/Navbar";
import { MobileNav } from "./components/layout/MobileNav";
import { Footer } from "./components/layout/Footer";
import { ProtectedRoute } from "./components/common/ProtectedRoute";

function AppShell() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route element={<HomePage />} path="/" />
        <Route element={<LoginPage />} path="/login" />
        <Route element={<RegisterPage />} path="/register" />
        <Route element={<ProfilePage />} path="/u/:username" />
        <Route element={<ProtectedRoute />}>
          <Route element={<ChatPage />} path="/chat" />
          <Route element={<ProfilePage />} path="/profile" />
        </Route>
        <Route element={<ProtectedRoute requireAdmin />}>
          <Route element={<AdminPage />} path="/admin" />
        </Route>
        <Route element={<Navigate replace to="/" />} path="/home" />
        <Route element={<NotFoundPage />} path="*" />
      </Routes>
      <Footer />
      <MobileNav />
    </>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppShell />
    </BrowserRouter>
  );
}
