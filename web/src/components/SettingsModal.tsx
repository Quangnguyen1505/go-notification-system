"use client";
import { useState, useEffect } from "react";

interface Props {
  open: boolean;
  onClose: () => void;
}

export default function SettingsModal({ open, onClose }: Props) {
  const [base, setBase] = useState("http://localhost:5000");

  useEffect(() => {
    setBase(localStorage.getItem("apiBase") || "http://localhost:5000");
  }, [open]);

  const save = () => {
    localStorage.setItem("apiBase", base.replace(/\/$/, ""));
    onClose();
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">API Settings</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xl">✕</button>
        </div>
        <label className="mb-1 block text-sm font-medium text-gray-700">
          Backend Base URL
        </label>
        <input
          type="text"
          value={base}
          onChange={(e) => setBase(e.target.value)}
          className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          placeholder="http://localhost:5000"
        />
        <p className="mt-1 text-xs text-gray-500">Proxy service that forwards to the gRPC notification service.</p>
        <div className="mt-5 flex justify-end gap-2">
          <button onClick={onClose} className="rounded-lg border border-gray-200 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50">
            Cancel
          </button>
          <button onClick={save} className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700">
            Save
          </button>
        </div>
      </div>
    </div>
  );
}
