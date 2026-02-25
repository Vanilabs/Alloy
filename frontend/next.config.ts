import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,

  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "https://3366-105-112-238-79.ngrok-free.app/api/:path*",
      },
    ];
  },
};

export default nextConfig;
