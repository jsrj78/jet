# JET web interface

The web interface is built with [ClojureScript][CLJS] (which compiles to
JavaScript) and [Reagent][REAG] (which is a lightweight wrapper around ReactJS).

This application was created using `lein new figwheel web -- --reagent`.

To launch this code in development mode, you need [Leiningen][LEIN] and a JVM:

    lein figwheel

The following command generates an optimised client-side-only static build:

    lein do clean, cljsbuild once min

The result will be a static set of files in `resources/public/`.

   [CLJS]: https://clojurescript.org
   [REAG]: https://reagent-project.github.io
   [LEIN]: https://leiningen.org
