forgetram

include engine.fs

c:begin
  _ :inlet  0 1 0 wire  eow
  _ :pass  0 2 0 wire  eow
  _ i>m :print  eow
c:end  eow

( pass: ) hex. memp @ hex.
( 11 ) 11 i>m 0 cg @ feed
( 22 ) 22 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  0 2 0 wire  eow
  11 i>m :print  eow
  22 i>m :print  eow
c:end  eow

( two: ) hex. memp @ hex.
( 11 456 22 456 ) 456 i>m 0 cg @ feed
( 11 789 22 789 ) 789 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  eow
  c:begin
    _ :inlet  0 1 0 wire  eow
    _ :outlet
  c:end  1 2 0 wire  eow
  111 i>m :print  eow
c:end  eow

( nest: ) hex. memp @ hex.
( 111 222 ) 222 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  eow
  123 i>m :swap  0 2 0 wire  1 3 0 wire  eow
  11 i>m :print  eow
  22 i>m :print  eow
c:end  eow

( swap: ) hex. memp @ hex.
( 22 333 11 123 ) 333 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  eow
  _ :change  0 2 0 wire  eow
  _ :print  eow
c:end  eow

( change: ) hex. memp @ hex.
( 0 ) 0 i>m 0 cg @ feed
( 1 ) 1 i>m 0 cg @ feed
( 2 ) 2 i>m 0 cg @ feed
      2 i>m 0 cg @ feed
( 3 ) 3 i>m 0 cg @ feed
      3 i>m 0 cg @ feed
( 0 ) 0 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  eow
  5 i>m :moses  0 2 0 wire  1 3 0 wire  eow
  11 i>m :print  eow
  22 i>m :print  eow
c:end  eow

( moses: ) hex. memp @ hex.
( 11 4 ) 4 i>m 0 cg @ feed
( 22 5 ) 5 i>m 0 cg @ feed
( 22 6 ) 6 i>m 0 cg @ feed

h.s

\ need to think about location of all data: flash vs ram
\ make sure this can run from power-up, i.e. all flash-based
