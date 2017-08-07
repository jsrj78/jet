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
  (update gob :outlets swap! into (repeat num []))
  gob)

(defn feed [gob inlet msg]
  ((get (:inlets gob) inlet) msg))

(defn emit [gob outlet msg]
  (doseq [[dst out] (get @(:outlets gob) outlet)]
    (feed dst out msg)))

(defn emitter [gob outlet]
  (fn [msg]
    (emit gob outlet msg)))

(defgadget :print
  (fn [label]
    (-> (init-gadget)
        (add-inlet (fn [msg]
                    (let [args (if label (cons label msg) msg)] 
                      (apply pr args)))))))

(defgadget :pass
  (fn []
    (let [gob (init-gadget)]
      (-> gob
          (add-outlets 1)
          (add-inlet (emitter gob 0))))))

(defgadget :inlet
  (fn []
    (let [gob (init-gadget)]
      (-> gob
          (add-outlets 1)
          (assoc :on-add #(add-inlet % (emitter gob 0)))))))

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

(defn add-wire [cob [srcg srco dstg dsti :as wire]]
  (let [src (get-in cob [:g srcg])
        dst (get-in cob [:g dstg])]
    (swap! (:outlets src) update srco conj [dst dsti])
    (update cob :w conj wire))) 
