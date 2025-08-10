export default function PolicyForm() {
    return (
        <form className="bg-white p-6 rounded shadow-md w-full max-w-lg">
            <h2 className="text-lg font-semibold mb-4">Add New Policy</h2>
            <div className="mb-3">
                <label className="block text-sm font-medium">Name</label>
                <input type="text" className="w-full mt-1 p-2 border rounded" />
            </div>
            <div className="mb-3">
                <label className="block text-sm font-medium">Namespace</label>
                <input type="text" className="w-full mt-1 p-2 border rounded" />
            </div>
            <div className="mb-3">
                <label className="block text-sm font-medium">Target Deployment</label>
                <input type="text" className="w-full mt-1 p-2 border rounded" />
            </div>
            <div className="mb-3">
                <label className="block text-sm font-medium">Metric Query</label>
                <input type="text" className="w-full mt-1 p-2 border rounded" />
            </div>
            <div className="flex gap-4 mb-4">
                <div className="flex-1">
                    <label className="block text-sm font-medium">Scale Up</label>
                    <input type="number" className="w-full mt-1 p-2 border rounded" />
                </div>
                <div className="flex-1">
                    <label className="block text-sm font-medium">Scale Down</label>
                    <input type="number" className="w-full mt-1 p-2 border rounded" />
                </div>
            </div>
            <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                Create Policy
            </button>
        </form>
    );
}
