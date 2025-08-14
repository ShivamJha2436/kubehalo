import Card from "../components/ui/Card";
import MetricChart from "../components/MetricChart";
import { getMockOverview } from "../lib/api";

export default async function Page() {
    // server component can call server-side helpers; using a mock for now
    const overview = await getMockOverview();

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-semibold">Overview</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card title="Scale Policies" value={String(overview.policies)} />
                <Card title="Total Deployments" value={String(overview.deployments)} />
                <Card title="Active Alerts" value={String(overview.alerts)} />
            </div>

            <div className="mt-4">
                <h2 className="text-lg font-medium mb-2">Cluster Metrics</h2>
                <div className="bg-white rounded-md p-4 shadow">
                    <MetricChart metric="up" label="Up (sample metric)" />
                </div>
            </div>
        </div>
    );
}
