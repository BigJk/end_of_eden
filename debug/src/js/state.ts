import m from "mithril";
import rwebsocket from "reconnecting-websocket";
import PubSub from "pubsub-js";

export const TOPIC_WS = "ws";

const socket = new rwebsocket(
  "ws://127.0.0.1:" + document.location.port + "/ws",
  undefined,
  {
    minReconnectionDelay: 1000,
    maxReconnectionDelay: 1000,
  },
);

socket.onopen = (event) => {
  PubSub.publish(`${TOPIC_WS}.open`, event);
  m.redraw();
};

socket.onclose = (event) => {
  PubSub.publish(`${TOPIC_WS}.close`, event);
  m.redraw();
};

socket.onerror = (event) => {
  PubSub.publish(`${TOPIC_WS}.error`, event);
  m.redraw();
};

socket.onmessage = (event) => {
  PubSub.publish(`${TOPIC_WS}.message`, {
    event: event,
    data: event.data + "\r\n",
  });
};

const state = {
  socket,
  pubsub: PubSub,
  topics: {
    TOPIC_WS,
  },
};

export default state;
