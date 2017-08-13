(ns app.db)

(def default-db
  {:selected-gadget nil
   :gadgets
   [[50 40 :obj :inlet]
    [120 40 [] :bang]
    [50 90 :obj :swap 555]
    [160 90 [] :bang]
    [50 140 :obj :print 1]
    [150 140 :obj :print 2]]
   :wires
   [[0 0 2 0]
    [1 0 2 0]
    [1 0 3 0]
    [2 0 4 0]
    [2 1 5 0]]})
