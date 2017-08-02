(ns purecss.core
  (:require
   [reagent.core :as reagent]))

(defonce app-state
  (reagent/atom {}))

(defn page [ratom]
  [:div.pure-g
   [:div.pure-u-1-5]
   [:div.pure-u-2-5
    [:h3 "Welcome to Reagent with Pure CSS!"]]
   [:div.pure-u-2-5
    [:p "Some more text..."]]])

(defn dev-setup []
  (when ^boolean js/goog.DEBUG
    (enable-console-print!)
    (println "dev mode")))

(defn reload []
  (reagent/render [page app-state]
                  (.getElementById js/document "app")))

(defn ^:export main []
  (dev-setup)
  (reload))
