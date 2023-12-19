import m from "mithril";

import "@shoelace-style/shoelace/dist/shoelace.js";
import "@shoelace-style/shoelace/dist/themes/light.css";

import "tachyons/css/tachyons.css";
import "xterm/css/xterm.css";
import "leaflet/dist/leaflet.css";

import Home from "./pages/home";
import Registered from "./pages/registered";
import StateVis from "./pages/state-vis";

m.route(document.getElementById("app")!, "/", {
  "/": Home,
  "/registered": Registered,
  "/state-vis": StateVis,
});
