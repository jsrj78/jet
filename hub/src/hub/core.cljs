(ns hub.core
  (:require [cljs.nodejs :as node]))

(def restify (node/require "restify"))
(def fs (node/require "fs"))
(def bootdir "./bootimages")

(node/enable-util-print!)
(println "[hub.core]" bootdir)

(defn rest-logger [req res next]
  (println (.path req))
  (next))

(defn file-details [name]
  (let [stat (.statSync fs (str bootdir "/" name))]
    {:name name :size (.-size stat) :date (.-mtime stat)}))

(defn list-files [req res next]
  (let [dirlist (.readdirSync fs bootdir)]
    (.send res (clj->js (map file-details dirlist)))))

(defn add-file [req res next]
  (.log js/console 111 (.keys js/Object req))
  (.log js/console (.-headers req))
  (.log js/console 222 (.-body req))
  (.log js/console 333 (.-params req))
  (.send res 200))

(defn create-server []
  (let [static-server (.serveStatic restify #js {:directory bootdir})]
    (doto (.createServer restify)
      (.use (.CORS restify))
      (.use (.bodyParser restify))
      (.use rest-logger)
      (.get "/" list-files)
      (.get "/index.txt" static-server)
      (.get #"^/.+\.bin$" static-server)
      (.post #"^/.+" add-file))))

(defn -main []
  (println "Hello world!")
  (let [server (create-server)]
    (.listen server 3000)))

(set! *main-cli-fn* -main)
