# 11-frame

A minimal [re-frame](https://github.com/Day8/re-frame) application:

```
$ tree
.
├── README.md
├── project.clj
├── public
│   ├── app.css
│   └── index.html
└── src
    └── app
        ├── core.cljs
        ├── db.cljs
        ├── events.cljs
        ├── subs.cljs
        └── views.cljs

3 directories, 9 files
```

## Development Mode

```
lein clean
lein figwheel dev
```

Figwheel will automatically push cljs changes to the browser.

Wait a bit, then browse to [http://localhost:3449](http://localhost:3449).

## Production Build

To compile clojurescript to javascript:

```
lein clean
lein cljsbuild once min
```
