(ns app.events
  (:require [re-frame.core :as rf]
            [app.db :as db]))

(rf/reg-event-db
  :initialize-db
  (fn  [_ _]
    db/default-db))

(rf/reg-event-db
  :move-gadget
  (fn [db [_ id dx dy]]
    (update-in db [:gadgets id]
                  (fn [[vhead vx vy & vtail]]
                    (into [vhead (+ vx dx) (+ vy dy)] vtail)))))

(rf/reg-event-db
  :select-gadget
  (fn [db [_ id]]
    (assoc db :selected-gadget id)))

(rf/reg-event-db
  :set-label-width
  (fn [db [_ id size]]
    (assoc-in db [:label-widths id] size)))
