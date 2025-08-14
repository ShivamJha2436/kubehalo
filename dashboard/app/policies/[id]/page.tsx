"use client"

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import PolicyForm from "../../../components/PolicyForm";
import { getPolicy, updatePolicy } from "../../../lib/api";
import { ScalePolicy } from "../../../types";

export default function PolicyDetail() {
    const params = useParams();
    const id = params?.id;
    const router = useRouter();
    const [policy, setPolicy] = useState<ScalePolicy | null>(null);

    useEffect(() => {
        if (!id) return;
        getPolicy(id).then((p) => setPolicy(p));
    }, [id]);

    async function handleSubmit(data: any) {
        if (!id) return;
        await updatePolicy(id, data);
        router.push("/policies");
    }

    if (!policy) return <div>Loading...</div>;

    return (
        <div>
            <h1 className="text-2xl font-semibold mb-4">Edit Policy — {policy.metadata.name}</h1>
            <div className="bg-white p-6 rounded shadow">
                <PolicyForm initialValues={policy} onSubmit={handleSubmit} />
            </div>
        </div>
    );
}
