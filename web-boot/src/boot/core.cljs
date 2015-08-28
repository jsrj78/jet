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
  (ajax/GET url {:handler #(swap! app-state assoc :index (parse-index %))}))

(defn get-files [url]
  (ajax/GET url {:handler #(swap! app-state assoc :files %)}))

(defn index-row [x]
  [:tr [:td [:code (x 0)]] [:td (x 1)] [:td (x 2)]])

(defn index-table []
  [:table
   [:thead
    [:tr [:th "Hardware ID"] [:th "S/W ID"] [:th "Filename"]]]
   [:tbody
    (for [x (:index @app-state)]
      ^{:key x} [index-row x])]])

(defn files-row [x]
  ;; TODO should change string keys to keywords during get-files reception
  (println x)
  [:tr [:td [:code (x "name")]] [:td (x "size")] [:td (x "date")]])

(defn files-table []
  [:table
   [:thead
    [:tr [:th "Filename"] [:th "Size"] [:th "Date"]]]
   [:tbody
    (for [x (:files @app-state)]
      ^{:key x} [files-row x])]])

(defn hello-world []
  [:div
   [:h1 "JeeBoot Configuration"]
   [:div
    [:h3 "Node map"]
    [index-table]]
   [:div
    [:h3 "Files"]
    [files-table]]
   [:p (:text @app-state)]])

(defn on-js-reload []
  ;; optionally touch app-state to force rerendering depending on the app
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
  )

(get-index "http://localhost:3000/index.txt")
(get-files "http://localhost:3000/")

(reagent/render-component [hello-world] (. js/document (getElementById "app")))
