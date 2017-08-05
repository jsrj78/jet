(ns app.subs
  (:require [re-frame.core :as rf]))

(rf/reg-sub
 :design
 (fn [db _]
   (:design db)))

(rf/reg-sub
 :gadget-num
 (fn [db [_ pos]]
   (nth (:design db) pos)))
