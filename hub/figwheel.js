try {
    require("source-map-support").install();
} catch(err) {
}
require("./out.dev/goog/bootstrap/nodejs.js");
require("./out.dev/hub.js");
goog.require("hub.dev");
goog.require("cljs.nodejscli");
