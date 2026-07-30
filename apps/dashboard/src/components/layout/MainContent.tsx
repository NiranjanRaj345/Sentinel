"use client";

import { ReactNode } from "react";

type MainContentProps = {
  children: ReactNode;
};

export function MainContent({ children }: MainContentProps) {
  return (
    <main className="min-h-screen overflow-auto bg-slate-950 p-6">
      <div className="mx-auto max-w-6xl space-y-6">{children}</div>
    </main>
  );
}
