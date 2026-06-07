"use client";
import { useState, useCallback } from "react";
import SendTab from "@/components/SendTab";
import BatchTab from "@/components/BatchTab";
import LookupTab from "@/components/LookupTab";
import HistoryTab from "@/components/HistoryTab";
import SettingsModal from "@/components/SettingsModal";
import Toast, { type ToastType } from "@/components/Toast";

type Tab = "send" | "batch" | "lookup" | "history";

const TABS: { id: Tab; label: string; icon: string }[] = [
  { id: "send", label: "Send", icon: "✉️" },
  { id: "batch", label: "Batch", icon: "📦" },
  { id: "lookup", label: "Lookup", icon: "🔍" },
  { id: "history", label: "All Notifications", icon: "🔔" },
];

const PAGE_META: Record<Tab, { title: string; subtitle: string }> = {
  send: { title: "Send Notification", subtitle: "Deliver messages via Email, SMS, or Push" },
  batch: { title: "Batch Send", subtitle: "Send multiple notifications in one request" },
  lookup: { title: "Lookup Status", subtitle: "Check the delivery status of a notification" },
  history: { title: "All Notifications", subtitle: "Browse, search and filter all sent notifications" },
};

interface ToastState { message: string; type: ToastType; key: number }

export default function Home() {
  const [tab, setTab] = useState<Tab>("send");
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [toast, setToast] = useState<ToastState | null>(null);

  const showToast = useCallback((message: string, type: ToastType) => {
    setToast({ message, type, key: Date.now() });
  }, []);

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="flex w-56 flex-col bg-white border-r border-gray-100 shadow-sm">
        <div className="flex items-center gap-2.5 px-5 py-5 border-b border-gray-100">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white text-sm">
            ✉
          </div>
          <span className="font-semibold text-gray-900">NotifyHub</span>
        </div>

        <nav className="flex-1 p-3 space-y-1">
          {TABS.map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`flex w-full items-center gap-2.5 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors
                ${tab === t.id
                  ? "bg-indigo-50 text-indigo-700"
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900"
                }`}
            >
              <span>{t.icon}</span>
              {t.label}
            </button>
          ))}
        </nav>

        <div className="border-t border-gray-100 p-4">
          <button
            onClick={() => setSettingsOpen(true)}
            className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-gray-500 hover:bg-gray-50 hover:text-gray-700"
          >
            ⚙️ Settings
          </button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <header className="border-b border-gray-100 bg-white px-8 py-5 shadow-sm">
          <h1 className="text-xl font-semibold text-gray-900">{PAGE_META[tab].title}</h1>
          <p className="mt-0.5 text-sm text-gray-500">{PAGE_META[tab].subtitle}</p>
        </header>

        <main className="flex-1 overflow-y-auto p-8">
          <div className="mx-auto max-w-3xl">
            {tab === "send" && <SendTab onToast={showToast} />}
            {tab === "batch" && <BatchTab onToast={showToast} />}
            {tab === "lookup" && <LookupTab onToast={showToast} />}
            {tab === "history" && <HistoryTab />}
          </div>
        </main>
      </div>

      <SettingsModal open={settingsOpen} onClose={() => setSettingsOpen(false)} />
      {toast && (
        <Toast
          key={toast.key}
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}
