(ns cards.flow
  (:require-macros [devcards.core :refer [defcard-rg deftest]]
                   [cljs.test :refer [testing is]])
  (:require [flow.core :as f]))

(deftest print-exists
  (let [f-print (f/make-gadget :print)] 
    (is f-print)
    (is (= "1 2 3\n"
           (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest print-with-label
  (let [f-print (f/make-gadget :print "hello")] 
    (is (= "\"hello\" 1 2 3\n"
           (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest inlet-and-print
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0]))]
    (is (= "1 2 3\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))

(deftest trivial-circuits

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :pass))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0]))]
    (is (= "1 2 3\n"
           (with-out-str (f/feed c 0 [1 2 3])))))

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :print :a))
              (f/add (f/make-gadget :print :b))
              (f/add-wire [0 0 1 0])
              (f/add-wire [0 0 2 0]))]
    (is (= ":a 1 2 3\n:b 1 2 3\n"
           (with-out-str (f/feed c 0 [1 2 3])))))

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :pass))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0])
              (f/add-wire [0 0 2 0]))]
    (is (= "1 2 3\n1 2 3\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))

(deftest nested-circuit
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (-> (f/make-circuit)
                         (f/add (f/make-gadget :inlet))
                         (f/add (f/make-gadget :outlet))
                         (f/add-wire [0 0 1 0])))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0]))]
    (is (= "1 2 3\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))
