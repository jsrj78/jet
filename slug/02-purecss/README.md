# 02-purecss


A [reagent](https://github.com/reagent-project/reagent) application
with [Pure CSS](https://purecss.io) as layout framework.

## Development Mode

```
lein clean
lein figwheel dev
```

Figwheel will automatically push cljs changes to the browser.  
Wait a bit, then browse to [http://localhost:3449](http://localhost:3449).

## Production Build

```
lein clean
lein cljsbuild once min
```
