(ns flow.core
  (:require [clojure.string :as s]))

(defonce registry (atom {})) 

; TODO use macro to inject the "fn" and use a symbol iso a keyword
(defn defgadget [key fun]
  (swap! registry assoc key fun))

(defn init-gadget []
  {:inlets []
   :outlets (atom [])
   :on-add identity})

(defn add-inlet [gob inlet-fn]
  (update gob :inlets conj inlet-fn))

(defn add-outlets [gob num]
  (swap! (:outlets gob) into (repeat num []))
  (- (count @(:outlets gob)) num))

(defn feed [gob inlet msg]
  ((get (:inlets gob) inlet) msg))

(defn emit [gob outlet msg]
  (doseq [[dst-gob dst-in] (get @(:outlets gob) outlet)]
    (feed dst-gob dst-in msg)))

(defn emitter [gob outlet]
  (fn [msg]
    (emit gob outlet msg)))

(defn make-gadget [k & args]
  (let [g-fn (k @registry)] 
    (if g-fn
      (apply g-fn args)
      (.error js/console "no such gadget:" key))))

(defn make-circuit []
  (-> (init-gadget)
      (assoc :gadgets []
             :wires []
             :notifiers (atom {}))))

(defn add [cob gob]
  (-> ((:on-add gob) cob)
      (update :gadgets conj gob))) 

(defn add-wire [cob [src-id src-out dst-id dst-in :as wire]]
  (let [src-gob (get-in cob [:gadgets src-id])
        dst-gob (get-in cob [:gadgets dst-id])]
    (swap! (:outlets src-gob) update src-out conj [dst-gob dst-in])
    (update cob :wires conj wire))) 

(defn on [cob topic on-fn]
  ; FIXME assoc should append to a vec, current code supports one fn per topic
  (swap! (:notifiers cob) assoc topic [on-fn]))

(defn notify [cob topic msg]
  (doseq [on-fn (get @(:notifiers cob) topic)] 
    (on-fn msg)))
