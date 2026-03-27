import type { Message } from "./message";

export interface WebSocketEvent {
  type: string;
  payload?: Message;
  message?: string;
}

export interface OutgoingWebSocketMessage {
  type?: string;
  receiver_id: string;
  content: string;
}
