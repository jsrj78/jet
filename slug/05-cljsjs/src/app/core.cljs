(ns app.core
    (:require [cljsjs.paho]))

(enable-console-print!)

;; define your app data so that it doesn't get over-written on reload
;;(defonce app-state (atom {:text "Hello world!"}))

(declare client send-message)

(def topic "/cljsjs/pahotest")

(println "Hello console!")

(defn on-connect []
  (println "Connected")
  (.subscribe client topic #js {:qos 0})
  (println "Subscribed")
  (send-message "Hello MQTT!" topic 0)
  (println "Sent message."))

(defn send-message [payload destination qos]
  (let [msg (Paho.MQTT.Message. payload)]
    (set! (.-destinationName msg) destination)
    (set! (.-qos msg) qos)
    (.send client msg)))

(defn connect []
  (let [mqtt (Paho.MQTT.Client. "test.mosquitto.org" 8080 "")
        connectOptions (js/Object.)]
       (set! (.-onConnectionLost mqtt) #(println %1 %2))
       (set! (.-onMessageArrived mqtt)
             #(println (str "Topic: " (.-destinationName %)
                            " Payload: " (.-payloadString %))))
       (set! (.-onSuccess connectOptions) on-connect)
       (set! (.-onFailure connectOptions ) #(println "Failure Connect: " %3))
       (.connect mqtt connectOptions)
       mqtt))

(def client (connect))

(defn on-js-reload []
  ;; optionally touch your app-state to force rerendering depending on
  ;; your application
  ;; (swap! app-state update-in [:__figwheel_counter] inc)
)
