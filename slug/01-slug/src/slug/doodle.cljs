(ns slug.doodle
  (:require-macros [devcards.core :refer [defcard
                                          defcard-doc
                                          defcard-rg
                                          deftest
                                          mkdn-pprint-source]]
                   [cljs.test :refer [testing is]])
  (:require [devcards.core]
            [reagent.core :as r]))

(defonce app-state (r/atom {:count 0}))

(defn on-click []
  (swap! app-state update-in [:count] inc))

(defn counter []
  [:div
    [:p "Current count: " (:count @app-state)]
    [:button {:on-click on-click} "Increment"]])

(defcard-doc
  "# My first devcard doodle"
  (mkdn-pprint-source counter))

(defcard-rg counter
  [counter])

(deftest silly-test
  (is (= 0 0)))

(defcard-doc
  (str "JavaScript date: *" (js/Date.) "*"))

(def patch
  [[:obj 75 40 :inlet]
   [:obj 75 91 :swap 123]
   [:obj 75 142 :print 1]
   [:obj 146 142 :print 2]
   [:connect 0 0 1 0]
   [:connect 1 0 2 0]
   [:connect 1 1 3 0]])

(defcard-doc
  (mkdn-pprint-source patch))

(defn as-svg [[_ x y] obj]
  [:rect {:x x :y y :width 30 :height 20
          :style {:outline "2px solid black" :fill "white"}}])

(defn svg-component []
  (into
    [:svg {:width "400" :height "200"
           :id "canvas"
           :style {:outline "1px solid black"}}
          [:circle {:cx 180 :cy 80 :r 30 :fill "red"}]
          [:rect {:x 250 :y 50 :width 50 :height 30 :fill "lightgreen"}]]
    (map as-svg (filter #(= (% 0) :obj) patch))))

(defcard-rg svg-test
  [svg-component])
