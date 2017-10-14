This is an exploration of how to implement a flow engine in Python.

Needs the pyserial and paho-mqtt packages.

Quick check that the base gadgets work:

    python test_base.py

Running as an MQTT service called "pyf-demo":

    python connect.py

Sample test, with a serial connection to an attached STM32 w/ Mecrisp Forth:

    python test_serial.py

Sample output:

    $ python test_serial.py
    reply: s/pyf-demo/test/out/0 "1 2 + . 3  ok."
    reply: s/pyf-demo/test/out/0 "11 22 + . 33  ok."

This is a dump of all MQTT messages exchanged in the above test:

    s/pyf-demo ["create", "test"]
    s/pyf-demo/test [["inlet"], ["serial", "/dev/cu.usbmodem34208131"], ["outlet"], [0, 0, 1, 0], [1, 0, 2, 0]]
    s/pyf-demo/test/in/0 "1 2 + ."
    s/pyf-demo/test/out/0 "1 2 + . 3  ok."
    s/pyf-demo/test/in/0 "11 22 + ."
    s/pyf-demo/test/out/0 "11 22 + . 33  ok."

Note the hardwired serial port in this example.
