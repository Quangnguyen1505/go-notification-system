"use client";
import { useEffect } from "react";

export type ToastType = "success" | "error" | "info";

interface Props {
  message: string;
  type: ToastType;
  onClose: () => void;
}

const colors: Record<ToastType, string> = {
  success: "bg-green-600",
  error: "bg-red-600",
  info: "bg-blue-600",
};

export default function Toast({ message, type, onClose }: Props) {
  useEffect(() => {
    const t = setTimeout(onClose, 4000);
    return () => clearTimeout(t);
  }, [onClose]);

  return (
    <div
      className={`fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-lg px-5 py-3 text-white shadow-lg ${colors[type]} animate-slide-in`}
    >
      <span className="text-sm font-medium">{message}</span>
      <button onClick={onClose} className="ml-2 text-white/80 hover:text-white">
        ✕
      </button>
    </div>
  );
}
