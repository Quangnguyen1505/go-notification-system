const DEFAULT_BASE = "http://localhost:5000";

function getBase(): string {
  if (typeof window !== "undefined") {
    return localStorage.getItem("apiBase") || DEFAULT_BASE;
  }
  return DEFAULT_BASE;
}

export type Channel = "email" | "sms" | "push";

export interface CreateNotificationRequest {
  user_id: string;
  channels: Channel[];
  subject: string;
  body: string;
  schedule_at?: number;
}

export interface CreateNotificationResponse {
  notification_id: string;
  status: string;
}

export interface BatchNotificationRequest {
  notifications: CreateNotificationRequest[];
}

export interface BatchResult {
  notification_id: string;
  status: string;
  error?: string;
}

export interface BatchNotificationResponse {
  results: BatchResult[];
}

export interface GetNotificationResponse {
  notification_id: string;
  status: string;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${getBase()}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || data.error || `HTTP ${res.status}`);
  return data as T;
}

export const api = {
  createNotification(body: CreateNotificationRequest) {
    return request<CreateNotificationResponse>("/v1/api/notification", {
      method: "POST",
      body: JSON.stringify(body),
    });
  },
  batchNotification(body: BatchNotificationRequest) {
    return request<BatchNotificationResponse>("/v1/api/notification:batch", {
      method: "POST",
      body: JSON.stringify(body),
    });
  },
  getNotification(id: string) {
    return request<GetNotificationResponse>(`/v1/api/notification/${id}`);
  },
};
