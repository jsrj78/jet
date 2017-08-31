(ns cards.flow
  (:require [cljs.test :refer-macros [is testing]]
            [devcards.core :refer-macros [defcard-rg deftest]]
            [flow.core :as flow]
            [flow.gadgets]))

(deftest core-gadgets-and-circuits

  (testing "print exists"
    (let [f-print (flow/make-gadget :print)] 
      (is f-print)
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed f-print 0 [1 2 3]))))))

  (testing "print scalar"
    (let [f-print (flow/make-gadget :print)] 
      (is f-print)
      (is (= "333\n"
            (with-out-str (flow/feed f-print 0 333))))))

  (testing "print with label"
    (let [f-print (flow/make-gadget :print :hello)] 
      (is (= "[:hello 1 2 3]\n"
            (with-out-str (flow/feed f-print 0 [1 2 3]))))))

  (testing "inlet and print"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3]))))))

  (testing "trivial circuits"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :pass))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3])))))
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :print :a))
                (flow/add (flow/make-gadget :print :b))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [0 0 2 0]))]
      (is (= "[:a 1 2 3]\n[:b 1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3])))))
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :pass))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0])
                (flow/add-wire [0 0 2 0]))]
      (is (= "[1 2 3]\n[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3]))))))

  (testing "nested circuit"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (-> (flow/make-circuit)
                              (flow/add (flow/make-gadget :inlet))
                              (flow/add (flow/make-gadget :outlet))
                              (flow/add-wire [0 0 1 0])))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3]))))))

  (testing "swap gadget"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :swap [1 2 3]))
                (flow/add (flow/make-gadget :print :a))
                (flow/add (flow/make-gadget :print :b))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0])
                (flow/add-wire [1 1 3 0]))]
      (is (= "[:b 111]\n[:a 1 2 3]\n"
            (with-out-str (flow/feed c 0 [111]))))
      (is (= "[:b 222]\n[:a 1 2 3]\n"
            (with-out-str (flow/feed c 0 [222]))))))

  (testing "bare send gadget"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :s :blah))
                (flow/add-wire [0 0 1 0]))]
      (flow/on c :blah prn)
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3]))))))

  (testing "send and receive"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :s :blah))
                (flow/add (flow/make-gadget :r :blah))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [2 0 3 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (flow/feed c 0 [1 2 3]))))))

  (testing "smooth gadget"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :smooth 3))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0]))]
      (is (= "[0]\n[25]\n[43]\n[57]\n[67]\n[75]\n[81]\n[85]\n[88]\n[91]\n[93]\n"
            (with-out-str (doseq [x (cons 0 (repeat 10 100))]
                            (flow/feed c 0 x)))))))

  (testing "change gadget"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :change))
                (flow/add (flow/make-gadget :print))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0]))]
      (is (= "[0]\n[1]\n[2]\n[3]\n[0]\n"
            (with-out-str (doseq [x [0 1 1 2 2 3 0]]
                            (flow/feed c 0 x)))))))

  (testing "moses gadget"
    (let [c (-> (flow/make-circuit)
                (flow/add (flow/make-gadget :inlet))
                (flow/add (flow/make-gadget :moses 5))
                (flow/add (flow/make-gadget :print :a))
                (flow/add (flow/make-gadget :print :b))
                (flow/add-wire [0 0 1 0])
                (flow/add-wire [1 0 2 0])
                (flow/add-wire [1 1 3 0]))]
      (is (= "[:a 4]\n[:b 5]\n[:b 6]\n[:b 5]\n[:a 4]\n"
            (with-out-str (doseq [x [4 5 6 5 4]]
                            (flow/feed c 0 x))))))))
