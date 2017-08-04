(ns app.events
  (:require [re-frame.core :as rf]
            [app.db :as db]))

(rf/reg-event-db
 :initialize-db
 (fn  [_ _]
   db/default-db))

(rf/reg-event-db
 :change-name
 (fn  [db [_ new-name]]
   (assoc db :name new-name)))
