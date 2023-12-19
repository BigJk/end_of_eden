import m from "mithril";

import Base from "src/js/components/base";
import StateVis from "src/js/components/state-vis";

export default (): m.Component => {
  return {
    view: () =>
      m(
        Base,
        {
          title: "State Visualization",
          subtitle:
            "This is a visualization of the current entities in the game state",
        },
        m(StateVis),
      ),
  };
};
