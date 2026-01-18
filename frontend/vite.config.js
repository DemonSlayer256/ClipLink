import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Proxy all API calls to Go backend in dev
      "/login": "http://localhost:8080",
      "/register": "http://localhost:8080",
      "/links": "http://localhost:8080",
      "/shorten": "http://localhost:8080",
      "/delete": "http://localhost:8080",
    },
  },
});
