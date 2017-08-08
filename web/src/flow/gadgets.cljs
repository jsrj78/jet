(ns flow.gadgets
  (:require [flow.core :as f]))

(f/defgadget :print
  (fn [label]
    (let [gob (f/init-gadget)]
      (f/add-inlet gob (fn [msg]
                        (let [args (if label (cons label msg) msg)] 
                          (prn (vec args))))))))

(f/defgadget :pass
  (fn []
    (let [gob (f/init-gadget)]
      (f/add-outlets gob 1)
      (f/add-inlet gob (f/emitter gob 0)))))

(f/defgadget :inlet
  (fn []
    (let [gob (f/init-gadget)]
      (f/add-outlets gob 1)
      (assoc gob :on-add (fn [cob]
                           (f/add-inlet cob (f/emitter gob 0)))))))

(f/defgadget :outlet
  (fn []
    (let [gob (f/init-gadget)]
      (assoc gob :on-add (fn [cob]
                          (let [off (count (:gadgets cob))]
                            (->> (f/add-outlets cob 1)
                                 (f/emitter cob)
                                 (f/add-inlet gob)
                                 (assoc-in cob [:gadgets off]))))))))

(f/defgadget :swap
  (fn [args]
    (let [val (atom args)
          gob (f/init-gadget)]
      (f/add-outlets gob 2)
      (-> gob
        (f/add-inlet (fn [msg]
                      (f/emit gob 1 msg)
                      (f/emit gob 0 @val)))
        (f/add-inlet #(reset! val %))))))

(f/defgadget :s
  (fn [topic]
    (let [parent (atom nil)]
      (-> (f/init-gadget)
          (f/add-inlet #(f/notify @parent topic %))
          (assoc :on-add #(reset! parent %))))))

(f/defgadget :r
  (fn [topic]
    (let [gob (f/init-gadget)]
      (f/add-outlets gob 1)
      (assoc gob :on-add (fn [cob]
                           (f/on cob topic (f/emitter gob 0))
                           cob)))))
(f/defgadget :smooth
  (fn [arg]
    (let [hist  (atom 0)
          order (atom arg)
          gob   (f/init-gadget)]
      (f/add-outlets gob 1)
      (-> gob
        (f/add-inlet (fn [msg]
                       (let [[o] @order]
                        (reset! hist (int (/ (+ (* o @hist) msg) (inc o))))
                        (f/emit gob 0 [@hist]))))
        (f/add-inlet #(reset! order %))))))

(f/defgadget :change
  (fn []
    (let [last  (atom nil)
          gob   (f/init-gadget)]
      (f/add-outlets gob 1)
      (f/add-inlet gob (fn [msg]
                        (if (not= msg @last)
                          (do
                            (reset! last msg)
                            (f/emit gob 0 [msg]))))))))
