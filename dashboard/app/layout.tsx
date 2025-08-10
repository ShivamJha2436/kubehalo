import "@/app/globals.css";
import Sidebar from "@/components/Sidebar";
import Navbar from "@/components/Navbar";

export const metadata = {
    title: "KubeHalo Dashboard",
    description: "Kubernetes Autoscaling Policies Dashboard",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en">
        <body className="flex bg-gray-100">
        <Sidebar />
        <div className="flex flex-col flex-1 min-h-screen">
            <Navbar />
            <main className="p-6">{children}</main>
        </div>
        </body>
        </html>
    );
}
