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
