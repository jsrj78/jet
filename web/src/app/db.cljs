(ns app.db)

#_(def default-db)
  {:label-widths {}
   :gadgets [[50 30 :obj :inlet]
             [120 40 [] :bang]
             [50 90 :obj :swap 555]
             [160 90 [] :bang]
             [200 75 [] :toggle]
             [70 140 [] :msg]
             [170 140 :obj :print 2]
             [50 180 :obj :print 1]]
   :wires [[0 0 2 0]
           [1 0 2 0]
           [1 0 3 0]
           [1 0 4 0]
           [2 0 5 0]
           [2 0 7 0]
           [2 1 6 0]]}

(def default-db
  {:label-widths {}
   :gadgets [[100 50 [] :bang]
             [100 100 :obj :print 123]]
   :wires [[0 0 1 0]]})
