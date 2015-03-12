(ns blink3.test-runner
  (:require
   [cljs.test :refer-macros [run-tests]]
   [blink3.core-test]))

(enable-console-print!)

(defn runner []
  (if (cljs.test/successful?
       (run-tests
        'blink3.core-test))
    0
    1))
