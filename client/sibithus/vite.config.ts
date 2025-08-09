// biome-ignore assist/source/organizeImports: <not much to do here>
import path from "path";
import react from "@vitejs/plugin-react-swc";
import { defineConfig, loadEnv } from 'vite'
import tailwindcss from "@tailwindcss/vite"

// https://vite.dev/config/
export default defineConfig(({ mode }) => {

  const env = loadEnv(mode, process.cwd(), '');
  return {
    base: env.VITE_BASE_URL || '/',
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
  }
});
