(ns hub.dev
  (:require [figwheel.client]
            [hub.core]))

(defn -main []
  (figwheel.client/start)
  (hub.core/-main))

(set! *main-cli-fn* -main)
