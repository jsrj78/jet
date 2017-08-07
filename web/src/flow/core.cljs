(ns flow.core
  (:require [clojure.string :as s]))

(defonce registry (atom {})) 

; TODO use macro to inject the "fn" and use a symbol iso a keyword
(defn defgadget [key fun]
  (swap! registry assoc key fun))

(defn new-gadget [num-outs]
  {:state (atom {})
   :inlets []
   :outlets (vec (repeat num-outs []))})

(defn emit [gadget outlet msg])

(defgadget :print
  (fn [label]
    (let [obj (new-gadget 0)
          ins [(fn [msg]
                (let [args (if label (cons label msg) msg)] 
                  (apply pr args)))]]
      (assoc obj :inlets ins))))

(defgadget :pass
  (fn []
    (let [obj (new-gadget 1)
          ins [(fn [msg]
                (emit obj 0 msg))]]
      (assoc obj :inlets ins))))

(defn lookup-gadget [key & args]
  (let [f (key @registry)] 
    (if f
      (apply f args)
      (.log js/console "no such gadget:" key))))

(defn feed [gadget inlet msg]
  (((:inlets gadget) inlet) msg))

(defn new-circuit [gadgets wires]
  (let [c (atom {:g gadgets :w wires})]
    c))
