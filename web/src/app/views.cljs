(ns app.views
  (:require [re-frame.core]
            [clojure.string :as s]
            [goog.events :as ev]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-name [[_ _ _ & cmd]]
  (subs (s/join " " cmd) 1))

(defn obj-id [[_ x y]]
  (str x "," y))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn client-xy [evt]
  [(.-clientX evt) (.-clientY evt)])

(defn drag-move-fn [id state]
  (fn [evt]
    (let [[ox oy]       (:pos @state)
          [cx cy :as c] (client-xy evt)]
      (swap! state assoc :pos c)
      (>evt [:move-gadget id (- cx ox) (- cy oy)]))))

(defn drag-end-fn [move-fn state]
  (fn [evt]
    (ev/unlisten js/window "mousemove" move-fn)
    (ev/unlisten js/window "mouseup" (:end @state))))

(defn drag-start [x y evt]
  (let [state   (atom {:pos (client-xy evt)})
        id      (js/parseInt (.-id (.-target evt)))
        move-fn (drag-move-fn id state)
        done-fn (drag-end-fn move-fn state)]
    (swap! state assoc :end done-fn)
    (ev/listen js/window "mousemove" move-fn)
    (ev/listen js/window "mouseup" done-fn)))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn obj-as-svg [id [_ x y & cmd :as obj]]
  ^{:key id}
  [:g.draggable {:on-mouse-down #(drag-start x y %)}
    [:rect.obj {:id id :x x :y y :width 65 :height 20}]
    [:text.obj {:x (+ x 5) :y (+ y 15)} (obj-name obj)]])

(defn wire-id [wire]
  (s/join ":" wire))

(defn wire-path [x1 y1 x2 y2] ;; either straight line or cubic bezier
; (s/join " " ["M" x1 y1 "L" x2 y2])
  (s/join " " ["M" x1 y1 "C" x1 (+ y1 25) x2 (- y2 25) x2 y2]))

(defn wire-as-svg [[src-pos src-out dst-pos dst-in :as wire]]
  (let [[_ sx sy] (<sub [:gadget-num src-pos])
        [_ dx dy] (<sub [:gadget-num dst-pos])]
    ^{:key (wire-id wire)}
    [:path.wire {:d (wire-path (+ sx (* 65 src-out)) (+ sy 20)
                               (+ dx (* 65 dst-in))  (+ dy 0))}])) 

(defn design-as-svg []
  (let [objs   (<sub [:gadgets])
        wires  (<sub [:wires])]
    [:svg {:width 300 :height 200}
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
