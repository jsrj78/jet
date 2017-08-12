forgetram

include engine.fs

c:begin
  _      :inlet  0 1 0 wire  0 2 0 wire  eow
  11 i>m :print  eow
  22 i>m :print  eow
c:end  eow

hex.
456 i>m 0 g-feed
789 i>m 0 g-feed





c:begin
  _ :inlet  0 1 0 wire  eow
  c:begin
    _ :inlet  0 1 0 wire  eow
    _ :outlet
  c:end  0 2 0 wire  eow
  111 i>m :print  eow
c:end  eow
222 i>m 0 g-feed





\ need to think about location of all data: flash vs ram
\ make sure this can run from power-up, i.e. all flash-based
