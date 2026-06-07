"use client";
import { useState, useEffect, useMemo } from "react";
import StatusBadge from "./StatusBadge";
import { getHistory, clearHistory, updateHistoryStatus, type HistoryEntry } from "@/lib/history";
import { api } from "@/lib/api";

const PAGE_SIZE = 10;

const CHANNEL_ICONS: Record<string, string> = { email: "✉️", sms: "📱", push: "🔔" };
const ALL_CHANNELS = ["email", "sms", "push"] as const;
const ALL_STATUSES = ["created", "pending", "sent", "failed"];

export default function HistoryTab() {
  const [entries, setEntries] = useState<HistoryEntry[]>([]);
  const [search, setSearch] = useState("");
  const [filterChannel, setFilterChannel] = useState<string>("all");
  const [filterStatus, setFilterStatus] = useState<string>("all");
  const [sortOrder, setSortOrder] = useState<"desc" | "asc">("desc");
  const [page, setPage] = useState(1);
  const [refreshingId, setRefreshingId] = useState<string | null>(null);

  useEffect(() => { setEntries(getHistory()); }, []);

  const filtered = useMemo(() => {
    let list = [...entries];

    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter(
        (e) =>
          e.userId.toLowerCase().includes(q) ||
          e.subject.toLowerCase().includes(q) ||
          e.notificationId.toLowerCase().includes(q)
      );
    }

    if (filterChannel !== "all") {
      list = list.filter((e) => e.channels.includes(filterChannel));
    }

    if (filterStatus !== "all") {
      list = list.filter((e) => e.status === filterStatus);
    }

    list.sort((a, b) => {
      const diff = new Date(a.sentAt).getTime() - new Date(b.sentAt).getTime();
      return sortOrder === "desc" ? -diff : diff;
    });

    return list;
  }, [entries, search, filterChannel, filterStatus, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE));
  const currentPage = Math.min(page, totalPages);
  const paged = filtered.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE);

  const clear = () => { clearHistory(); setEntries([]); };

  const refreshStatus = async (entry: HistoryEntry) => {
    if (!entry.notificationId) return;
    setRefreshingId(entry.notificationId);
    try {
      const res = await api.getNotification(entry.notificationId);
      const updated = updateHistoryStatus(entry.notificationId, res.status);
      setEntries(updated);
    } catch {
      // silently ignore if lookup fails
    } finally {
      setRefreshingId(null);
    }
  };

  const resetFilters = () => {
    setSearch(""); setFilterChannel("all"); setFilterStatus("all"); setPage(1);
  };

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="rounded-2xl bg-white p-5 shadow-sm border border-gray-100">
        <div className="mb-4 flex items-center justify-between">
          <div>
            <h2 className="text-base font-semibold text-gray-900">All Notifications</h2>
            <p className="text-xs text-gray-400 mt-0.5">
              {filtered.length} of {entries.length} entries
            </p>
          </div>
          {entries.length > 0 && (
            <button onClick={clear} className="text-sm text-red-400 hover:text-red-600">
              Clear All
            </button>
          )}
        </div>

        <div className="flex flex-wrap gap-3">
          {/* Search */}
          <div className="relative flex-1 min-w-[180px]">
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm">🔍</span>
            <input
              type="text"
              placeholder="Search user ID, subject, notification ID…"
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              className="w-full rounded-lg border border-gray-200 pl-8 pr-3 py-2 text-sm focus:border-indigo-400 focus:outline-none focus:ring-1 focus:ring-indigo-400"
            />
          </div>

          {/* Channel filter */}
          <select
            value={filterChannel}
            onChange={(e) => { setFilterChannel(e.target.value); setPage(1); }}
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-700 focus:border-indigo-400 focus:outline-none"
          >
            <option value="all">All Channels</option>
            {ALL_CHANNELS.map((c) => (
              <option key={c} value={c}>{CHANNEL_ICONS[c]} {c.charAt(0).toUpperCase() + c.slice(1)}</option>
            ))}
          </select>

          {/* Status filter */}
          <select
            value={filterStatus}
            onChange={(e) => { setFilterStatus(e.target.value); setPage(1); }}
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-700 focus:border-indigo-400 focus:outline-none"
          >
            <option value="all">All Statuses</option>
            {ALL_STATUSES.map((s) => (
              <option key={s} value={s}>{s.charAt(0).toUpperCase() + s.slice(1)}</option>
            ))}
          </select>

          {/* Sort */}
          <button
            onClick={() => setSortOrder((o) => (o === "desc" ? "asc" : "desc"))}
            className="flex items-center gap-1.5 rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 hover:bg-gray-50"
          >
            {sortOrder === "desc" ? "↓ Newest" : "↑ Oldest"}
          </button>

          {/* Reset */}
          {(search || filterChannel !== "all" || filterStatus !== "all") && (
            <button
              onClick={resetFilters}
              className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-500 hover:bg-gray-50"
            >
              ✕ Reset
            </button>
          )}
        </div>
      </div>

      {/* Table */}
      <div className="rounded-2xl bg-white shadow-sm border border-gray-100 overflow-hidden">
        {entries.length === 0 ? (
          <div className="flex flex-col items-center py-20 text-gray-400">
            <span className="text-5xl mb-3">🔔</span>
            <p className="text-sm font-medium">No notifications yet</p>
            <p className="text-xs mt-1">Send a notification and it will appear here.</p>
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center py-16 text-gray-400">
            <span className="text-4xl mb-3">🔍</span>
            <p className="text-sm">No results match your filters.</p>
            <button onClick={resetFilters} className="mt-2 text-xs text-indigo-500 hover:underline">
              Clear filters
            </button>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-100 bg-gray-50/60 text-left text-xs text-gray-500">
                    <th className="px-4 py-3 whitespace-nowrap font-medium">Time</th>
                    <th className="px-4 py-3 font-medium">User ID</th>
                    <th className="px-4 py-3 font-medium">Channels</th>
                    <th className="px-4 py-3 font-medium">Subject</th>
                    <th className="px-4 py-3 font-medium">Status</th>
                    <th className="px-4 py-3 font-medium">Notification ID</th>
                    <th className="px-4 py-3 font-medium"></th>
                  </tr>
                </thead>
                <tbody>
                  {paged.map((e) => (
                    <tr key={e.id} className="border-b border-gray-50 hover:bg-indigo-50/30 transition-colors">
                      <td className="px-4 py-3 whitespace-nowrap text-xs text-gray-500">
                        {new Date(e.sentAt).toLocaleString()}
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-gray-600 max-w-[130px] truncate" title={e.userId}>
                        {e.userId}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex gap-1">
                          {e.channels.map((c) => (
                            <span key={c} title={c} className="text-base leading-none">
                              {CHANNEL_ICONS[c] ?? c}
                            </span>
                          ))}
                        </div>
                      </td>
                      <td className="px-4 py-3 max-w-[180px] truncate text-gray-700" title={e.subject}>
                        {e.subject}
                      </td>
                      <td className="px-4 py-3">
                        <StatusBadge status={e.status} />
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-gray-400 max-w-[150px] truncate" title={e.notificationId}>
                        {e.notificationId || "—"}
                      </td>
                      <td className="px-4 py-3">
                        {e.notificationId && (
                          <button
                            onClick={() => refreshStatus(e)}
                            disabled={refreshingId === e.notificationId}
                            title="Refresh status from server"
                            className="text-xs text-indigo-400 hover:text-indigo-600 disabled:opacity-40"
                          >
                            {refreshingId === e.notificationId ? (
                              <span className="inline-block h-3.5 w-3.5 animate-spin rounded-full border-2 border-indigo-400 border-t-transparent" />
                            ) : "↻"}
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-between border-t border-gray-100 px-4 py-3">
                <span className="text-xs text-gray-500">
                  Page {currentPage} of {totalPages}
                </span>
                <div className="flex gap-2">
                  <button
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="rounded-lg border border-gray-200 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50 disabled:opacity-40"
                  >
                    ← Prev
                  </button>
                  <button
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                    className="rounded-lg border border-gray-200 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50 disabled:opacity-40"
                  >
                    Next →
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
