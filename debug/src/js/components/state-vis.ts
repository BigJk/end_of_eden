import m from "mithril";
import L from "leaflet";

export default (): m.Component => {
  return {
    oncreate: ({ dom }) => {
      let map = L.map(dom as HTMLElement, {
        center: [40.75, -74.2],
        zoom: 13,
      });

      let imageUrl = "/api/svg",
        imageBounds = [
          [40.712216, -74.22655],
          [40.773941, -74.12544],
        ];

      L.imageOverlay(imageUrl, imageBounds).addTo(map);
    },
    view: () => m("div", { style: { height: "800px" } }),
  };
};
