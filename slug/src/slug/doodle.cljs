(ns slug.doodle
  (:require-macros [devcards.core :refer [defcard
                                          defcard-doc
                                          defcard-rg
                                          deftest
                                          mkdn-pprint-source]]
                   [cljs.test :refer [testing is]])
  (:require [devcards.core]
            [reagent.core :as reagent]))

(defonce app-state (reagent/atom {:count 0}))

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
  [[:obj 75 101 :swap 123]
   [:obj 75 142 :print 1]
   [:obj 146 142 :print 2]
   [:obj 75 60 :inlet]
   [:connect 0 0 1 0]
   [:connect 0 1 2 0]
   [:connect 3 0 0 0]])

(defcard-doc
  (mkdn-pprint-source patch))

(defn svg-component []
  [:svg {:width "400" :height "250"
         :id "canvas"
         :style {:outline "1px solid black"}}
    [:circle {:cx 100 :cy 100 :r 30 :fill "red"}]
    [:rect {:x 250 :y 150 :width 50 :height 30 :fill "lightgreen"}]])

(defcard-rg svg-test
  [svg-component])
