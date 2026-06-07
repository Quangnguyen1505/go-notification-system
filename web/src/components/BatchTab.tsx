"use client";
import { useState } from "react";
import ChannelPicker from "./ChannelPicker";
import StatusBadge from "./StatusBadge";
import type { Channel, BatchResult, CreateNotificationRequest } from "@/lib/api";
import { api } from "@/lib/api";
import { addHistory } from "@/lib/history";

interface Props {
  onToast: (msg: string, type: "success" | "error") => void;
}

interface Row {
  id: string;
  userId: string;
  channels: Channel[];
  subject: string;
  body: string;
}

function emptyRow(): Row {
  return { id: crypto.randomUUID(), userId: "", channels: [], subject: "", body: "" };
}

export default function BatchTab({ onToast }: Props) {
  const [rows, setRows] = useState<Row[]>([emptyRow()]);
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<BatchResult[] | null>(null);

  const update = (id: string, patch: Partial<Row>) =>
    setRows((r) => r.map((row) => (row.id === id ? { ...row, ...patch } : row)));

  const removeRow = (id: string) => setRows((r) => r.filter((row) => row.id !== id));

  const send = async () => {
    const invalid = rows.find((r) => !r.userId || !r.channels.length || !r.subject || !r.body);
    if (invalid) { onToast("Fill in all fields for every row", "error"); return; }

    setLoading(true);
    try {
      const notifications: CreateNotificationRequest[] = rows.map((r) => ({
        user_id: r.userId,
        channels: r.channels,
        subject: r.subject,
        body: r.body,
      }));
      const res = await api.batchNotification({ notifications });
      setResults(res.results);
      res.results.forEach((r, i) => {
        addHistory({
          id: crypto.randomUUID(),
          sentAt: new Date().toISOString(),
          userId: rows[i]?.userId ?? "",
          channels: rows[i]?.channels ?? [],
          subject: rows[i]?.subject ?? "",
          status: r.status,
          notificationId: r.notification_id,
        });
      });
      onToast(`Batch sent: ${res.results.length} notifications`, "success");
    } catch (err) {
      onToast((err as Error).message, "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-base font-semibold text-gray-900">Batch Notifications</h2>
          <button
            onClick={() => setRows((r) => [...r, emptyRow()])}
            className="flex items-center gap-1 rounded-lg border border-indigo-200 bg-indigo-50 px-3 py-1.5 text-sm font-medium text-indigo-700 hover:bg-indigo-100"
          >
            + Add Row
          </button>
        </div>

        <div className="space-y-4">
          {rows.map((row, i) => (
            <div key={row.id} className="rounded-xl border border-gray-200 p-4">
              <div className="mb-3 flex items-center justify-between">
                <span className="text-sm font-medium text-gray-500">#{i + 1}</span>
                {rows.length > 1 && (
                  <button
                    onClick={() => removeRow(row.id)}
                    className="text-xs text-red-400 hover:text-red-600"
                  >
                    Remove
                  </button>
                )}
              </div>
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <input
                  placeholder="User ID"
                  value={row.userId}
                  onChange={(e) => update(row.id, { userId: e.target.value })}
                  className="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
                <input
                  placeholder="Subject"
                  value={row.subject}
                  onChange={(e) => update(row.id, { subject: e.target.value })}
                  className="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div className="mt-3">
                <ChannelPicker value={row.channels} onChange={(ch) => update(row.id, { channels: ch })} />
              </div>
              <textarea
                placeholder="Message body"
                value={row.body}
                rows={2}
                onChange={(e) => update(row.id, { body: e.target.value })}
                className="mt-3 w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 resize-none"
              />
            </div>
          ))}
        </div>

        <div className="mt-5 flex justify-end gap-3">
          <button
            onClick={() => { setRows([emptyRow()]); setResults(null); }}
            className="rounded-lg border border-gray-200 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50"
          >
            Clear All
          </button>
          <button
            onClick={send}
            disabled={loading}
            className="flex items-center gap-2 rounded-lg bg-indigo-600 px-5 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
          >
            {loading ? (
              <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
            ) : "→"}
            Send Batch
          </button>
        </div>
      </div>

      {results && (
        <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
          <h2 className="mb-4 text-base font-semibold text-gray-900">Batch Results</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-100 text-left text-xs text-gray-500">
                  <th className="pb-2 pr-4">#</th>
                  <th className="pb-2 pr-4">Notification ID</th>
                  <th className="pb-2 pr-4">Status</th>
                  <th className="pb-2">Error</th>
                </tr>
              </thead>
              <tbody>
                {results.map((r, i) => (
                  <tr key={i} className="border-b border-gray-50">
                    <td className="py-2 pr-4 text-gray-500">{i + 1}</td>
                    <td className="py-2 pr-4 font-mono text-xs text-gray-700 break-all">{r.notification_id || "—"}</td>
                    <td className="py-2 pr-4"><StatusBadge status={r.status} /></td>
                    <td className="py-2 text-red-500 text-xs">{r.error || "—"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
