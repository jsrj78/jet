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

(defn list-files []
  (filter #(not= % "index.txt") (.readdirSync fs bootdir)))

(defn send-file-list [req res next]
  (let [withinfo (fn [name]
                   (let [stat (.statSync fs (str bootdir "/" name))]
                     {:name name :size (.-size stat) :date (.-mtime stat)}))]
    (.send res (clj->js (map withinfo (list-files))))))

(defn parse-leading-int [str]
  (let [digits (re-find  #"\d+" str)]
    (if digits (js/parseInt digits 10))))

(defn highest-seqnum []
  (apply max 999 (map parse-leading-int (list-files))))

(println (highest-seqnum))

(defn add-file [req res next]
  (let [id (inc (highest-seqnum))
        params (.-params req)
        bytes (.-bytes params)]
    (.log js/console 333 id params (identity bytes))
    (.writeFileSync fs (str bootdir "/" id ".bin") (identity bytes))
    (.send res 200)))

(defn create-server []
  (let [static-server (.serveStatic restify #js {:directory bootdir})]
    (doto (.createServer restify)
      (.use (.CORS restify))
      (.use (.bodyParser restify))
      (.use rest-logger)
      (.get "/" send-file-list)
      (.get "/index.txt" static-server)
      (.get #"^/.+\.bin$" static-server)
      (.post #"^/.+\.bin$" add-file))))

(defn -main []
  (println "Hello world!")
  (let [server (create-server)]
    (.listen server 3000)))

(set! *main-cli-fn* -main)
