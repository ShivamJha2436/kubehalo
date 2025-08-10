type Props = {
    policy: {
        name: string;
        namespace: string;
        targetDeployment: string;
        metricQuery: string;
        thresholds: { scaleUp: number; scaleDown: number };
    };
};

export default function PolicyCard({ policy }: Props) {
    return (
        <div className="bg-white rounded-lg p-4 shadow hover:shadow-lg transition">
            <h3 className="text-lg font-semibold mb-1">{policy.name}</h3>
            <p><strong>Namespace:</strong> {policy.namespace}</p>
            <p><strong>Target:</strong> {policy.targetDeployment}</p>
            <p><strong>Metric:</strong> {policy.metricQuery}</p>
            <p><strong>Scale Up:</strong> {policy.thresholds.scaleUp}</p>
            <p><strong>Scale Down:</strong> {policy.thresholds.scaleDown}</p>
        </div>
    );
}
