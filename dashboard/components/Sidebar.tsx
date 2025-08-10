import Image from "next/image";

export default function Sidebar() {
    return (
        <aside className="w-64 hidden md:block bg-white h-screen shadow-md p-4">
            <div className="flex items-center gap-2">
                <Image src="/kubehalo.png" alt="KubeHalo" width={96} height={26} />
                <h2 className="font-bold text-lg">KubeHalo</h2>
            </div>
            <nav className="mt-6">
                <ul>
                    <li className="text-sm text-gray-700 font-medium py-2">Dashboard</li>
                    <li className="text-sm text-gray-700 font-medium py-2">Policies</li>
                </ul>
            </nav>
        </aside>
    );
}
