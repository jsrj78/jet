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
  [[:rect {:x x :y y :width 60 :height 20
           :style {:outline "1px solid black" :fill "white"}}]
   [:text {:x (+ x 1) :y (+ y 15)} (str cmd)]])

(defn svg-wire [ovec [src-pos src-outlet dst-pos dst-inlet]]
  (let [[sx sy] (get ovec src-pos)
        [dx dy] (get ovec dst-pos)]
    [:line {:x1 (+ sx (* 60 src-outlet))
            :y1 (+ sy 20)
            :x2 (+ dx (* 60 dst-inlet))
            :y2 (+ dy 0)
            :stroke-width "1" :stroke "black"}]))

(defn objects [design]
  (let [obj-vec (atom [])
        results (atom [])]
    (doseq [[t & r] design]
      (case t
        :obj (do (swap! obj-vec conj r)
                 (swap! results into (svg-obj r)))
        :connect (swap! results conj (svg-wire @obj-vec r))))
    (.log js/console "obj-vec:" (str @obj-vec))
    @results))

(defn page [ratom]
  (let [o (objects (:design @app-state))]
    [:div
     [:h3 "Welcome to JET/Slug"]
     [:table
      [:tbody
       [:tr
        [:td {:width "50%"
              :style {:vertical-align "top"}} [into (svg) o]]
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
