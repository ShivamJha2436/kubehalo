export default function LoadingSpinner() {
    return (
        <div className="col-span-2 flex justify-center items-center py-12">
            <div className="w-10 h-10 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
    );
}
