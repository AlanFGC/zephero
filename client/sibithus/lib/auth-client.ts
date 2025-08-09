import { createAuthClient } from "better-auth/react";
export const authClient = createAuthClient({
  baseURL: process.env.VITE_AUTH_BASE_URL || "/",
});
