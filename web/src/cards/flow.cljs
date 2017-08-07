(ns cards.flow
  (:require-macros [devcards.core :refer [defcard-rg deftest]]
                   [cljs.test :refer [testing is]])
  (:require [flow.core :as f]
            [flow.gadgets]))

(deftest print-exists
  (let [f-print (f/make-gadget :print)] 
    (is f-print)
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest print-with-label
  (let [f-print (f/make-gadget :print :hello)] 
    (is (= "[:hello 1 2 3]\n"
           (with-out-str (f/feed f-print 0 [1 2 3]))))))

(deftest inlet-and-print
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0]))]
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))

(deftest trivial-circuits

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :pass))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0]))]
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3])))))

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :print :a))
              (f/add (f/make-gadget :print :b))
              (f/add-wire [0 0 1 0])
              (f/add-wire [0 0 2 0]))]
    (is (= "[:a 1 2 3]\n[:b 1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3])))))

  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :pass))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0])
              (f/add-wire [0 0 2 0]))]
    (is (= "[1 2 3]\n[1 2 3]\n"
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
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))

(deftest swap-gadget
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :swap [1 2 3]))
              (f/add (f/make-gadget :print :a))
              (f/add (f/make-gadget :print :b))
              (f/add-wire [0 0 1 0])
              (f/add-wire [1 0 2 0])
              (f/add-wire [1 1 3 0]))]
    (is (= "[:b 111]\n[:a 1 2 3]\n"
           (with-out-str (f/feed c 0 [111]))))
    (is (= "[:b 222]\n[:a 1 2 3]\n"
           (with-out-str (f/feed c 0 [222]))))))

(deftest bare-send-gadget
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :s :blah))
              (f/add-wire [0 0 1 0]))]
    (f/on c :blah prn)
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))

(deftest send-and-receive
  (let [c (-> (f/make-circuit)
              (f/add (f/make-gadget :inlet))
              (f/add (f/make-gadget :s :blah))
              (f/add (f/make-gadget :r :blah))
              (f/add (f/make-gadget :print))
              (f/add-wire [0 0 1 0])
              (f/add-wire [2 0 3 0]))]
    (is (= "[1 2 3]\n"
           (with-out-str (f/feed c 0 [1 2 3]))))))
