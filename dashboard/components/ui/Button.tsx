"use client"

import Link from "next/link";
import React from "react";

type Props = {
    children: React.ReactNode;
    href?: string;
    onClick?: () => void;
    className?: string;
};

export default function Button({ children, href, onClick, className = "" }: Props) {
    const base = "inline-flex items-center px-4 py-2 rounded-md text-sm font-medium";
    const style = "bg-indigo-600 text-white hover:bg-indigo-700";
    const classes = `${base} ${style} ${className}`;

    if (href) {
        return (
            <Link href={href} className={classes}>
                {children}
            </Link>
        );
    }
    return (
        <button onClick={onClick} className={classes}>
            {children}
        </button>
    );
}
