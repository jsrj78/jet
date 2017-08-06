(defproject app "0.1.0-SNAPSHOT"
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [org.clojure/clojurescript "1.9.854"]
                 [reagent "0.7.0"]
                 [re-frame "0.9.4"]]

  :plugins [[lein-cljsbuild "1.1.7"]]

  :min-lein-version "2.7.1"

  :clean-targets ^{:protect false} ["public/js" "target"]

  :figwheel {:css-dirs ["public"]}

  :profiles {:dev {:dependencies [[binaryage/devtools "0.9.4"]]
                   :plugins      [[lein-figwheel "0.5.12"]]}}

  :resource-paths ["."]

  :cljsbuild
  {:builds
   [{:id           "dev"
     :source-paths ["src"]
     :figwheel     {:on-jsload "app.core/mount-root"}

     :compiler     {:main                 app.core
                    :output-to            "public/js/app.js"
                    :output-dir           "public/js/out"
                    :asset-path           "js/out"
                    :source-map-timestamp true
                    :preloads             [devtools.preload]
                    :external-config      {:devtools/config
                                           {:features-to-install :all}}}}

    {:id           "min"
     :source-paths ["src"]
     :compiler     {:main            app.core
                    :output-to       "public/js/app.js"
                    :optimizations   :advanced
                    :closure-defines {goog.DEBUG false}
                    :pretty-print    false}}]})
