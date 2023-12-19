import m from "mithril";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";

import state from "src/js/state";

export default (): m.Component => {
  const fitAddon = new FitAddon();

  let term: Terminal | null = null;
  let input = "";

  return {
    oncreate: ({ dom }) => {
      term = new Terminal();
      term.loadAddon(fitAddon);
      term.open(dom.querySelector(".term") as HTMLElement);
      term.write("\x1b[?25l"); // hide cursor

      state.pubsub.subscribe(
        `${state.topics.TOPIC_WS}.message`,
        (topic, msg) => {
          term?.write(msg.data as string);
        },
      );
    },
    view: () =>
      m("sl-card.w-100", [
        m("div[slot='header']", "Terminal into LUA Interpreter"), //
        m(
          "div.w-100[slot='image']",
          m(
            "sl-resize-observer.w-100",
            {
              "onsl-resize": () => {
                fitAddon?.fit();
              },
            },
            m("div.term"),
          ),
        ), //
        m("div.flex", [
          m("sl-input.flex-grow-1.mr2", {
            value: input,
            oninput: (e: InputEvent) => (input = e.target?.value ?? ""),
          }), //
          m(
            "sl-button",
            {
              onclick: () => {
                state.socket.send(input);
                input = "";
              },
            },
            "Execute Command",
          ),
        ]),
      ]),
  };
};
