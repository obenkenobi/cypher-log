import type { GatsbyConfig } from "gatsby";
import { createProxyMiddleware } from "http-proxy-middleware";

const config: GatsbyConfig = {
  siteMetadata: {
    title: `Cypher Log`,
    siteUrl: `https://www.cypherlog.com`
  },
  graphqlTypegen: true,
  plugins: ["gatsby-plugin-postcss"],

  pathPrefix: `/ui`,

  // Proxy to ui server
  developMiddleware: (app: any) => {
    app.use("/auth", createProxyMiddleware({
      target: "https://localhost:8080",
      secure: false, // Do not reject self-signed certificates.
      pathRewrite: {
        "/auth": "/auth",
      },
    }));
    app.use("/api", createProxyMiddleware({
      target: "https://localhost:8080",
      secure: false, // Do not reject self-signed certificates.
      pathRewrite: {
        "/api": "/api",
      },
    }));
  },
};

export default config;
