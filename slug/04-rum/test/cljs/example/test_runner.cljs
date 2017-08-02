(ns example.test-runner
  (:require
   [doo.runner :refer-macros [doo-tests]]
   [example.core-test]
   [example.common-test]))

(enable-console-print!)

(doo-tests 'example.core-test
           'example.common-test)
