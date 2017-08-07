(ns cards.flow
  (:require-macros [devcards.core :refer [defcard-rg deftest]]
                   [cljs.test :refer [testing is]])
  (:require [flow.core :as f]))

(deftest print-exists
  (let [f-print (f/lookup-gadget :print)] 
    (is f-print)
    (is (= "1 2 3" (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest print-with-label
  (let [f-print (f/lookup-gadget :print "hello")] 
    (is (= "\"hello\" 1 2 3" (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest pass-and-print
  (let [f-print (f/lookup-gadget :print) 
        circuit (-> (f/new-circuit)
                    (f/add (f/lookup-gadget :inlet))
                    (f/add f-print)
                    (f/add-wire [0 0 1 0]))]
    (.log js/console circuit)
    (is (= "1 2 3" (with-out-str (f/feed circuit 0 [1 2 3]))))))
