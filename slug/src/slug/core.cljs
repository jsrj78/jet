(ns slug.core
  (:require
   [reagent.core :as r]))

(defonce app-state (r/atom {}))

(defn page [ratom]
  [:h1 "Welcome to JET/Slug"])

(defn dev-setup []
  (when ^boolean js/goog.DEBUG
    (enable-console-print!)
    (println "dev mode")))

(defn reload []
  (r/render [page app-state] (.getElementById js/document "app")))

(defn ^:export main []
  (dev-setup)
  (reload))
