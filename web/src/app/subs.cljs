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
 (fn [db [_ pos]]
   (nth (:gadgets db) pos)))

(rf/reg-sub
 :current-gadget
 (fn [db _]
   (:selected-gadget db)))
