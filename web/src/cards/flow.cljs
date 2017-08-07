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
  (let [f-pass  (f/lookup-gadget :pass)
        f-print (f/lookup-gadget :print) 
        circuit (f/new-circuit [f-pass f-print]
                               [[0 0 1 0]])]
    (is (and f-pass f-print))
    (is (= "\"hello\" 1 2 3" (with-out-str (f/feed f-print 0 [1 2 3]))))))
