"use client";

import { Card } from "@/components/ui/Card";

type MetricCardProps = {
  title: string;
  value: string;
  description: string;
  className?: string;
};

export function MetricCard({
  title,
  value,
  description,
  className,
}: MetricCardProps) {
  return (
    <Card className={className}>
      <div className="space-y-1">
        <p className="text-xs text-slate-400">{title}</p>
        <p className="text-2xl font-semibold text-white">{value}</p>
        <p className="text-xs text-slate-400">{description}</p>
      </div>
    </Card>
  );
}
