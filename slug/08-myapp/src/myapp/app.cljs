(ns myapp.app
  (:require [reagent.core :as reagent :refer [atom]]))

(defn some-component []
  [:div
   [:h3 "I am a component!"]
   [:p.someclass
    "I have " [:strong "bold"]
    [:span {:style {:color "red"}} " and red"]
    " text."]])

(defn calling-component []
  [:div "Parent component"
   [some-component]])

(defn init []
  (.log js/console (type (range 10)))
  (.log js/console (range 10))

  (reagent/render-component [calling-component]
                            (.getElementById js/document "container")))
