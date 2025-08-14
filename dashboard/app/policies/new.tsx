"use client"

import PolicyForm from "../../components/PolicyForm";
import { useRouter } from "next/navigation";
import { createPolicy } from "../../lib/api";

export default function NewPolicyPage() {
    const router = useRouter();

    async function handleSubmit(data: any) {
        await createPolicy(data);
        router.push("/policies");
    }

    return (
        <div>
            <h1 className="text-2xl font-semibold mb-4">Create Scale Policy</h1>
            <div className="bg-white p-6 rounded shadow">
                <PolicyForm onSubmit={handleSubmit} />
            </div>
        </div>
    );
}
