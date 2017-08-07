(ns flow.core
  (:require [clojure.string :as s]))

(defonce registry (atom {})) 

; TODO use macro to inject the "fn" and use a symbol iso a keyword
(defn defgadget [key fun]
  (swap! registry assoc key fun))

(defn new-gadget []
  {;:state (atom {})
   :inlets []
   :outlets []
   :on-added identity})

(defn add-inlet [gob f]
  (update gob :inlets conj f))

(defn add-outlets [gob num]
  (update gob :outlets into (repeat num [])))

(defn feed [gob inlet msg]
  ((get (:inlets gob) inlet) msg))

(defn emit [gob outlet msg]
  (.log js/console "emit:" gob (get-in gob [:outlets outlet]))
  (doseq [[dst out] (get-in gob [:outlets outlet])]
    (feed dst out msg)))

(defn emitter [gob outlet]
  (fn [msg]
    (emit gob outlet msg)))

(defgadget :print
  (fn [label]
    (-> (new-gadget)
        (add-inlet (fn [msg]
                    (let [args (if label (cons label msg) msg)] 
                      (apply pr args)))))))

(defgadget :pass
  (fn []
    (let [gob (new-gadget)]
      (-> gob
          (add-outlets 1)
          (add-inlet (emitter gob 0))))))

(defgadget :inlet
  (fn []
    (let [gob (new-gadget)]
      (-> gob
          (add-outlets 1)
          (assoc :on-added #(add-inlet % (emitter gob 0)))))))

(defgadget :outlet
  (fn []
    (let [gob (new-gadget)]
      (-> gob
          (assoc :on-added #(let [n (count (:outlets %))
                                  nobj (add-inlet gob (emitter % n))]
                              (-> %
                                  ;;; FIXME
                                  (add-outlets 1))))))))

(defn lookup-gadget [key & args]
  (let [f (key @registry)] 
    (if f
      (apply f args)
      (.log js/console "no such gadget:" key))))

(defn new-circuit []
  (assoc (new-gadget) :g [] :w []))

(defn add [cob gob]
  (.log js/console "oa:" gob)
  (-> cob
      (update :g conj gob) 
      ((:on-added gob))))

(defn add-wire [cob [srcg srco dstg dsti :as wire]]
  (let [src (get-in cob [:g srcg]) 
        dst (get-in cob [:g dstg])]
    (-> cob
        (update :w conj [wire])
        (update-in [:g srcg :outlets] conj [dst dsti]))))
