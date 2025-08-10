export interface ScalePolicy {
    name: string;
    namespace: string;
    targetDeployment: string;
    targetNamespace: string;
    metricQuery: string;
    thresholds: {
        scaleUp: number;
        scaleDown: number;
    };
    status?: string;
}
