(ns app.views
  (:require [re-frame.core]
            [clojure.string :as s]
            [goog.events :as events])
  (:import [goog.events EventType]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-name [[_ _ _ & cmd]]
  (subs (s/join " " cmd) 1))

(defn obj-id [[_ x y]]
  (str x "," y))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn bounding-client-xy [evt]
  (let [rect (.getBoundingClientRect (.-target evt))]
    [(.-left rect) (.-top rect)]))

(defn client-xy [evt]
  [(.-clientX evt) (.-clientY evt)])

(defn drag-move-fn [gid]
  (fn [evt]
    (let [[bx by] (bounding-client-xy evt)
          [cx cy] (client-xy evt)]
      (>evt [:move-obj gid (- cx bx) (- cy by)]))))

(defn drag-end-fn [drag-move drag-end-atom]
  (fn [evt]
    (events/unlisten js/window EventType.MOUSEMOVE drag-move)
    (events/unlisten js/window EventType.MOUSEUP @drag-end-atom)))

(defn drag-start [x y evt]
  (.log js/console "ds:" (client-xy evt) (.-id (.-target evt)))
  (let [gid (js/parseInt (.-id (.-target evt)))
        drag-move (drag-move-fn gid)
        drag-end-atom (atom nil)
        drag-end (drag-end-fn drag-move drag-end-atom)]
    (reset! drag-end-atom drag-end)
    (events/listen js/window EventType.MOUSEMOVE drag-move)
    (events/listen js/window EventType.MOUSEUP drag-end)))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn obj-as-svg [idx [_ x y & cmd :as obj]]
  ^{:key idx}
  [:g.draggable {:on-mouse-down #(drag-start x y %)}
    [:rect.obj {:id idx :x x :y y :width 65 :height 20}]
    [:text.obj {:x (+ x 5) :y (+ y 15)} (obj-name obj)]])

(defn wire-id [[_ & args]]
  (s/join ":" args))

(defn wire-as-svg [[_ src-pos src-out dst-pos dst-in :as wire]]
  (let [[_ sx sy] (<sub [:gadget-num src-pos])
        [_ dx dy] (<sub [:gadget-num dst-pos])]
    ^{:key (wire-id wire)}
    [:line.wire {:x1 (+ sx (* 65 src-out)) :y1 (+ sy 20)
                 :x2 (+ dx (* 65 dst-in))  :y2 (+ dy 0)}]))

(defn design-as-svg []
  (let [design (<sub [:design])
        objs   (filter #(= (first %) :obj) design)
        wires  (filter #(= (first %) :connect) design)] 
    [:svg {:width 400 :height 200}
      (map-indexed obj-as-svg objs)
      ; can't leave reactive refs in a lazy sequence
      (doall (map wire-as-svg wires))]))

(defn main-panel []
  [:div#main

    [:p.pure-g.pure-u-1
      "Hello " [:b "SVG"]]

    [:div.pure-g.pure-u-1
      [design-as-svg]]

    [:p.pure-g.pure-u-1
      [:small (pr-str @re-frame.db/app-db)]]])

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(comment "unused code:"

  (defn by-id [id]
    (.getElementById js/document id))

  (defn obj-id-as-xy [id]
    (mapv js/parseInt (s/split id ","))))
