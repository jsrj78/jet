(ns frame.views
  (:require [re-frame.core :as rf]))

(defn doodle-panel []
  (let [name (rf/subscribe [:name])]
    (fn []
      [:div
       "Hello from: " [:b @name]
       [:hr]
       [:button {:on-click #(rf/dispatch [:change-name "abc"])}
                "Change name"]])))

(defn main-panel []
  doodle-panel)
