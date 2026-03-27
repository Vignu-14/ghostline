import {
  createContext,
  useEffect,
  useState,
  type PropsWithChildren,
} from "react";
import * as authService from "../services/authService";
import type { LoginRequest, RegisterRequest, User } from "../types/user";

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  login: (input: LoginRequest) => Promise<void>;
  register: (input: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  async function refreshUser() {
    try {
      const response = await authService.getCurrentUser();
      setUser(response.user);
    } catch {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }

  async function handleLogin(input: LoginRequest) {
    const response = await authService.login(input);
    setUser(response.user);
  }

  async function handleRegister(input: RegisterRequest) {
    const response = await authService.register(input);
    setUser(response.user);
  }

  async function handleLogout() {
    await authService.logout();
    setUser(null);
  }

  useEffect(() => {
    void refreshUser();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        login: handleLogin,
        register: handleRegister,
        logout: handleLogout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
