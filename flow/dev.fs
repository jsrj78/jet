forgetram

include engine.fs

c:begin
  _ :inlet  0 1 0 wire  0 2 0 wire  eow
  11 i>m :print  eow
  22 i>m :print  eow
c:end  eow

hex. memp @ hex.
456 i>m 0 cg @ feed
789 i>m 0 cg @ feed

c:begin
  _ :inlet  0 1 0 wire  eow
  c:begin
    _ :inlet  0 1 0 wire  eow
    _ :outlet
  c:end  1 2 0 wire  eow
  111 i>m :print  eow
c:end  eow

hex. memp @ hex.
222 i>m 0 cg @ feed

h.s

\ need to think about location of all data: flash vs ram
\ make sure this can run from power-up, i.e. all flash-based
