// electron.vite.config.ts
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react-swc";
import { defineConfig } from "electron-vite";
import flow from "rollup-plugin-flow";
import { ViteImageOptimizer } from "vite-plugin-image-optimizer";
import { viteStaticCopy } from "vite-plugin-static-copy";
import svgr from "vite-plugin-svgr";
import tsconfigPaths from "vite-tsconfig-paths";
var electron_vite_config_default = defineConfig({
  main: {
    root: ".",
    build: {
      rollupOptions: {
        input: {
          index: "emain/emain.ts"
        }
      },
      outDir: "dist/main"
    },
    plugins: [tsconfigPaths(), flow()],
    resolve: {
      alias: {
        "@": "frontend"
      }
    },
    define: {
      "process.env.WS_NO_BUFFER_UTIL": "true",
      "process.env.WS_NO_UTF_8_VALIDATE": "true"
    }
  },
  preload: {
    root: ".",
    build: {
      sourcemap: true,
      rollupOptions: {
        input: {
          index: "emain/preload.ts",
          "preload-webview": "emain/preload-webview.ts"
        },
        output: {
          format: "cjs"
        }
      },
      outDir: "dist/preload"
    },
    plugins: [tsconfigPaths(), flow()]
  },
  renderer: {
    root: ".",
    build: {
      target: "es6",
      sourcemap: true,
      outDir: "dist/frontend",
      rollupOptions: {
        input: {
          index: "index.html"
        }
      }
    },
    server: {
      open: false
    },
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ["mixed-decls"]
        }
      }
    },
    plugins: [
      ViteImageOptimizer(),
      tsconfigPaths(),
      flow(),
      svgr({
        svgrOptions: { exportType: "default", ref: true, svgo: false, titleProp: true },
        include: "**/*.svg"
      }),
      react({}),
      tailwindcss(),
      viteStaticCopy({
        targets: [
          { src: "node_modules/monaco-editor/min/vs/*", dest: "monaco" },
          { src: "public/fontawesome/webfonts/*", dest: "fontawesome/webfonts" },
          { src: "public/icons/*", dest: "icons" }
        ]
      })
    ]
  }
});
export {
  electron_vite_config_default as default
};
