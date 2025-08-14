"use client"

import MetricChart from "../../components/MetricChart";

export default function MetricsPage() {
    return (
        <div>
            <h1 className="text-2xl font-semibold mb-4">Metrics</h1>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-white p-4 rounded shadow">
                    <MetricChart metric="up" label="Up metric" />
                </div>
                <div className="bg-white p-4 rounded shadow">
                    <MetricChart metric="cpu_usage" label="CPU usage (sample)" />
                </div>
            </div>
        </div>
    );
}
