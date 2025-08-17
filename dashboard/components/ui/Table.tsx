import React from "react";

export default function Table({ children }: { children: React.ReactNode }) {
    return (
        <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">{children}</table>
        </div>
    );
}
