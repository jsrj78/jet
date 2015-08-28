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

(defn create-server []
  (let [server (.createServer restify)
        static-server (.serveStatic restify #js {:directory bootdir})]
    (.use server rest-logger)
    (.use server (.CORS restify))
    (.get server "/" list-files)
    (.get server "/index.txt" static-server)
    (.get server #"^/.*\.bin$" static-server)
    server))

(defn -main []
  (println "Hello world!")
  (let [server (create-server)]
    (.listen server 3000)))

(set! *main-cli-fn* -main)
