(ns app.subs
  (:require [re-frame.core :as rf]))

(rf/reg-sub
  :gadgets
  (fn [db _]
    (:gadgets db)))

(rf/reg-sub
  :wires
  (fn [db _]
    (:wires db)))

(rf/reg-sub
  :gadget-num
  (fn [db [_ id]]
    (get-in db [:gadgets id])))

(rf/reg-sub
  :current-gadget
  (fn [db _]
    (:selected-gadget db)))

(rf/reg-sub
  :rect-width
  (fn [db [_ id]]
    (+ (get-in db [:label-widths id]) 11)))

(rf/reg-sub
  :num-iolets
  (fn [[_ id]]
    (rf/subscribe [:gadget-num id]))
  (fn [obj]
    (case (nth obj 3) ; TODO hard-coded for now
      :inlet  [0 1]
      :r      [0 1]
      :moses  [2 2]
      :swap   [2 2]
      :metro  [2 1]
      :print  [1 0]
      :outlet [1 0]
      :s      [1 0]
              [1 1])))

(defn spread-xy [n x y w]
  (let [s (/ (- w 5) (dec n))]
    (mapv #(list % y) (range (+ x 2.5) (+ x w) s))))

(rf/reg-sub
  :gadget-coords
  (fn [[_ id]]
    [(rf/subscribe [:gadget-num id])
     (rf/subscribe [:rect-width id])
     (rf/subscribe [:num-iolets id])])
  (fn [[[x y typ & cmd :as obj] w [ni no]] _]
     [[x y w 20]
      (spread-xy ni x y w)
      (spread-xy no x (+ y 20) w)]))
