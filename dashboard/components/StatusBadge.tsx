type StatusBadgeProps = {
    type: "scaleUp" | "scaleDown";
    value: number;
};

export default function StatusBadge({ type, value }: StatusBadgeProps) {
    const color = type === "scaleUp" ? "green" : "red";
    return (
        <span className={`text-sm font-medium text-${color}-700 bg-${color}-100 px-2 py-1 rounded-md`}>
      {type === "scaleUp" ? "⬆ Scale Up:" : "⬇ Scale Down:"} {value}
    </span>
    );
}
