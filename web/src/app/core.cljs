(ns app.core
  (:require [reagent.core :as r]
            [re-frame.core :as rf]
            [app.events]
            [app.subs]
            [app.views :as views]))

(defn mount-root []
  (rf/clear-subscription-cache!)
  (r/render [views/app-page] (.getElementById js/document "app")))

(def debug?
  ^boolean goog.DEBUG)

(defn ^:export init []
  (rf/dispatch-sync [:initialize-db])
  (when debug?
    (enable-console-print!)
    (println "dev mode"))
  (mount-root))
