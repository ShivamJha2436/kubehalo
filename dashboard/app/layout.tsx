import "@/app/globals.css";
import Navbar from "../components/Navbar";
import Sidebar from "../components/Sidebar";

export const metadata = {
    title: "KubeHalo Dashboard",
    description: "KubeHalo — Kubernetes autoscaling dashboard",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en">
        <body className="min-h-screen bg-gray-50 text-gray-900">
        <div className="flex h-screen">
            <Sidebar />
            <div className="flex-1 flex flex-col">
                <Navbar />
                <main className="p-6 overflow-auto">{children}</main>
            </div>
        </div>
        </body>
        </html>
    );
}
