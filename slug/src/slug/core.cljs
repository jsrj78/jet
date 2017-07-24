(ns slug.core
  (:require
   [reagent.core :as r]))

(defonce app-state (r/atom {:design [[:obj 75 40 :inlet]
                                     [:obj 75 91 :swap 123]
                                     [:obj 75 142 :print 1]
                                     [:obj 146 142 :print 2]
                                     [:connect 0 0 1 0]
                                     [:connect 1 0 2 0]
                                     [:connect 1 1 3 0]]}))

(defn svg []
  [:svg {:width "600" :height "400"
         :style {:outline "1px solid black"}}
   [:circle {:cx 300 :cy 200 :r 30 :fill "red"}]])

(defn svg-obj [[x y & cmd]]
  [:rect {:x x :y y :width 50 :height 20}])

(defn svg-wire [[from-obj from-outlet to-obj to-inlet]]
  [:line {:x1 from-obj :y1 from-outlet :x2 to-obj :y2 to-inlet}])

(defn objects [design]
  (let [gadgets (atom [])
        wires   (atom [])]
    (doseq [[t & r] design]
      (case t
        :obj     (swap! gadgets conj (svg-obj r))
        :connect (swap! wires   conj (svg-wire r))))
    (concat @gadgets @wires)))

(defn page [ratom]
  (let [o (objects (:design @app-state))]
    [:div
     [:h3 "Welcome to JET/Slug"]
     [:table
      [:tbody
       [:tr
        [:td {:width "50%"} (svg)]
        [:td {:width "1%"}]
        [:td {:style {:vertical-align "top"}}
         [:pre (with-out-str (cljs.pprint/pprint @app-state))]
         [:hr]
         [:pre (with-out-str (cljs.pprint/pprint o))]]]]]]))


(defn dev-setup []
  (when ^boolean js/goog.DEBUG
    (enable-console-print!)
    (println "dev mode")))

(defn reload []
  (r/render [page app-state] (.getElementById js/document "app")))

(defn ^:export main []
  (dev-setup)
  (reload))
