# JET/Slug

My explorations into the world of ClojureScript, Figwheel, and Reagent.  
Each "slug" is a self-contained project where I _slowly_ try something new:

* [01-slug](01-slug) - Set up a Leiningen project with Figwheel (live coding),
  Reagent (React wrapper), and DevCards (web-based test environment). Based on
  the [reagent-figwheel](https://github.com/gadfly361/reagent-figwheel)
  template.
* [02-purecss](02-purecss) - A basic reagent setup with [Pure
  CSS](https://purecss.io) as layout framework.
* [03-paho](03-paho) - Minimal MQTT client in pure JavaScript, example from the
  [Eclipse Paho](https://www.eclipse.org/paho/clients/js/) website.
* [04-rum](04-rum) - Simple [Rum](https://github.com/tonsky/rum) example, which
  is a wrapper around React similar to Reagent.  Project created with
  [Chestnut](https://github.com/plexus/chestnut).
* [05-cljsjs](05-cljsjs) - Running the Paho MQTT client from ClojureScript,
  using [this cljsjs
  example](https://github.com/cljsjs/packages/tree/master/paho). Project created
  with `lein new figwheel app`, which picked up all the latest packages.
* [06-boot](06-boot) - Set up a reagent app via [Boot](http://boot-clj.com),
  using `boot -d boot/new new -t tenzing -n your-app -a +reagent`.
* [07-devtools](07-devtools) - Demo of the
  [cljs-devtools](https://github.com/binaryage/cljs-devtools/blob/master/docs/installation.md)
  package, to show much better debugging info in the JavaScript console.
* [08-myapp](08-myapp) - A fairly comfortable dev setup with Reagent and
  Devtools.  Uses Boot, supports a browser REPL.
