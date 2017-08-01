; lein new reagent-figwheel slug +devcards
; https://github.com/gadfly361/reagent-figwheel

(defproject slug "0.1.0-SNAPSHOT"
  :dependencies [[org.clojure/clojure "1.8.0" :scope "provided"]
                 [org.clojure/clojurescript "1.9.671" :scope "provided"]
                 [reagent "0.7.0"]
                 [devcards "0.2.3" :exclusions [cljsjs/react]]]

  :min-lein-version "2.7.1"

  :plugins [[lein-cljsbuild "1.1.5"]]

  :clean-targets ^{:protect false} ["resources/public/js" "target"]

  :figwheel {:css-dirs ["resources/public/css"]}

  :profiles
  {:dev
   {:dependencies []
    :plugins      [[lein-figwheel "0.5.11"]]}}

  :cljsbuild
  {:builds
   [{:id           "dev"
     :source-paths ["src"]
     :figwheel     {:on-jsload "slug.core/reload"}
     :compiler     {:main                 slug.core
                    :optimizations        :none
                    :output-to            "resources/public/js/app.js"
                    :output-dir           "resources/public/js/dev"
                    :asset-path           "js/dev"
                    :source-map-timestamp true}}

    {:id           "devcards"
     :source-paths ["src"]
     :figwheel     {:devcards true}
     :compiler     {:main                 "slug.core-card"
                    :optimizations        :none
                    :output-to            "resources/public/js/devcards.js"
                    :output-dir           "resources/public/js/devcards"
                    :asset-path           "js/devcards"
                    :source-map-timestamp true}}

    {:id           "hostedcards"
     :source-paths ["src"]
     :compiler     {:main          "slug.core-card"
                    :optimizations :advanced
                    :devcards      true
                    :output-to     "resources/public/js/devcards.js"
                    :output-dir    "resources/public/js/hostedcards"}}

    {:id           "min"
     :source-paths ["src"]
     :compiler     {:main            slug.core
                    :optimizations   :advanced
                    :output-to       "resources/public/js/app.js"
                    :output-dir      "resources/public/js/min"
                    :elide-asserts   true
                    :closure-defines {goog.DEBUG false}
                    :pretty-print    false}}]})
