(ns hub.core
  (:require [cljs.nodejs :as node]))

(def restify (node/require "restify"))
(def fs (node/require "fs"))
(def bootdir "./bootimages")

(node/enable-util-print!)
(println "[hub.core]" bootdir)

(defn rest-logger [req res next]
  (println (.-method req) (.path req))
  (next))

(defn list-files []
  (filter #(not= % "index.txt") (.readdirSync fs bootdir)))

(defn send-file-list [req res next]
  (let [withinfo (fn [name]
                   (let [stat (.statSync fs (str bootdir "/" name))]
                     {:name name :size (.-size stat) :date (.-mtime stat)}))]
    (.send res (clj->js (map withinfo (list-files))))))

(defn parse-leading-int [s]
  (let [digits (re-find  #"\d+" s)]
    (if digits (js/parseInt digits 10))))

(defn update-index! [req res next]
  (.writeFileSync fs (str bootdir "/index.txt") (.. req -params -text))
  (.send res 200))

(defn highest-seqnum []
  (apply max 999 (map parse-leading-int (list-files))))

(defn add-file! [req res next]
  (let [id (inc (highest-seqnum))
        params (.-params req)
        name (.-name params)
        bytes (js/Buffer. (.-bytes params) "base64")]
    (.writeFileSync fs (str bootdir "/" id ".bin") bytes)
    (.send res 200 #js {:id id})))

(defn delete-file! [req res next]
  (.unlinkSync fs (str bootdir (.-url req)))
  (.send res 200))

(defn create-server []
  (let [static-server (.serveStatic restify #js {:directory bootdir})]
    (doto (.createServer restify)
      (.use (.CORS restify))
      (.use (.bodyParser restify))
      (.use rest-logger)
      (.get "/" send-file-list)
      (.get "/index.txt" static-server)
      (.post "/index.txt" update-index!)
      (.get #"^/.+\.bin$" static-server)
      (.post #"^/.+\.bin$" add-file!)
      (.del #"^/.+\.bin$" delete-file!))))

(defn -main []
  (println "server starting on http://localhost:3000/ ...")
  (let [server (create-server)]
    (.listen server 3000)))

(set! *main-cli-fn* -main)
