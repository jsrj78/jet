(ns app.core
  (:require [reagent.core :as reagent]
            [re-frame.core :as rf]
            [app.events]
            [app.subs]
            [app.views :as views]))

(defn mount-root []
  (rf/clear-subscription-cache!)
  (reagent/render [views/app-page] (.getElementById js/document "app")))

(def debug?
  ^boolean goog.DEBUG)

(defn ^:export init []
  (when debug?
    (enable-console-print!)
    (println "dev mode"))
  (rf/dispatch-sync [:initialize-db])
  (mount-root))
