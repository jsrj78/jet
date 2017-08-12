# JET/Flow

* implemented in Forth (for
  [Mecrisp-Stellaris](https://github.com/jeelabs/mecrisp-stellaris))
* just the bare essentials so far: message flow and nested circuits
* `gadgets.fs` contains some example gadget & handler defininitions

### Sample code

Using [folie](http://folie.jeelabs.org), this is the output from running
`dev.fs` on a HyTiny F103 board:

```
!s dev.fs
1> dev.fs 11: ( pass: ) 2000040C 20000424  ok.
1> dev.fs 12: ( 11 ) 11  ok.
1> dev.fs 13: ( 22 ) 22  ok.
1> dev.fs 21: ( two: ) 20000450 20000468  ok.
1> dev.fs 22: ( 11 456 22 456 ) 11 456 22 456  ok.
1> dev.fs 23: ( 11 789 22 789 ) 11 789 22 789  ok.
1> dev.fs 34: ( nest: ) 200004B0 200004C8  ok.
1> dev.fs 35: ( 111 222 ) 111 222  ok.
1> dev.fs 44: ( swap: ) 20000504 20000520  ok.
1> dev.fs 45: ( 22 333 11 123 ) 22 333 11 123  ok.
1> dev.fs 53: ( change: ) 2000054C 20000564  ok.
1> dev.fs 54: ( 0 ) 0  ok.
1> dev.fs 55: ( 1 ) 1  ok.
1> dev.fs 56: ( 2 ) 2  ok.
1> dev.fs 58: ( 3 ) 3  ok.
1> dev.fs 60: ( 0 ) 0  ok.
1> dev.fs 69: ( moses: ) 200005A0 200005BC  ok.
1> dev.fs 70: ( 11 4 ) 11 4  ok.
1> dev.fs 71: ( 22 5 ) 22 5  ok.
1> dev.fs 72: ( 22 6 ) 22 6  ok.
1> dev.fs 74: .s Stack: [0 ]  TOS: 42  *>
 ok.
```

Comments were added to generate output of the form "(expected) actual".
