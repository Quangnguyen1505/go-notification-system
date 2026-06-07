"use client";
import { useState } from "react";
import StatusBadge from "./StatusBadge";
import type { GetNotificationResponse } from "@/lib/api";
import { api } from "@/lib/api";

interface Props {
  onToast: (msg: string, type: "success" | "error") => void;
}

export default function LookupTab({ onToast }: Props) {
  const [id, setId] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<GetNotificationResponse | null>(null);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await api.getNotification(id.trim());
      setResult(res);
    } catch (err) {
      onToast((err as Error).message, "error");
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
        <h2 className="mb-5 text-base font-semibold text-gray-900">Check Notification Status</h2>
        <form onSubmit={submit} className="flex gap-3">
          <input
            required
            value={id}
            onChange={(e) => setId(e.target.value)}
            placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
            className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm font-mono focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          />
          <button
            type="submit"
            disabled={loading}
            className="flex items-center gap-2 rounded-lg bg-indigo-600 px-5 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
          >
            {loading ? (
              <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
            ) : "🔍"} Lookup
          </button>
        </form>
      </div>

      {result && (
        <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-semibold text-gray-900">Notification Details</h2>
            <StatusBadge status={result.status} />
          </div>
          <dl className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div className="rounded-lg bg-gray-50 px-4 py-3">
              <dt className="text-xs font-medium text-gray-500">Notification ID</dt>
              <dd className="mt-1 font-mono text-sm text-gray-900 break-all">{result.notification_id}</dd>
            </div>
            <div className="rounded-lg bg-gray-50 px-4 py-3">
              <dt className="text-xs font-medium text-gray-500">Status</dt>
              <dd className="mt-1 text-sm text-gray-900 capitalize">{result.status}</dd>
            </div>
          </dl>
        </div>
      )}
    </div>
  );
}
