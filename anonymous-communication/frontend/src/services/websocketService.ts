import type { OutgoingWebSocketMessage, WebSocketEvent } from "../types/websocket";
import { API_BASE_URL } from "../utils/constants";

function buildWebSocketURL() {
  const origin = typeof window === "undefined" ? "http://localhost:5173" : window.location.origin;
  const url = new URL(API_BASE_URL || origin, origin);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  url.pathname = "/ws/chat";
  url.search = "";
  return url.toString();
}

export class ChatSocket {
  private socket: WebSocket | null = null;

  connect(onMessage: (event: WebSocketEvent) => void, onClose?: () => void) {
    this.socket = new WebSocket(buildWebSocketURL());
    this.socket.onmessage = (event) => {
      onMessage(JSON.parse(event.data) as WebSocketEvent);
    };
    this.socket.onclose = () => {
      if (onClose) {
        onClose();
      }
    };
  }

  send(payload: OutgoingWebSocketMessage) {
    this.socket?.send(JSON.stringify(payload));
  }

  close() {
    this.socket?.close();
    this.socket = null;
  }
}
