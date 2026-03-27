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
    const wsURL = buildWebSocketURL();
    
    // Get JWT token from localStorage
    const token = localStorage.getItem('auth_token');
    
    // Add token as query parameter if available
    const urlWithToken = token ? `${wsURL}?token=${encodeURIComponent(token)}` : wsURL;
    
    this.socket = new WebSocket(urlWithToken);
    
    this.socket.onopen = () => {
      // Send authentication message with token
      if (token) {
        this.socket?.send(JSON.stringify({
          type: 'authenticate',
          token: token
        }));
      }
    };
    
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
