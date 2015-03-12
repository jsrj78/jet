(ns blink3.core
  (:require [om.core :as om :include-macros true]
            [om.dom :as dom :include-macros true]))

(defonce app-state (atom {:text "JET Blink 3"}))

(defn main []
  (om/root
    (fn [app owner]
      (reify
        om/IRender
        (render [_]
          (dom/div nil
            (dom/h2 nil (:text app))
            (dom/input #js {:type "checkbox"} "Enable")))))
    app-state
    {:target (. js/document (getElementById "app"))}))
