(ns app.views
  (:require [re-frame.core]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-as-svg [[_ x y & cmd]]
  ^{:key (+ (* x 10000) y)}
  [:g
    [:rect.obj {:x x :y y :width 70 :height 20}]
    [:text {:x (+ x 2) :y (+ y 15)} (str cmd)]])

(defn wire-as-svg [ovec [_ src-pos src-outlet dst-pos dst-inlet]]
  (let [[_ sx sy] (get ovec src-pos)
        [_ dx dy] (get ovec dst-pos)]
    ^{:key (str src-pos "," src-outlet "," dst-pos "," dst-pos)}
    [:line {:x1 (+ sx (* 70 src-outlet)) :y1 (+ sy 20)
            :x2 (+ dx (* 70 dst-inlet))  :y2 (+ dy 0)
            :stroke "black"}]))

(defn design-as-svg [design]
  (let [obj-v (filterv #(= (first %) :obj) design)
        wire-v (filterv #(= (first %) :connect) design)] 
    [:svg {:width 400 :height 200}
      (map obj-as-svg obj-v)
      (map #(wire-as-svg obj-v %) wire-v)]))

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
