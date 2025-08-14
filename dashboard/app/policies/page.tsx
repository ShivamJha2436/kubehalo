"use client"

import { useEffect, useState } from "react";
import PolicyTable from "../../components/PolicyTable";
import { fetchPolicies } from "../../lib/api";
import Button from "../../components/ui/Button";
import { ScalePolicy } from "../../types";

export default function PoliciesPage() {
    const [policies, setPolicies] = useState<ScalePolicy[] | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchPolicies().then((p) => {
            setPolicies(p);
            setLoading(false);
        });
    }, []);

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-semibold">Scale Policies</h1>
                <div>
                    <Button href="/policies/new">Create policy</Button>
                </div>
            </div>

            <div>
                <div className="bg-white rounded shadow p-4">
                    <PolicyTable policies={policies} loading={loading} />
                </div>
            </div>
        </div>
    );
}
