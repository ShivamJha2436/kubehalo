"use client"

import Table from "./ui/Table";
import Spinner from "./ui/Spinner";
import Link from "next/link";
import { ScalePolicy } from "../types";

export default function PolicyTable({
                                        policies,
                                        loading,
                                    }: {
    policies: ScalePolicy[] | null;
    loading: boolean;
}) {
    if (loading) return <Spinner />;

    if (!policies || policies.length === 0)
        return <div className="text-sm text-gray-500">No policies found. Create one to get started.</div>;

    return (
        <Table>
            <thead className="bg-gray-50">
            <tr>
                <th className="px-4 py-2 text-left text-sm">Name</th>
                <th className="px-4 py-2 text-left text-sm">Target</th>
                <th className="px-4 py-2 text-left text-sm">Query</th>
                <th className="px-4 py-2 text-left text-sm">Threshold</th>
                <th className="px-4 py-2 text-left text-sm">Replicas</th>
                <th className="px-4 py-2 text-left text-sm">Actions</th>
            </tr>
            </thead>
            <tbody className="divide-y">
            {policies.map((p) => (
                <tr key={p.metadata.name}>
                    <td className="px-4 py-3 text-sm">{p.metadata.name}</td>
                    <td className="px-4 py-3 text-sm">{p.spec.targetRef.name}</td>
                    <td className="px-4 py-3 text-sm">{p.spec.metric.query}</td>
                    <td className="px-4 py-3 text-sm">{String(p.spec.metric.threshold)}</td>
                    <td className="px-4 py-3 text-sm">
                        {p.spec.minReplicas} - {p.spec.maxReplicas}
                    </td>
                    <td className="px-4 py-3 text-sm">
                        <Link href={`/policies/${p.metadata.name}`} className="text-indigo-600 hover:underline">
                            View
                        </Link>
                    </td>
                </tr>
            ))}
            </tbody>
        </Table>
    );
}
