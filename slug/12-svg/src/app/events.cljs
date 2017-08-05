(ns app.events
  (:require [re-frame.core :as rf]
            [app.db :as db]))

(rf/reg-event-db
  :initialize-db
  (fn  [_ _]
    db/default-db))

(rf/reg-event-db
  :move-obj
  (fn [db [_ idx dx dy]]
    (.log js/console "move:" idx dx dy)
    (update-in db [:design idx]
                  (fn [[vhead vx vy & vtail :as v]]
                    (into [vhead (+ vx dx) (+ vy dy)] vtail)))))
