export type APIStatus = "success" | "error";

export interface APIResponse<T> {
  status: APIStatus;
  message?: string;
  data?: T;
  error?: string;
  details?: Record<string, string>;
}

export class APIError extends Error {
  status: number;
  details?: Record<string, string>;

  constructor(message: string, status: number, details?: Record<string, string>) {
    super(message);
    this.name = "APIError";
    this.status = status;
    this.details = details;
  }
}
