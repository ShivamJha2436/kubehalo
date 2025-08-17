export type ScalePolicy = {
    apiVersion: string;
    kind: string;
    metadata: {
        name: string;
        namespace?: string;
    };
    spec: {
        targetRef: { kind: string; name: string; namespace?: string };
        metric: { name: string; query: string; threshold: number };
        minReplicas: number;
        maxReplicas: number;
        scaleUp: { step: number };
        scaleDown: { step: number };
    };
};
