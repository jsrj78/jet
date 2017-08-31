(ns flow.gadgets
  (:require [flow.core :as flow]))

(flow/defgadget :print
  (fn [label]
    (let [gob (flow/init-gadget)]
      (flow/add-inlet gob #(prn (or label :print) %)))))

(flow/defgadget :pass
  (fn []
    (let [gob (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (flow/add-inlet gob (flow/emitter gob 0)))))

(flow/defgadget :inlet
  (fn []
    (let [gob (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (assoc gob :on-add (fn [cob]
                           (flow/add-inlet cob (flow/emitter gob 0)))))))

(flow/defgadget :outlet
  (fn []
    (let [gob (flow/init-gadget)]
      (assoc gob :on-add (fn [cob]
                          (let [off (count (:gadgets cob))]
                            (->> (flow/add-outlets cob 1)
                                 (flow/emitter cob)
                                 (flow/add-inlet gob)
                                 (assoc-in cob [:gadgets off]))))))))

(flow/defgadget :swap
  (fn [args]
    (let [*val (atom args)
          gob (flow/init-gadget)]
      (flow/add-outlets gob 2)
      (-> gob
          (flow/add-inlet (fn [msg]
                            (flow/emit gob 1 msg)
                            (flow/emit gob 0 @*val)))
          (flow/add-inlet #(reset! *val %))))))

(flow/defgadget :s
  (fn [topic]
    (let [*parent (atom nil)]
      (-> (flow/init-gadget)
          (flow/add-inlet #(flow/notify @*parent topic %))
          (assoc :on-add #(reset! *parent %))))))

(flow/defgadget :r
  (fn [topic]
    (let [gob (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (assoc gob :on-add (fn [cob]
                           (flow/on cob topic (flow/emitter gob 0))
                           cob)))))

(flow/defgadget :smooth
  (fn [arg]
    (let [*hist  (atom 0)
          *order (atom arg)
          gob    (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (-> gob
          (flow/add-inlet (fn [msg]
                            (let [o @*order
                                  h @*hist]
                              (reset! *hist (int (/ (+ (* o h) msg) (inc o))))
                              (flow/emit gob 0 [@*hist]))))
          (flow/add-inlet #(reset! *order %))))))

(flow/defgadget :change
  (fn []
    (let [*last (atom nil)
          gob   (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (flow/add-inlet gob (fn [msg]
                            (if (not= msg @*last)
                              (do
                                (reset! *last msg)
                                (flow/emit gob 0 [msg]))))))))

(flow/defgadget :moses
  (fn [arg]
    (let [*split (atom arg)
          gob    (flow/init-gadget)]
      (flow/add-outlets gob 1)
      (-> gob
          (flow/add-inlet (fn [msg]
                            (if (< msg @*split)
                              (flow/emit gob 0 [msg])
                              (flow/emit gob 1 [msg]))))
          (flow/add-inlet #(reset! *split %))))))
