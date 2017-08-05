(ns app.views
  (:require [re-frame.core]
            [clojure.string :as s]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-name [[_ _ _ & cmd]]
  (subs (s/join " " cmd) 1))

(defn obj-id [[_ x y]]
  (+ x (* 10000 y)))

(defn obj-as-svg [[_ x y & cmd :as obj]]
  ^{:key (obj-id obj)}
  [:g.obj {:on-click #(.log js/console "click:" (.-target %))}
    [:rect.obj {:x x :y y :width 65 :height 20}]
    [:text.obj {:x (+ x 5) :y (+ y 15)} (obj-name obj)]])

(defn wire-id [[_ & args]]
  (s/join "," args))

(defn wire-as-svg [ovec [_ src-pos src-out dst-pos dst-in :as wire]]
  (let [[_ sx sy] (nth ovec src-pos)
        [_ dx dy] (nth ovec dst-pos)]
    ^{:key (wire-id wire)}
    [:line.wire {:x1 (+ sx (* 65 src-out)) :y1 (+ sy 20)
                 :x2 (+ dx (* 65 dst-in))  :y2 (+ dy 0)}]))

(defn design-as-svg []
  (let [design (<sub [:design])
        objs (filterv #(= (first %) :obj) design)
        wires (filter #(= (first %) :connect) design)] 
    [:svg {:width 400 :height 200}
      (map obj-as-svg objs)
      (map (partial wire-as-svg objs) wires)]))

(defn main-panel []
  [:div#main

    [:p.pure-g.pure-u-1
      "Hello " [:b "SVG"]]

    [:div.pure-g.pure-u-1
      [design-as-svg]]

    [:p.pure-g.pure-u-1
      [:small (pr-str @re-frame.db/app-db)]]])
