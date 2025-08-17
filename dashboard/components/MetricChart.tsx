"use client"

import { useEffect, useState } from "react";

/**
 * Minimal mock chart using SVG lines — lightweight and dependency-free.
 * Replace with chart.js/recharts if you prefer.
 */

type Props = { metric: string; label?: string };

export default function MetricChart({ metric, label }: Props) {
    const [points, setPoints] = useState<number[]>([]);

    useEffect(() => {
        // mock data simulator
        const initial = Array.from({ length: 20 }).map(() => Math.random() * 1.2);
        setPoints(initial);

        const t = setInterval(() => {
            setPoints((p) => [...p.slice(1), Math.max(0, p[p.length - 1] + (Math.random() - 0.5) * 0.2)]);
        }, 2000);

        return () => clearInterval(t);
    }, [metric]);

    const max = Math.max(...points, 1);
    const width = 600;
    const height = 120;
    const step = width / (points.length - 1 || 1);

    return (
        <div>
            <div className="flex items-center justify-between mb-2">
                <div className="font-medium">{label ?? metric}</div>
                <div className="text-sm text-gray-500">live</div>
            </div>
            <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-32">
                <polyline
                    fill="none"
                    stroke="#4f46e5"
                    strokeWidth={2}
                    points={points.map((v, i) => `${i * step},${height - (v / max) * (height - 10)}`).join(" ")}
                />
            </svg>
        </div>
    );
}
