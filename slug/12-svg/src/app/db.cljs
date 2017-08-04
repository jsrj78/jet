(ns app.db)

(def default-db
  {:name "SVG demo"
   :design
   [[:obj 75 40 :inlet]
    [:obj 75 91 :swap 123]
    [:obj 75 142 :print 1]
    [:obj 146 142 :print 2]
    [:connect 0 0 1 0]
    [:connect 1 0 2 0]
    [:connect 1 1 3 0]]})
