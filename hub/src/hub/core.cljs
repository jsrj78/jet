(ns hub.core
  (:require [cljs.nodejs :as node]))

(def restify (node/require "restify"))

(node/enable-util-print!)
(println "[hub.core]")

(defn rest-logger [req res next]
  (println (.path req))
  (next))

(defn create-server []
  (let [server (.createServer restify)
        static-server (.serveStatic restify #js {:directory "./bootimages"})]
    (.use server rest-logger)
    (.use server (.CORS restify))
    (.get server "/index.txt" static-server)
    (.get server #"^/.*\.bin$" static-server)
    server))

(defn -main []
  (println "Hello world!")
  (let [server (create-server)]
    (.listen server 3000)))

(set! *main-cli-fn* -main)
