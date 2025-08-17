import { ScalePolicy } from "../types";

// Simple in-memory mock store
let store: ScalePolicy[] = [
    {
        apiVersion: "kubehalo.sh/v1",
        kind: "ScalePolicy",
        metadata: { name: "example-policy", namespace: "default" },
        spec: {
            targetRef: { kind: "Deployment", name: "test-deployment", namespace: "default" },
            metric: { name: "up-metric", query: "up", threshold: 0 },
            minReplicas: 1,
            maxReplicas: 5,
            scaleUp: { step: 1 },
            scaleDown: { step: 1 },
        },
    },
];

export async function fetchPolicies(): Promise<ScalePolicy[]> {
    await new Promise((r) => setTimeout(r, 300));
    return store;
}

export async function getPolicy(name: string): Promise<ScalePolicy | null> {
    await new Promise((r) => setTimeout(r, 200));
    return store.find((s) => s.metadata.name === name) ?? null;
}

export async function createPolicy(data: ScalePolicy | any): Promise<ScalePolicy> {
    const p: ScalePolicy = {
        apiVersion: "kubehalo.sh/v1",
        kind: "ScalePolicy",
        metadata: { name: data.metadata.name, namespace: "default" },
        spec: data.spec,
    };
    store = [p, ...store];
    await new Promise((r) => setTimeout(r, 200));
    return p;
}

export async function updatePolicy(name: string, data: any): Promise<ScalePolicy | null> {
    const idx = store.findIndex((s) => s.metadata.name === name);
    if (idx === -1) return null;
    store[idx] = { ...store[idx], spec: data.spec, metadata: { ...store[idx].metadata, ...data.metadata } };
    await new Promise((r) => setTimeout(r, 200));
    return store[idx];
}

export async function getMockOverview() {
    await new Promise((r) => setTimeout(r, 100));
    return { policies: store.length, deployments: 3, alerts: 0 };
}
