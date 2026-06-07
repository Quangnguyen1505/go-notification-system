"use client";
import type { Channel } from "@/lib/api";

interface Props {
  value: Channel[];
  onChange: (channels: Channel[]) => void;
}

const CHANNELS: { id: Channel; label: string; icon: string; color: string }[] = [
  { id: "email", label: "Email", icon: "✉️", color: "border-blue-500 bg-blue-50 text-blue-700" },
  { id: "sms", label: "SMS", icon: "📱", color: "border-green-500 bg-green-50 text-green-700" },
  { id: "push", label: "Push", icon: "🔔", color: "border-purple-500 bg-purple-50 text-purple-700" },
];

export default function ChannelPicker({ value, onChange }: Props) {
  const toggle = (ch: Channel) => {
    onChange(value.includes(ch) ? value.filter((c) => c !== ch) : [...value, ch]);
  };

  return (
    <div className="flex gap-3 flex-wrap">
      {CHANNELS.map((ch) => {
        const active = value.includes(ch.id);
        return (
          <button
            key={ch.id}
            type="button"
            onClick={() => toggle(ch.id)}
            className={`flex items-center gap-2 rounded-xl border-2 px-5 py-3 font-medium transition-all select-none
              ${active ? ch.color + " border-opacity-100 shadow-sm scale-105" : "border-gray-200 bg-white text-gray-500 hover:border-gray-300"}`}
          >
            <span className="text-xl">{ch.icon}</span>
            {ch.label}
            {active && <span className="ml-1 text-xs">✓</span>}
          </button>
        );
      })}
    </div>
  );
}
