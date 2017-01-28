(ns web.core
  (:require [reagent.core :as reagent :refer [atom]]))

(enable-console-print!)

(println "This text is printed from src/web/core.cljs.")

;; define your app data so that it doesn't get over-written on reload

(defonce app-state (atom {:text "Hello from JET"}))

(defn hello-world []
  [:h1.red (:text @app-state)])

(reagent/render-component [hello-world]
                          (. js/document (getElementById "app")))

(defn on-js-reload []
  ;; optionally touch your app-state to force rerendering depending on
  ;; your application
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
)

