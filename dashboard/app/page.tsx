"use client"

import { useEffect, useState } from "react";
import PolicyCard from "@/components/PolicyCard";
import LoadingSpinner from "@/components/LoadingSpinner";
import Chart from "@/components/Chart";

const mockPolicies = [
    {
        name: "autoscale-nginx",
        namespace: "default",
        targetDeployment: "nginx-deployment",
        metricQuery: "avg(rate(http_requests_total[2m]))",
        thresholds: { scaleUp: 0.75, scaleDown: 0.25 },
    },
    {
        name: "autoscale-api",
        namespace: "production",
        targetDeployment: "api-deployment",
        metricQuery: "sum(rate(api_requests_total[1m]))",
        thresholds: { scaleUp: 0.9, scaleDown: 0.4 },
    },
];

export default function Dashboard() {
    const [policies, setPolicies] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setTimeout(() => {
            setPolicies(mockPolicies); // Replace with API later
            setLoading(false);
        }, 1000);
    }, []);

    return (
        <>
            <h1 className="text-2xl font-semibold mb-4">Scale Policies</h1>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {loading ? (
                    <LoadingSpinner />
                ) : (
                    policies.map((policy, i) => (
                        <PolicyCard key={i} policy={policy} />
                    ))
                )}
            </div>

            <div className="mt-12">
                <h2 className="text-xl font-semibold mb-2">Metric Overview</h2>
                <Chart />
            </div>
        </>
    );
}
