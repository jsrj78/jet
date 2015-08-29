(ns ^:figwheel-always boot.core
    (:require [reagent.core :as reagent]
              [ajax.core :as ajax]
              [clojure.string :as str]))

(enable-console-print!)
(println "[jet/boot]")

(def server-url "http://localhost:3000/")

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

(defn get-index []
  (let [url (str server-url "index.txt")]
    (ajax/GET url {:handler #(swap! app-state assoc :index (parse-index %))})))

(defn get-files []
  (ajax/GET server-url {:handler #(swap! app-state assoc :files %)}))

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
  [:tr [:td [:code (x "date")]] [:td (x "size")] [:td (x "name")]])

(defn files-table []
  [:table
   [:thead
    [:tr [:th "Date"] [:th "Size"] [:th "Filename"]]]
   [:tbody
    (for [x (sort #(compare (%2 "date") (%1 "date")) (:files @app-state))]
      ^{:key x} [files-row x])]])

(defn upload-file [file]
  (let [name (.-name file)
        date (.-lastModified file)
        reader (js/FileReader.)]
    (set! (.-onloadend reader)
          (fn []
            (let [bytes (.-result reader)
                  desc {:name name :date date :bytes (js/btoa bytes)}
                  url (str server-url name)]
              (ajax/POST url {:format :json
                              :params desc
                              :handler get-files}))))
    (.readAsBinaryString reader file)))

(defn drop-handler [evt]
  (.preventDefault evt)
  (let [files (.. evt -dataTransfer -files)]
    (dotimes [i (.-length files)]
      (upload-file (.item files i)))))

(defn hello-world []
  [:div
   [:h1 "JeeBoot Configuration"]
   [:div
    [:h3 "Node map"]
    [index-table]]
   [:div
    [:h3 "Firmware images"]
    [:div {:id "drop" :on-drop drop-handler
           :on-drag-over #(.preventDefault %)}
     "Drop new files here ... (only *.bin accepted)"]
    [files-table]]
   [:p (:text @app-state)]])

(defn on-js-reload []
  ;; optionally touch app-state to force rerendering depending on the app
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
  )

(get-index)
(get-files)

(defn get-by-id [id]
  (. js/document (getElementById id)))

(reagent/render-component [hello-world] (get-by-id "app"))
