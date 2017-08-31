(ns app.flows 
  (:require [re-frame.core :as rf]
            [flow.core :as flow]
            [flow.gadgets]))

(defn create-engine [gadgets wires]
  (-> (flow/make-circuit)
      (flow/add (flow/make-gadget :r 0))
      (flow/add (flow/make-gadget :print 123))
      (flow/add-wire [0 0 1 0])))

(def init-engine
  (rf/->interceptor
    :id    :init-engine
    :after (fn [{:keys [db] :as context}]
             (let [{:keys [gadgets wires]} db
                   circuit (create-engine gadgets wires)] 
               (assoc-in context [:effects :db :engine] circuit))))) 

(rf/reg-fx
  :send-bang
  (fn [[cob id]]
    (flow/notify cob 0 [])))
