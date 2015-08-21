(defproject hub "0.2.0-SNAPSHOT"
  :description "JeeLabs Embello Toolkit Hub"
  :url "https://github.com/jeelabs/jet"
  :license {:name "The Unlicense"
            :url "http://unlicense.org"}

  :dependencies [[org.clojure/clojure "1.7.0"]
                 [org.clojure/clojurescript "1.7.107"]
                 [figwheel "0.3.5"]]

  :plugins [[lein-cljsbuild "1.0.6"]
            [lein-figwheel "0.3.5"]]

  :source-paths ["src"]

  :clean-targets ["out.dev"
                  "out.prod"
                  "server.js"]

  :cljsbuild {
    :builds [{:id "dev"
              :source-paths ["src" "src.dev"]
              :compiler {
                :output-to "out.dev/hub.js"
                :output-dir "out.dev"
                :target :nodejs
                :optimizations :none
                :source-map true}}
             {:id "prod"
              :source-paths ["src"]
              :compiler {
                :output-to "server.js"
                :output-dir "out.prod"
                :target :nodejs
                :optimizations :simple}}]})
