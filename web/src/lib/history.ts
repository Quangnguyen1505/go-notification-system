export interface HistoryEntry {
  id: string;
  sentAt: string;
  userId: string;
  channels: string[];
  subject: string;
  status: string;
  notificationId: string;
}

const KEY = "noti_history";

export function getHistory(): HistoryEntry[] {
  if (typeof window === "undefined") return [];
  try {
    return JSON.parse(localStorage.getItem(KEY) || "[]");
  } catch {
    return [];
  }
}

export function addHistory(entry: HistoryEntry) {
  const h = getHistory();
  localStorage.setItem(KEY, JSON.stringify([entry, ...h].slice(0, 100)));
}

export function updateHistoryStatus(notificationId: string, status: string) {
  const h = getHistory();
  const updated = h.map((e) =>
    e.notificationId === notificationId ? { ...e, status } : e
  );
  localStorage.setItem(KEY, JSON.stringify(updated));
  return updated;
}

export function clearHistory() {
  localStorage.removeItem(KEY);
}
