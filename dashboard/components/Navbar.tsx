"use client"

import React from "react";

export default function Navbar() {
    return (
        <header className="flex items-center justify-between px-4 py-3 bg-white border-b">
            <div className="flex items-center space-x-4">
                <button className="md:hidden p-2 rounded-md hover:bg-gray-100">☰</button>
                <div className="text-lg font-semibold">KubeHalo</div>
            </div>

            <div className="flex items-center space-x-3">
                <div className="text-sm text-gray-600">Shivam</div>
                <div className="w-8 h-8 bg-gray-200 rounded-full" />
            </div>
        </header>
    );
}
