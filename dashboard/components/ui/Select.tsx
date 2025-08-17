"use client"

import React from "react";

type Props = React.SelectHTMLAttributes<HTMLSelectElement> & {
    label?: string;
};

export default function Select({ label, children, ...rest }: Props) {
    return (
        <label className="block">
            {label && <div className="text-sm text-gray-600 mb-1">{label}</div>}
            <select {...rest} className="w-full px-3 py-2 border rounded-md text-sm">
                {children}
            </select>
        </label>
    );
}
