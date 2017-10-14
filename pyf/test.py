from __future__ import print_function

import flow, gadgets
import time

print(">>> metro gadget, 3 fast + 2 slow ticks")

c = flow.Circuit()
c.add('inlet')
c.add('inlet')
c.add('metro', 300)
c.add('print', 'timer:')
c.wire(0, 0, 2, 0)
c.wire(1, 0, 2, 1)
c.wire(2, 0, 3, 0)

time.sleep(1)
c.feed(0, 900)   # slow down the timer
time.sleep(2)
c.feed(0, None)  # cancel the timer
