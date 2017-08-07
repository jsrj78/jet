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
  (dec (count @(:outlets gob))))

(defn feed [gob inlet msg]
  ((get (:inlets gob) inlet) msg))

(defn emit [gob outlet msg]
  (doseq [[dst-gob dst-in] (get @(:outlets gob) outlet)]
    (feed dst-gob dst-in msg)))

(defn emitter [gob outlet]
  (fn [msg]
    (emit gob outlet msg)))

(defgadget :print
  (fn [label]
    (-> (init-gadget)
        (add-inlet (fn [msg]
                    (let [args (if label (cons label msg) msg)] 
                      (apply prn args)))))))

(defgadget :pass
  (fn []
    (let [gob (init-gadget)]
      (add-outlets gob 1)
      (add-inlet gob (emitter gob 0)))))

(defgadget :inlet
  (fn []
    (let [gob (init-gadget)]
      (add-outlets gob 1)
      (assoc gob :on-add #(add-inlet % (emitter gob 0))))))

(defgadget :outlet
  (fn []
    (let [gob (init-gadget)]
      (assoc gob :on-add (fn [cob]
                          (let [off (add-outlets cob 1)
                                gid (count (:g cob))
                                nob (add-inlet gob (emitter cob off))]
                            (assoc-in cob [:g gid] nob)))))))

(defn make-gadget [k & args]
  (let [g-fn (k @registry)] 
    (if g-fn
      (apply g-fn args)
      (.error js/console "no such gadget:" key))))

(defn make-circuit []
  (-> (init-gadget)
      (assoc :g [] :w [])))

(defn add [cob gob]
  (-> ((:on-add gob) cob)
      (update :g conj gob))) 

(defn add-wire [cob [src-id src-out dst-id dst-in :as wire]]
  (let [src-gob (get-in cob [:g src-id])
        dst-gob (get-in cob [:g dst-id])]
    (swap! (:outlets src-gob) update src-out conj [dst-gob dst-in])
    (update cob :w conj wire))) 
