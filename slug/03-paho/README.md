# 03-paho

Minimal MQTT client in pure JavaScript, example from the [Eclipse
Paho](https://www.eclipse.org/paho/clients/js/) website.

This expects a Mosquitto broker to be running locally, with websockets and http
enabled.  
In the "Extra listeners" sections of `/usr/local/etc/mosquitto/mosquitto.conf`:

    listener 9000
    protocol websockets
    http_dir .../jet/slug/03-paho     <=  adjust as needed

The above is for Mosquitto 1.4.14, which was installed on macOS using:

    brew install mosquitto
    brew services start mosquitto

Sample console output in a web browser pointed at <http://localhost:9000/>:

    onConnect
    onMessageArrived:Hello
