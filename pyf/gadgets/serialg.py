from __future__ import print_function

import flow
import threading
import serial

class SerialG(flow.Gadget):
    def __init__(self, dev):
        flow.Gadget.__init__(self, 1)
        try:
            import serial
        except Exception as e:
            print("PySerial package not found")
            raise
        self.ser = serial.Serial(dev, 115200)
        t = threading.Thread(target=self.reader)
        t.daemon = True
        t.start()

    def reader(self):
        while True:
            line = self.ser.readline()
            self.emit(0, line[:-1])

    def feed(self, inum, msg):
        self.ser.write(msg.encode('utf-8') + b'\r')

flow.registry['serial'] = SerialG
