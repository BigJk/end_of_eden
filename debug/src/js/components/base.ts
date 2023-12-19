import m from "mithril";

import "@vaadin/app-layout";
import "@vaadin/app-layout/vaadin-drawer-toggle.js";
import "@vaadin/tabs";

const tabs = [
  {
    title: "Registered",
    link: "/registered",
    icon: "lumo:user",
  },
  {
    title: "State Visualization",
    link: "/state-vis",
    icon: "lumo:eye",
  },
];

const getTabIndex = (path: string): number => {
  const tab = tabs.find((t) => t.link === path);
  if (tab) {
    return tabs.indexOf(tab);
  }
  return 0;
};

export default (): m.Component<{ title: string; subtitle: string }> => {
  return {
    view: ({ children, attrs }) => {
      return m("vaadin-app-layout", [
        m("vaadin-drawer-toggle", { slot: "navbar touch-optimized" }),
        m("h3", { slot: "navbar touch-optimized" }, "End Of Eden Dev Tools"),
        m(
          "vaadin-tabs",
          {
            orientation: "vertical",
            slot: "drawer",
            selected: getTabIndex(m.route.get()),
          },
          tabs.map((t) => {
            return m("vaadin-tab", [
              m("a", { href: "#!" + t.link }, [
                m("vaadin-icon", { icon: t.icon }),
                t.title,
              ]),
            ]);
          }),
        ),
        m(
          "div.pa3",
          m("div", [
            m("div.f3.pb3", attrs.title, m("div.f6.pt2", attrs.subtitle)),
            children,
          ]),
        ),
      ]);
    },
  };
};
