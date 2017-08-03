# 10-testme

See <https://github.com/vallard/Reagent-Example>. Launch as `lein figwheel`.

For production, run `lein uberjar`, then launch using Java as back-end:

    java $JVM_OPTS -cp target/testme.jar clojure.main -m testme.server

This demo carries lots of extra luggage around, the generated jar file is 22 MB.
