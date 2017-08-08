(ns cards.smoke
  (:require [cljs.test :refer-macros [is testing]]
            [devcards.core :refer-macros [defcard-rg deftest]]))

(deftest smoke-test
  (testing "First tests"
    (is (= 0 0))
    (is (= 1 1)))
  (testing "next tests"
    #_(is (= 1 2))
    (is (= 2 2))))

(defn twice [n]
  (+ n n))

(deftest twice-test
  (is (= (twice 12) 24)))
