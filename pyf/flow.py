#from __future__ import print_function

registry = {}

class Gadget:
    "A gadget responds to messages fed to it and can emit more messages."

    def __init__(self, onum=0):
        self.outlets = [[] for _ in range(onum)]

    def feed(self, inum, msg):
        raise NotImplementedError

    def emit(self, onum, msg):
        for (g,i) in self.outlets[onum]:
            g.feed(i, msg)

class Circuit(Gadget):
    "A circuit is a gadget built from other gadgets and interconnecting wires."

    def __init__(self):
        self.inlets = []
        self.outlets = []
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
