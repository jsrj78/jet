(ns app.views
  (:require [re-frame.core]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-as-svg [[_ x y]]
  ^{:key (+ (* x 10000) y)}
  [:rect.obj {:x x :y y :width 30 :height 20}])

(defn design-as-svg [design]
  (let [obj-v (filterv #(= (first %) :obj) design)] 
    [:svg {:width 400 :height 200}
      [:circle.test {:cx 180 :cy 80 :r 30}]
      [:rect.test {:x 250 :y 50 :width 50 :height 30}]
      (map obj-as-svg obj-v)]))

(defn main-panel []
  (let [my-name (<sub [:name]) 
        design (<sub [:design])]
    [:div#main

      [:div.pure-g.pure-u-1
        [:h3 "Hello " my-name]
        [:input  {:type "text"
                  :value my-name
                  :on-change #(>evt [:change-name (.. % -target -value)])}]]

      [:p.pure-g.pure-u-1
        [design-as-svg design]]

      [:small.pure-g.pure-u-1
        (pr-str @re-frame.db/app-db)]]))
