(ns ^:figwheel-always boot.core
    (:require [reagent.core :as reagent]
              [ajax.core :as ajax]
              [clojure.string :as str]))

(enable-console-print!)
(println "[jet/boot]")

(def server-url "http://localhost:3000/")

(defonce app-state (reagent/atom {}))

(defn parse-one [line]
  (let [re #"^([0-9A-F]{32}) ([-+]) (\d*) ?(\d*)$"
        [_ hwid new id prev] (re-find re line)]
    (if hwid [hwid (= new "+") id prev])))

(defn parse-index [text]
  (->> text
       str/split-lines
       (map parse-one)
       (filter identity)
       vec))

(defn get-index []
  (let [url (str server-url "index.txt")]
    (ajax/GET url {:handler #(swap! app-state assoc :index (parse-index %))})))

(defn get-files []
  (ajax/GET server-url {:handler #(swap! app-state assoc :files %)}))

(defn delete-file! [file]
  (ajax/DELETE (str server-url file) {:handler get-files}))

(defn change-index! [newindex]
  (swap! app-state assoc :index newindex)
  (let [as-text (fn [[hwid new id prev]]
                  (str hwid (if new " + " " - ") id " " prev "\n"))
        newtext (str/join (map as-text newindex))]
    (ajax/POST (str server-url "index.txt") {:format :json
                                             :params {:text newtext}
                                             :handler get-index})))

(defn toggle-new! [hwid]
  (change-index! (mapv (fn [item]
                         (if (= (item 0) hwid)
                           (assoc item 1 (not (item 1)))
                           item))
                       (:index @app-state))))

(defn index-row [x]
  [:tr
   [:td [:code (x 0)]]
   [:td [:input {:type "checkbox"
                 :checked (if (x 1) "checked")
                 :on-change #(toggle-new! (x 0))}]]
   [:td (x 2)]
   [:td (x 3)]])

(defn index-table []
  [:table
   [:thead
    [:tr [:th "Hardware ID"] [:th "New"] [:th "ID"] [:th "Previous"]]]
   [:tbody
    (for [x (:index @app-state)]
      ^{:key x} [index-row x])]])

(defn id-unused? [id]
  (= 0 (count (filter #(or (= id (% 2)) (= id (% 3))) (:index @app-state)))))

(defn files-row [x]
  ;; TODO should change string keys to keywords during get-files reception
  (let [name (x "name")
        size (x "size")
        date (x "date")
        id (re-find #"\d+" name)]
    [:tr
     [:td date]
     [:td size]
     [:td id]
     [:td (if (id-unused? id)
            [:button {:on-click #(delete-file! (x "name"))} "delete"])]]))

(defn files-table []
  [:table
   [:thead
    [:tr [:th "Date"] [:th "Size"] [:th "ID"]]]
   [:tbody
    (for [x (sort #(compare (%2 "date") (%1 "date")) (:files @app-state))]
      ^{:key x} [files-row x])]])

(defn update-index! [req res next]
  (get-files)
  (change-index! (mapv (fn [[hwid new id prev]]
                         (if new
                           [hwid new (req "id") id]
                           [hwid new id prev]))
                       (:index @app-state))))

(defn upload-file! [file]
  (let [name (.-name file)
        date (.-lastModified file)
        reader (js/FileReader.)
        on-done (fn []
                  (let [bytes (.-result reader)
                        desc {:name name :date date :bytes (js/btoa bytes)}
                        url (str server-url name)]
                    (ajax/POST url {:format :json
                                    :params desc
                                    :handler update-index!})))]
    (set! (.-onloadend reader) on-done)
    (.readAsBinaryString reader file)))

(defn drop-handler [evt]
  (.preventDefault evt)
  (let [files (.. evt -dataTransfer -files)]
    (dotimes [i (.-length files)]
      (upload-file! (.item files i)))))

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
    [files-table]]])

(defn on-js-reload []
  ;; optionally touch app-state to force rerendering depending on the app
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
  )

(get-index)
(get-files)

(defn get-by-id [id]
  (. js/document (getElementById id)))

(reagent/render-component [hello-world] (get-by-id "app"))
