"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [client] = useState(() => queryClient);

  return (
    <html lang="en">
      <body suppressHydrationWarning>
        <QueryClientProvider client={client}>{children}</QueryClientProvider>
      </body>
    </html>
  );
}
