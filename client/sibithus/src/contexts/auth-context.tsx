import { createContext, useContext, useMemo } from 'react';
import { createAuthClient } from 'better-auth/react';
import { usernameClient } from 'better-auth/client/plugins';

type AuthClient = ReturnType<typeof createAuthClient<{ plugins: [ReturnType<typeof usernameClient>] }>>;

// biome-ignore lint/suspicious/noExplicitAny: <Will be defined when we call this>
const AuthContext = createContext<AuthClient>(undefined as any);

const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const authClient = useMemo(() => {
    return createAuthClient({
      baseURL: import.meta.env.VITE_AUTH_BASE_URL,
      plugins: [
        usernameClient()
      ],
    });
  }, []);

  return (
    <AuthContext value={authClient}>
      {children}
    </AuthContext>
  );
}

const useAuth = () => {
  return useContext(AuthContext);
}

export {
  AuthProvider,
  useAuth,
  AuthContext,
}