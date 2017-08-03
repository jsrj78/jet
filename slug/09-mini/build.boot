(set-env!
 :source-paths    #{"src"}
 :resource-paths  #{"resources"}
 :dependencies '[[adzerk/boot-cljs          "2.0.0"      :scope "test"]
                 [adzerk/boot-cljs-repl     "0.3.3"      :scope "test"]
                 [adzerk/boot-reload        "0.5.1"      :scope "test"]
                 [pandeiro/boot-http        "0.8.3"      :scope "test"]
                 [com.cemerick/piggieback   "0.2.1"      :scope "test"]
                 [org.clojure/tools.nrepl   "0.2.13"     :scope "test"]
                 [weasel                    "0.7.0"      :scope "test"]
                 [org.clojure/clojurescript "1.9.562"]
                 [rum                       "0.10.7"]])

(require
 '[adzerk.boot-cljs      :refer [cljs]]
 '[adzerk.boot-cljs-repl :refer [cljs-repl start-repl]]
 '[adzerk.boot-reload    :refer [reload]]
 '[pandeiro.boot-http    :refer [serve]])

(deftask development []
  (task-options! cljs {:optimizations :none}
                 reload {:on-jsload 'mini.app/init})
  identity)

(deftask production []
  (task-options! cljs {:optimizations :advanced})
  identity)

(deftask dev
  "Start up in development mode"
  []
  (comp (development)
        (serve)
        (watch)
        (cljs-repl)
        (reload)
        (cljs)))

(deftask prod
  "Build in production mode"
  []
  (comp (production)
        (cljs)
        (target)))

;; vim: ft=clojure :
