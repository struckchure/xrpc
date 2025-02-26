import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "xRPC",
  description: "RPC Framework for Golang",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [{ text: "Guide", link: "/getting-started" }],

    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "What is xRPC?", link: "/what-is-xrpc" },
          { text: "Features", link: "/features" },
          { text: "Architecture", link: "/architecture" },
          { text: "Comparison", link: "/comparison" },
        ],
      },
      {
        text: "Guide",
        items: [
          { text: "Getting Started", link: "/getting-started" },
          { text: "Installation", link: "/installation" },
          { text: "Server", link: "/server" },
          { text: "Client", link: "/client" },
          { text: "Middleware", link: "/middleware" },
          { text: "Error Handling", link: "/error-handling" },
          { text: "Logging", link: "/logging" },
          { text: "Testing", link: "/testing" },
          { text: "Benchmarks", link: "/benchmarks" },
          { text: "FAQ", link: "/faq" },
        ],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/struckchure/xrpc" },
      { icon: "x", link: "https://x.com/struckchure" },
    ],
  },
});
