# JET/Slug

My explorations into the world of ClojureScript, Figwheel, and Reagent.  
Each "slug" is a self-contained project where I _slowly_ try something new:

* [01-slug](01-slug) - First try with Leiningen + Figwheel + Reagent + DevCards.
  Based on [reagent-figwheel](https://github.com/gadfly361/reagent-figwheel).
* [02-purecss](02-purecss) - Minimal Reagent setup with [Pure
  CSS](https://purecss.io) as layout framework.
* [03-paho](03-paho) - MQTT client in pure JavaScript, example from the [Eclipse
  Paho](https://www.eclipse.org/paho/clients/js/) website.
* [04-rum](04-rum) - Simple [Rum](https://github.com/tonsky/rum) example, which
  is a wrapper around React similar to Reagent.  Project created with
  [Chestnut](https://github.com/plexus/chestnut).
* [05-cljsjs](05-cljsjs) - Use Paho MQTT client in ClojureScript, using [this
  cljsjs example](https://github.com/cljsjs/packages/tree/master/paho). Created
  with `lein new figwheel app`.
* [06-boot](06-boot) - Set up a reagent app via [Boot](http://boot-clj.com),
  using `boot -d boot/new new -t tenzing -n your-app -a +reagent`.
* [07-devtools](07-devtools) - Demo of the
  [cljs-devtools](https://github.com/binaryage/cljs-devtools/blob/master/docs/installation.md)
  package, to show much better debugging info in the JavaScript console.
* [08-myapp](08-myapp) - A fairly comfortable dev setup with Reagent and
  Devtools.  Uses Boot, supports a browser REPL.
* [09-mini](09-mini) - A Rum-based project in 5 files. Tweaked from `boot -d
  boot/new new -t tenzing -n mini -a +rum`.
* [10-testme](10-testme) - Testme app from
  [Reagent-Example](https://github.com/vallard/Reagent-Example) with Reagent,
  Re-Frame, Secretary, Figwheel, and Compojure server.
* [11-frame](11-frame) - Re-frame template w/ DevTools, created with `lein new
  re-frame frame`. Could include other
  [profiles](https://github.com/Day8/re-frame-template).
* [12-svg](12-svg) - Extend the `11-frame` example with some SVG rendering and
  Pure CSS styling.
