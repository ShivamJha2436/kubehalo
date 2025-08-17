import React from "react";

export default function Card({ title, value }: { title: string; value: string }) {
    return (
        <div className="bg-white rounded-md p-4 shadow">
            <div className="text-sm text-gray-500">{title}</div>
            <div className="mt-2 text-2xl font-semibold">{value}</div>
        </div>
    );
}
