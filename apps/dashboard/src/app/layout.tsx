"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { Toast } from "@/components/toast/Toast";

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
        <QueryClientProvider client={client}>
          {children}
          <Toast />
        </QueryClientProvider>
      </body>
    </html>
  );
}
