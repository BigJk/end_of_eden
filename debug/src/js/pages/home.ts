import m from "mithril";

import Base from "src/js/components/base";

export default (): m.Component => {
  return {
    view: () => m(Base, "Hello World!"),
  };
};
