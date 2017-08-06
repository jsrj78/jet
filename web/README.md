# JET web interface

Derived from `slug/12-frame`.

### Development mode

Once [Leiningen](https://leiningen.org) has been installed, run:

```
lein figwheel
```

Wait for the `Prompt will show ...` message, then click on this link:
<http://localhost:3449](http://localhost:3449/>.

### Production build

To generate a self-contained production release:

```
lein clean
lein cljsbuild once min
```

The result ends up in `public/` and can be used with any HTTP server:

```
tree public/
public/
├── app.css
├── index.html
├── js
│   └── app.js
└── pure-min.css

1 directory, 4 files
```
