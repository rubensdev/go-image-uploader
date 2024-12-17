import { defineConfig } from "vite";
import path from "node:path";

export default defineConfig({
  server: {
    origin: "",
  },
  build: {
    manifest: true,
    outDir: path.resolve(__dirname, "../dist"),
    // rollupOptions: {
    //   input: {
    //     index: "./src/main.js",
    //     about: "./src/about/about.js",
    //   },
    // },
  },
});
