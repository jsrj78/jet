(ns app.events
  (:require [re-frame.core :as rf]
            [app.db :as db]))

(rf/reg-event-db
  :initialize-db
  (fn  [_ _]
    db/default-db))

(rf/reg-event-db
  :move-gadget
  (fn [db [_ oid dx dy]]
    (update-in db [:design oid]
                  (fn [[vhead vx vy & vtail :as v]]
                    (into [vhead (+ vx dx) (+ vy dy)] vtail)))))
