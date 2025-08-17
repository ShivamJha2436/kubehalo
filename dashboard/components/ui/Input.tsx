"use client"

import React from "react";

type Props = React.InputHTMLAttributes<HTMLInputElement> & {
    label?: string;
    className?: string;
};

export default function Input({ label, className = "", ...rest }: Props) {
    return (
        <label className="block">
            {label && <div className="text-sm text-gray-600 mb-1">{label}</div>}
            <input
                {...rest}
                className={`w-full px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring ${className}`}
            />
        </label>
    );
}
