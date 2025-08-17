"use client"

import React, { useState } from "react";
import Input from "./ui/Input";
import Select from "./ui/Select";
import Button from "./ui/Button";
import { ScalePolicy } from "../types";

type Props = {
    initialValues?: ScalePolicy | null;
    onSubmit: (data: any) => Promise<void>;
};

export default function PolicyForm({ initialValues, onSubmit }: Props) {
    const [name, setName] = useState(initialValues?.metadata.name ?? "");
    const [deployment, setDeployment] = useState(initialValues?.spec.targetRef.name ?? "");
    const [query, setQuery] = useState(initialValues?.spec.metric.query ?? "up");
    const [threshold, setThreshold] = useState(String(initialValues?.spec.metric.threshold ?? 0));
    const [minReplicas, setMinReplicas] = useState(String(initialValues?.spec.minReplicas ?? 1));
    const [maxReplicas, setMaxReplicas] = useState(String(initialValues?.spec.maxReplicas ?? 5));

    async function submit(e: React.FormEvent) {
        e.preventDefault();
        await onSubmit({
            metadata: { name },
            spec: {
                targetRef: { kind: "Deployment", name: deployment, namespace: "default" },
                metric: { name: "metric", query, threshold: Number(threshold) },
                minReplicas: Number(minReplicas),
                maxReplicas: Number(maxReplicas),
                scaleUp: { step: 1 },
                scaleDown: { step: 1 },
            },
        });
    }

    return (
        <form onSubmit={submit} className="space-y-4">
            <Input label="Policy name" value={name} onChange={(e) => setName(e.target.value)} />
            <Input label="Target deployment" value={deployment} onChange={(e) => setDeployment(e.target.value)} />
            <Input label="PromQL query" value={query} onChange={(e) => setQuery(e.target.value)} />
            <Input label="Threshold" value={threshold} onChange={(e) => setThreshold(e.target.value)} />
            <div className="grid grid-cols-2 gap-3">
                <Input label="Min replicas" value={minReplicas} onChange={(e) => setMinReplicas(e.target.value)} />
                <Input label="Max replicas" value={maxReplicas} onChange={(e) => setMaxReplicas(e.target.value)} />
            </div>
            <div className="flex items-center justify-end">
                <Button type="submit">Save policy</Button>
            </div>
        </form>
    );
}
