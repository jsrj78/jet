(ns app.events
  (:require [re-frame.core :as rf]
            [app.db :as db]))

(rf/reg-event-db
  :initialize-db
  (fn  [_ _]
    db/default-db))

(rf/reg-event-db
  :move-obj
  (fn [db [_ x y dx dy]]
    (let [mover (fn [[vhead vx vy & vtail :as v]]
                  (if (and (= vx x) (= vy y))
                      (into [vhead (+ x dx) (+ y dy)] vtail)
                      v))]
      (assoc db :design (mapv mover (:design db))))))
