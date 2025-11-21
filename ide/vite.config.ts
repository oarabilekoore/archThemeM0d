import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [react(), tailwindcss(), tsconfigPaths()],
  resolve: {
    alias: {
      "@": new URL("./", import.meta.url).pathname,
    },
  },
  server: {
    // This makes the connection work!
    proxy: {
      // Any request starting with /files, /read, or /update goes to Go
      "/files": "http://localhost:8080",
      "/read": "http://localhost:8080",
      "/update": "http://localhost:8080",
    },
  },
});
