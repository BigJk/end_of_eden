import m from "mithril";

import Base from "src/js/components/base";

import "active-table/dist/activeTable.js";

export default (): m.Component => {
  const data = {
    artifacts: null,
    cards: null,
    statusEffects: null,
    events: null,
    storyTeller: null,
  } as {
    artifacts: any;
    cards: any;
    statusEffects: any;
    events: any;
    storyTeller: any;
  };

  const fetchData = (type: string) => {
    return m.request({
      method: "GET",
      url: "/api/registered/" + type,
    });
  };

  const toTable = (data: any) => {
    if (!data) {
      return null;
    }

    let values = Object.values(data);
    if (values.length === 0) {
      return [];
    }

    let keys = Object.keys(values[0]);
    values = values.map((v: any) => {
      return Object.values(v);
    });

    return [keys, ...values];
  };

  const fetchAll = () => {
    const types = ["artifacts", "cards", "status_effects", "events"];
    Promise.all(types.map((t) => fetchData(t))).then((d) => {
      d.forEach((d: any, i) => {
        data[types[i]] = d;
      });
      console.log("redraw");
      m.redraw();
    });
  };

  return {
    oninit: () => {
      fetchAll();
    },
    view: () => {
      if (data.artifacts === null) {
        return null;
      }

      return m(
        Base,
        {
          title: "Registered",
          subtitle: "Registered objects in the lua engine",
        },
        m("sl-tab-group", [
          m('sl-tab[slot="nav"][panel="artifacts"]', "Artifacts"),
          m('sl-tab[slot="nav"][panel="cards"]', "Cards"),
          m('sl-tab[slot="nav"][panel="status_effects"]', "Status Effects"),
          m('sl-tab[slot="nav"][panel="events"]', "Events"),
          m('sl-tab[slot="nav"][panel="story_teller"]', "Story Teller"),

          Object.keys(data).map((k) => {
            return m(
              `sl-tab-panel[name="${k}"]`,
              m("active-table", {
                tableStyle: {
                  width: "100%",
                },
                rowDropdown: {
                  isMoveAvailable: false,
                  isInsertUpAvailable: false,
                  isInsertDownAvailable: false,
                  isDeleteAvailable: false,
                  canEditHeaderRow: false,
                },
                pagination: true,
                isCellTextEditable: false,
                isHeaderTextEditable: false,
                displayAddNewRow: false,
                displayAddNewColumn: false,
                content: toTable(data[k]),
              }),
            );
          }),
        ]),
      );
    },
  };
};
