(ns mini.app
    (:require [rum.core :as rum]))

(rum/defc label [text]
  [:div
   [:h2 "A label"]
   [:p text]])

(defn init []
  (rum/mount (label "blah") (. js/document (getElementById "app"))))
