import { useEffect, useRef, useState } from "react";
import { ChatSocket } from "../services/websocketService";
import type { OutgoingWebSocketMessage, WebSocketEvent } from "../types/websocket";

export function useWebSocket(enabled: boolean) {
  const socketRef = useRef<ChatSocket | null>(null);
  const [events, setEvents] = useState<WebSocketEvent[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    if (!enabled) {
      socketRef.current?.close();
      socketRef.current = null;
      setIsConnected(false);
      return;
    }

    const socket = new ChatSocket();
    socketRef.current = socket;
    socket.connect(
      (event) => {
        setEvents((current) => [...current, event]);
        setIsConnected(true);
      },
      () => {
        setIsConnected(false);
      },
    );

    return () => {
      socket.close();
      socketRef.current = null;
      setIsConnected(false);
    };
  }, [enabled]);

  function send(payload: OutgoingWebSocketMessage) {
    socketRef.current?.send(payload);
  }

  return {
    events,
    isConnected,
    send,
  };
}
