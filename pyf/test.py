from __future__ import print_function

import flow
import gadgets

print(">>> pass gadget")

c = flow.Circuit()
c.add('inlet')
c.add('pass')
c.add('print', 1)
c.add('print', 2)
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)
c.wire(1, 0, 3, 0)
c.wire(0, 0, 3, 0)

c.feed(0, 'bingo')

print(">>> swap gadget")

c = flow.Circuit()
c.add('inlet')
c.add('swap', [1, 2, 3])
c.add('print', 'a')
c.add('print', 'b')
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)
c.wire(1, 1, 3, 0)

for msg in [111, 222]:
    c.feed(0, msg)

print(">>> s and r gadgets")

c = flow.Circuit()
c.add('inlet')
c.add('s', 'blah')
c.add('r', 'blah')
c.add('print')
c.wire(0, 0, 1, 0)
c.wire(2, 0, 3, 0)

c.feed(0, [1, 2, 3])
c.feed(0, [4, 5, 6])

print(">>> smooth gadget")

c = flow.Circuit()
c.add('inlet')
c.add('smooth', 3)
c.add('print')
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)

for msg in [0] + 10*[100]:
    c.feed(0, msg)

print(">>> change gadget")

c = flow.Circuit()
c.add('inlet')
c.add('change')
c.add('print')
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)

for msg in [0, 1, 1, 2, 2, 3, 0]:
    c.feed(0, msg)

print(">>> moses gadget")

c = flow.Circuit()
c.add('inlet')
c.add('moses', 5)
c.add('print', 'a')
c.add('print', 'b')
c.wire(0, 0, 1, 0)
c.wire(1, 0, 2, 0)
c.wire(1, 1, 3, 0)

for msg in [4, 5, 6, 5, 4]:
    c.feed(0, msg)
