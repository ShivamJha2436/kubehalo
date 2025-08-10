import { Line } from "react-chartjs-2";
import { Chart as ChartJS, LineElement, PointElement, CategoryScale, LinearScale } from "chart.js";

ChartJS.register(LineElement, PointElement, CategoryScale, LinearScale);

export default function Chart() {
    const data = {
        labels: ["00:00", "01:00", "02:00", "03:00", "04:00"],
        datasets: [
            {
                label: "http_requests_total",
                data: [50, 80, 65, 90, 120],
                fill: false,
                borderColor: "#3B82F6",
                tension: 0.2,
            },
        ],
    };

    return (
        <div className="bg-white rounded-lg p-4 shadow-md">
            <Line data={data} />
        </div>
    );
}
