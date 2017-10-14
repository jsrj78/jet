from __future__ import print_function

import flow

class PrintG(flow.Gadget):
    def __init__(self, label=None):
        self.label = label

    def feed(self, inum, msg):
        if self.label:
            print(self.label, end=' ')
        print(msg)

flow.registry['print'] = PrintG

class PassG(flow.Gadget):
    def __init__(self):
        flow.Gadget.__init__(self, 1)

    def feed(self, inum, msg):
        self.emit(0, msg)

flow.registry['pass'] = PassG

class InletG(flow.Gadget):
    def __init__(self):
        flow.Gadget.__init__(self, 1)

    def onAdd(self, cob):
        cob.inlets.append(self)

flow.registry['inlet'] = InletG

class OutletG(flow.Gadget):
    def feed(self, inum, msg):
        self.parent.emit(self.onum, msg)

    def onAdd(self, cob):
        self.parent = cob
        self.onum = len(cob.outlets)
        cob.outlets.append([])

flow.registry['outlet'] = OutletG

class SwapG(flow.Gadget):
    def __init__(self, val=None):
        flow.Gadget.__init__(self, 2)
        self.val = val

    def feed(self, inum, msg):
        if inum == 0:
            self.emit(1, msg)
            self.emit(0, self.val)
        else:
            self.val = msg

flow.registry['swap'] = SwapG

class SendG(flow.Gadget):
    def __init__(self, topic):
        flow.Gadget.__init__(self)
        self.topic = topic

    def feed(self, inum, msg):
        self.parent.notify(self.topic, msg)

    def onAdd(self, cob):
        self.parent = cob

flow.registry['s'] = SendG

class ReceiveG(flow.Gadget):
    def __init__(self, topic):
        flow.Gadget.__init__(self, 1)
        self.topic = topic

    def onAdd(self, cob):
        cob.on(self.topic, lambda msg: self.emit(0, msg))

flow.registry['r'] = ReceiveG

class SmoothG(flow.Gadget):
    def __init__(self, order):
        flow.Gadget.__init__(self, 1)
        self.hist = 0
        self.order = order

    def feed(self, inum, msg):
        if inum == 0:
            o = self.order
            self.hist = ((o * self.hist) + msg) // (o + 1)
            self.emit(0, self.hist)
        else:
            self.order = msg

flow.registry['smooth'] = SmoothG

class ChangeG(flow.Gadget):
    def __init__(self):
        flow.Gadget.__init__(self, 1)
        self.last = None

    def feed(self, inum, msg):
        if msg != self.last:
            self.last = msg
            self.emit(0, msg)

flow.registry['change'] = ChangeG

class MosesG(flow.Gadget):
    def __init__(self, split):
        flow.Gadget.__init__(self, 2)
        self.split = split

    def feed(self, inum, msg):
        if inum == 0:
            if msg < self.split:
                self.emit(0, msg)
            else:
                self.emit(1, msg)
        else:
            self.split = msg

flow.registry['moses'] = MosesG
