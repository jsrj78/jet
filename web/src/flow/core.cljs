(ns flow.core
  (:require [clojure.string :as s]))

(defonce registry (atom {})) 

; TODO use macro to inject the "fn" and use a symbol iso a keyword
(defn defgadget [key fun]
  (swap! registry assoc key fun))

(defn init-gadget []
  (atom {:inlets []
         :outlets []
         :on-added identity}))

(defn add-inlet [gob inlet-fn]
  (swap! gob update :inlets conj inlet-fn)
  gob)

(defn add-outlets [gob num]
  (swap! gob update :outlets into (repeat num []))
  gob)

(defn feed [gob inlet msg]
  ((get (:inlets @gob) inlet) msg))

(defn emit [gob outlet msg]
  (doseq [[dst out] (get-in @gob [:outlets outlet])]
    (feed dst out msg)))

(defn emitter [gob outlet]
  (fn [msg]
    (emit gob outlet msg)))

(defgadget :print
  (fn [label]
    (let [gob (init-gadget)]
      (-> gob
          (add-inlet (fn [msg]
                      (let [args (if label (cons label msg) msg)] 
                        (apply pr args))))))))

(defgadget :pass
  (fn []
    (let [gob (init-gadget)]
      (add-outlets gob 1)
      (add-inlet gob (emitter gob 0)))))

(defgadget :inlet
  (fn []
    (let [gob (init-gadget)]
      (add-outlets gob 1)
      (swap! gob assoc :on-added (fn [cob]
                                  (add-inlet cob (emitter gob 0))))
      gob)))

(defgadget :outlet
  (fn []
    (let [gob (init-gadget)]
      (swap! gob assoc :on-added #(let [n (count (:outlets @%))
                                        nobj (add-inlet gob (emitter % n))]
                                    (-> %
                                        ;;; FIXME
                                        (add-outlets 1))))
      gob)))

(defn make-gadget [k & args]
  (let [g-fn (k @registry)] 
    (if g-fn
      (apply g-fn args)
      (.error js/console "no such gadget:" key))))

(defn make-circuit []
  (let [cob (init-gadget)]
    (swap! cob assoc :g [] :w [])
    cob))

(defn add [cob gob]
  (swap! cob update :g conj gob) 
  ((:on-added @gob) cob))

(defn add-wire [cob [srcg srco dstg dsti :as wire]]
  (let [src (get-in @cob [:g srcg])
        dst (get-in @cob [:g dstg])]
    (swap! cob update :w conj wire) 
    (swap! src update-in [:outlets srco] conj [dst dsti])
    cob))
