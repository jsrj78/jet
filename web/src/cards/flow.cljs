(ns cards.flow
  (:require [cljs.test :refer-macros [is testing]]
            [devcards.core :refer-macros [defcard-rg deftest]]
            [flow.core :as f]
            [flow.gadgets]))

(deftest core-gadgets-and-circuits

  (testing "print exists"
    (let [f-print (f/make-gadget :print)] 
      (is f-print)
      (is (= "[1 2 3]\n"
            (with-out-str (f/feed f-print 0 [1 2 3]))))))

  (testing "print scalar"
    (let [f-print (f/make-gadget :print)] 
      (is f-print)
      (is (= "333\n"
            (with-out-str (f/feed f-print 0 333))))))

  (testing "print with label"
    (let [f-print (f/make-gadget :print :hello)] 
      (is (= "[:hello 1 2 3]\n"
            (with-out-str (f/feed f-print 0 [1 2 3]))))))

  (testing "inlet and print"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :print))
                (f/add-wire [0 0 1 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (f/feed c 0 [1 2 3]))))))

  (testing "trivial circuits"
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

  (testing "nested circuit"
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

  (testing "swap gadget"
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

  (testing "bare send gadget"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :s :blah))
                (f/add-wire [0 0 1 0]))]
      (f/on c :blah prn)
      (is (= "[1 2 3]\n"
            (with-out-str (f/feed c 0 [1 2 3]))))))

  (testing "send and receive"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :s :blah))
                (f/add (f/make-gadget :r :blah))
                (f/add (f/make-gadget :print))
                (f/add-wire [0 0 1 0])
                (f/add-wire [2 0 3 0]))]
      (is (= "[1 2 3]\n"
            (with-out-str (f/feed c 0 [1 2 3]))))))

  (testing "smooth gadget"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :smooth 3))
                (f/add (f/make-gadget :print))
                (f/add-wire [0 0 1 0])
                (f/add-wire [1 0 2 0]))]
      (is (= "[0]\n[25]\n[43]\n[57]\n[67]\n[75]\n[81]\n[85]\n[88]\n[91]\n[93]\n"
            (with-out-str (doseq [x (cons 0 (repeat 10 100))]
                            (f/feed c 0 x)))))))

  (testing "change gadget"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :change))
                (f/add (f/make-gadget :print))
                (f/add-wire [0 0 1 0])
                (f/add-wire [1 0 2 0]))]
      (is (= "[0]\n[1]\n[2]\n[3]\n[0]\n"
            (with-out-str (doseq [x [0 1 1 2 2 3 0]]
                            (f/feed c 0 x)))))))

  (testing "moses gadget"
    (let [c (-> (f/make-circuit)
                (f/add (f/make-gadget :inlet))
                (f/add (f/make-gadget :moses 5))
                (f/add (f/make-gadget :print :a))
                (f/add (f/make-gadget :print :b))
                (f/add-wire [0 0 1 0])
                (f/add-wire [1 0 2 0])
                (f/add-wire [1 1 3 0]))]
      (is (= "[:a 4]\n[:b 5]\n[:b 6]\n[:b 5]\n[:a 4]\n"
            (with-out-str (doseq [x [4 5 6 5 4]]
                            (f/feed c 0 x))))))))
