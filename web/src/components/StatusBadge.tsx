interface Props { status: string }

const styles: Record<string, string> = {
  pending: "bg-yellow-100 text-yellow-800",
  sent: "bg-green-100 text-green-800",
  failed: "bg-red-100 text-red-800",
  success: "bg-green-100 text-green-800",
};

export default function StatusBadge({ status }: Props) {
  const cls = styles[status?.toLowerCase()] ?? "bg-gray-100 text-gray-700";
  return (
    <span className={`inline-flex items-center rounded-full px-3 py-0.5 text-xs font-semibold capitalize ${cls}`}>
      {status || "—"}
    </span>
  );
}
