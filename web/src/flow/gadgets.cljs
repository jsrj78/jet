(ns flow.gadgets
  (:require [flow.core :as f]))

(f/defgadget :print
  (fn [label]
    (-> (f/init-gadget)
        (f/add-inlet (fn [msg]
                      (let [args (if label (cons label msg) msg)] 
                        (apply prn args)))))))

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
