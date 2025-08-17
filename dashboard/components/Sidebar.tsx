"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";

const items = [
    { href: "/", label: "Overview" },
    { href: "/policies", label: "Scale Policies" },
    { href: "/metrics", label: "Metrics" },
    { href: "/settings", label: "Settings" },
];

export default function Sidebar() {
    const pathname = usePathname();

    return (
        <aside className="w-64 bg-white border-r hidden md:block">
            <div className="h-full flex flex-col">
                <div className="p-4 border-b">
                    <img src="/kubehalo.png" alt="KubeHalo" className="h-6 w-6" />
                </div>
                <nav className="p-4 space-y-1">
                    {items.map((it) => {
                        const active = pathname === it.href;
                        return (
                            <Link
                                key={it.href}
                                href={it.href}
                                className={`block px-3 py-2 rounded-md text-sm ${active ? "bg-indigo-50 text-indigo-700" : "text-gray-700 hover:bg-gray-50"}`}
                            >
                                {it.label}
                            </Link>
                        );
                    })}
                </nav>
                <div className="mt-auto p-4 text-xs text-gray-500">KubeHalo &middot; Dashboard</div>
            </div>
        </aside>
    );
}
