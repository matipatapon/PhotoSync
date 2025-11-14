import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import fs from 'fs';

const CERT = "../localhost.crt"
const KEY = "../localhost.key"

export default defineConfig({
    server:{
        https:{
            key: fs.readFileSync(KEY),
            cert: fs.readFileSync(CERT),
        }
    },
  plugins: [reactRouter(), tailwindcss()],
});
