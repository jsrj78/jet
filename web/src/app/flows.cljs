(ns app.flows 
  (:require [re-frame.core :as rf]
            [flow.core :as flow]
            [flow.gadgets]))

(defn map-gadget-to-engine [id obj]
  (if (= (nth obj 2) :obj)
    (subvec obj 3)
    (case (nth obj 3)
      :msg [:message (nth obj 4)]
           [:pass id])))

(defn create-engine [gadgets wires]
  (let [g  (map-indexed map-gadget-to-engine gadgets)
        c0 (flow/make-circuit)
        c1 (reduce #(flow/add %1 (apply flow/make-gadget %2)) c0 g)]
    (reduce flow/add-wire c1 wires)))

(def init-engine
  (rf/->interceptor
    :id    :init-engine
    :after (fn [{:keys [effects] :as context}]
             (let [{:keys [db]} effects
                   {:keys [gadgets wires]} db
                   circuit (create-engine gadgets wires)] 
               (assoc-in context [:effects :db :engine] circuit))))) 

(rf/reg-fx
  :send-bang
  (fn [[cob id]]
    (flow/feed (get-in cob [:gadgets id]) 0 [])))
