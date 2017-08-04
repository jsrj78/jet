(ns app.subs
  (:require [re-frame.core :as rf]))

(rf/reg-sub
 :name
 (fn [db _]
   (:name db)))

(rf/reg-sub
 :design
 (fn [db _]
   (:design db)))
