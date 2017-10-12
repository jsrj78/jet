#!/usr/bin/env python
from __future__ import print_function

registry = {}

class Gadget:
    "A gadget responds to messages fed to it and can emit more messages"

    def __init__(self, onum=0):
        self.outlets = [[] for _ in range(onum)]

    def feed(self, inum, msg):
        raise NotImplementedError

    def emit(self, onum, msg):
        for (g,i) in self.outlets[onum]:
            g.feed(i, msg)

class Circuit(Gadget):
    def __init__(self):
        self.inlets = []
        self.gadgets = []
        self.wiring = []
        self.notifiers = {}

    def add(self, gname, *args):
        gob = registry[gname](*args)
        if hasattr(gob, 'onAdd'):
            gob.onAdd(self)
        self.gadgets.append(gob)

    def wire(self, sId, sOut, dId, dIn):
        sGob = self.gadgets[sId]
        dGob = self.gadgets[dId]
        sGob.outlets[sOut].append((dGob, dIn))
        w = (sId, sOut, dId, dIn)
        self.wiring.append(w)

    def on(self, topic, onFn):
        self.notifiers.setdefault(topic, []).append(onFn)

    def notify(self, topic, msg):
        for f in self.notifiers[topic]:
            f(msg)

    def feed(self, inum, msg):
        self.inlets[inum].emit(0, msg)

# Gadgets ---------------------------------------------------------------------

class PrintG(Gadget):
    def __init__(self, label=None):
        self.label = label

    def feed(self, inum, msg):
        if self.label:
            print(self.label, end=' ')
        print(msg)

registry['print'] = PrintG

class PassG(Gadget):
    def __init__(self):
        Gadget.__init__(self, 1)

    def feed(self, inum, msg):
        self.emit(0, msg)

registry['pass'] = PassG

class InletG(Gadget):
    def __init__(self):
        Gadget.__init__(self, 1)

    def onAdd(self, cob):
        cob.inlets.append(self)

registry['inlet'] = InletG

class OutletG(Gadget):
    def feed(self, inum, msg):
        self.parent.emit(self.onum, msg)

    def onAdd(self, cob):
        self.parent = cob
        self.onum = len(cob.outlets)
        cob.outlets.append([])

registry['outlet'] = OutletG

class SwapG(Gadget):
    def __init__(self, val=None):
        Gadget.__init__(self, 2)
        self.val = val

    def feed(self, inum, msg):
        if inum == 0:
            self.emit(1, msg)
            self.emit(0, val)
        else:
            val = msg

registry['swap'] = SwapG

class SendG(Gadget):
    def __init__(self, topic):
        self.topic = topic

    def feed(self, inum, msg):
        self.parent.notify(self.topic, msg)

    def onAdd(self, cob):
        self.parent = cob

registry['s'] = SendG

class ReceiveG(Gadget):
    def __init__(self, order):
        Gadget.__init__(self, 1)
        self.hist = 0
        self.order = order

    def feed(self, cob):
        cob.on(self.topic, lambda msg: self.emit(0, msg))

registry['r'] = ReceiveG

class SmoothG(Gadget):
    def __init__(self, topic):
        Gadget.__init__(self, 1)
        self.last = None

    def feed(self, inum, msg):
        if inum == 0:
            o = self.order
            self.hist = ((o * h) + msg) / (o + 1)
            self.emit(0, self.hist)
        else:
            self.order = msg

registry['smooth'] = SmoothG

class ChangeG(Gadget):
    def __init__(self, topic):
        Gadget.__init__(self, 1)
        self.last = None

    def feed(self, inum, msg):
        if msg <> self.last:
            self.last = msg
            self.emit(0, msg)

registry['change'] = ChangeG

class MosesG(Gadget):
    def __init__(self, split):
        Gadget.__init__(self, 2)
        self.split = split

    def feed(self, inum, msg):
        if inum == 0:
            if msg < self.split:
                self.emit(0, msg)
            else:
                self.emit(1, msg)
        else:
            self.split = msg

registry['moses'] = MosesG

# -----------------------------------------------------------------------------

g = PrintG('ha')
g.feed(0, 123)

c = Circuit()
c.add('inlet')
c.add('pass')
c.add('print', 1)
c.add('print', 2)
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)
c.wire(1, 0, 3, 0)
c.wire(0, 0, 3, 0)

c.feed(0, 'bingo')
