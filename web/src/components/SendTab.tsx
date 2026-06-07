"use client";
import { useState } from "react";
import ChannelPicker from "./ChannelPicker";
import StatusBadge from "./StatusBadge";
import type { Channel, CreateNotificationResponse } from "@/lib/api";
import { api } from "@/lib/api";
import { addHistory } from "@/lib/history";

interface Props {
  onToast: (msg: string, type: "success" | "error") => void;
}

export default function SendTab({ onToast }: Props) {
  const [userId, setUserId] = useState("");
  const [channels, setChannels] = useState<Channel[]>([]);
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const [scheduleAt, setScheduleAt] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<CreateNotificationResponse | null>(null);

  const clear = () => {
    setUserId(""); setChannels([]); setSubject(""); setBody(""); setScheduleAt(""); setResult(null);
  };

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!channels.length) { onToast("Select at least one channel", "error"); return; }
    setLoading(true);
    try {
      const payload = {
        user_id: userId,
        channels,
        subject,
        body,
        ...(scheduleAt ? { schedule_at: Math.floor(new Date(scheduleAt).getTime() / 1000) } : {}),
      };
      const res = await api.createNotification(payload);
      setResult(res);
      addHistory({
        id: crypto.randomUUID(),
        sentAt: new Date().toISOString(),
        userId,
        channels,
        subject,
        status: res.status,
        notificationId: res.notification_id,
      });
      onToast("Notification sent!", "success");
    } catch (err) {
      onToast((err as Error).message, "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
        <h2 className="mb-5 text-base font-semibold text-gray-900">New Notification</h2>
        <form onSubmit={submit} className="space-y-5">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">
                User ID <span className="text-red-500">*</span>
              </label>
              <input
                required
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder="00000000-0000-0000-0000-000000000001"
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">Schedule At (optional)</label>
              <input
                type="datetime-local"
                value={scheduleAt}
                onChange={(e) => setScheduleAt(e.target.value)}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-gray-700">
              Channels <span className="text-red-500">*</span>
            </label>
            <ChannelPicker value={channels} onChange={setChannels} />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">
              Subject <span className="text-red-500">*</span>
            </label>
            <input
              required
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="Notification subject"
              className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">
              Message Body <span className="text-red-500">*</span>
            </label>
            <textarea
              required
              rows={5}
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="Enter your notification message…"
              className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 resize-none"
            />
            <p className="mt-1 text-right text-xs text-gray-400">{body.length} characters</p>
          </div>

          <div className="flex justify-end gap-3">
            <button type="button" onClick={clear} className="rounded-lg border border-gray-200 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50">
              Clear
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex items-center gap-2 rounded-lg bg-indigo-600 px-5 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
            >
              {loading ? (
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
              ) : "→"}
              Send Notification
            </button>
          </div>
        </form>
      </div>

      {result && (
        <div className="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-semibold text-gray-900">Result</h2>
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
