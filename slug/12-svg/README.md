# 12-svg

Adapted from the `11-frame` minimal re-frame application.

## Development Mode

```
lein figwheel
```

Wait a bit, then browse to [http://localhost:3449](http://localhost:3449).

## Production Build

To compile clojurescript to javascript:

```
lein do clean, cljsbuild once min
```
