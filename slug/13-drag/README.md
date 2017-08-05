# 13-svg

A simplified version of the `12-svg` SVG demo.

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
