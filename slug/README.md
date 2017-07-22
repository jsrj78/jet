# JET/Slug

And exploration into the world of ClojureScript, Figwheel, and Reagent.

[![license](https://img.shields.io/github/license/jeelabs/jet.svg)](http://unlicense.org)

App code is `src/slug/core.cljs`, first devcard is `src/slug/doodle.cljs`.

Live app, see <http://localhost:3449/>:

    lein figwheel

Live devcards, see <http://localhost:3449/cards.html>:

    lein figwheel devcards

Static devcards, see `resources/public/cards.html`:

    lein cljsbuild once hostedcards

Static app, see `resources/public/index.html`:

    lein cljsbuild once min

Clean up temp files:

    lein clean

See https://github.com/gadfly361/reagent-figwheel on GitHub.
