;; unused code

(defn by-id [id]
  (.getElementById js/document id))

(defn obj-id-as-xy [id]
  (mapv js/parseInt (s/split id ",")))

(defn bounding-client-xy [evt]
  (let [rect (.getBoundingClientRect (.-target evt))]
    [(.-left rect) (.-top rect)]))
