(ns app.db)

(def default-db
  {:selected-gadget nil
   :gadgets
   [[:obj 50 40 :inlet]
    [:obj 50 90 :swap 555]
    [:obj 50 140 :print 1]
    [:obj 150 140 :print 12]]
   :wires
   [[0 0 1 0]
    [1 0 2 0]
    [1 1 3 0]]})
