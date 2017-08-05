(ns app.core
  (:require [reagent.core :as r]
            [goog.events :as ev]))

(defonce app-db (r/atom {:obj [[50 40 "inlet #0"]
                               [50 90 "swap 5"]
                               [50 140 "print 1"]
                               [150 140 "print 2"]]
                         :wire [[0 0 1 0]
                                [1 0 2 0]
                                [1 1 3 0]]}))

(defn obj-num [id]
  (get-in @app-db [:obj id]))

(defn move-obj [id dx dy]
  (swap! app-db update-in [:obj id]
                (fn [[vx vy vtail]]
                  [(+ vx dx) (+ vy dy) vtail])))

(defn client-xy [evt]
  [(.-clientX evt) (.-clientY evt)])

(defn drag-move-fn [id state]
  (fn [evt]
    (let [[ox oy]       (:pos @state)
          [cx cy :as c] (client-xy evt)]
      (swap! state assoc :pos c)
      (move-obj id (- cx ox) (- cy oy)))))

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

(defn obj-as-svg [id [x y cmd]]
  ^{:key id}
  [:g.draggable {:on-mouse-down #(drag-start x y %)}
    [:rect.obj {:id id :x x :y y :width 65 :height 20}]
    [:text.obj {:x (+ x 5) :y (+ y 15)} cmd]])

(defn wire-as-svg [[src-id src-out dst-id dst-in :as wire]]
  (let [[sx sy] (obj-num src-id)
        [dx dy] (obj-num dst-id)]
    ^{:key (str wire)}
    [:line.wire {:x1 (+ sx (* 65 src-out)) :y1 (+ sy 20)
                 :x2 (+ dx (* 65 dst-in))  :y2 (+ dy 0)}]))

(defn main []
  [:div
    [:h1 "Hello SVG"]
    [:svg {:width 300 :height 200}
      (map-indexed obj-as-svg (:obj @app-db))
      ; don't leave reactive refs inside a lazy sequence
      (doall (map wire-as-svg (:wire @app-db)))]])

(def debug? ^boolean goog.DEBUG)

(defn mount-root []
  (r/render [main] (.getElementById js/document "app")))

(defn ^:export init []
  (when debug?
    (enable-console-print!)
    (println "dev mode"))
  (mount-root))
