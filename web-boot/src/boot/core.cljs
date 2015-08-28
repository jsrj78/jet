(ns ^:figwheel-always boot.core
    (:require [reagent.core :as reagent]
              [ajax.core :as ajax]
              [clojure.string :as str]))

(enable-console-print!)
(println "[jet/boot]")

(defonce app-state (reagent/atom {:text "Hello world!"}))

(defn parse-one [line]
  (let [[_ hw sw nm] (re-find #"^([0-9A-F]{32})\s*=\s*(\d+)\s*(\S+)?" line)]
    (if (and hw sw)
      [hw sw nm])))

(defn parse-index [text]
  (->> text
       str/split-lines
       (map parse-one)
       (filter identity)))

(defn get-index [url]
  (ajax/GET url {:handler
                 (fn [res]
                   (swap! app-state assoc :index (parse-index res)))}))

(defn index-item [x]
  [:tr [:td [:code (x 0)]] [:td (x 1)] [:td (x 2)]])

(defn index-list []
  [:table
   [:thead
    [:tr [:th "Hardware ID"] [:th "S/W ID"] [:th "Filename"]]]
   [:tbody
    (for [x (:index @app-state)]
      ^{:key x} [index-item x])]])

(defn hello-world []
  [:div
   [:h1 "JeeBoot Configuration"]
   [index-list]
   [:p (:text @app-state)]])

(defn on-js-reload []
  ;; optionally touch app-state to force rerendering depending on the app
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
  )

(get-index "http://localhost:3000/index.txt")

(reagent/render-component [hello-world] (. js/document (getElementById "app")))
