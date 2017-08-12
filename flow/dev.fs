forgetram

include engine.fs

here hex. h.s
c:begin
  0   :inlet  0 1 0 wire  0 2 0 wire  eow
  123 :print  eow
  321 :print  eow
c:end  eow
here hex. h.s

dup 20 dump

dup 456 0 rot feed
    789 0 rot feed





\ new design, more behind-the-scenes magic:
\
\   c:begin
\     _ :inlet  0 1 0 wire  eow
\     c:begin
\       _ :inlet  0 1 0 wire  eow
\       _ :outlet
\     c:end  0 2 0 wire  eow
\     111 i>m :print  eow
\   c:end  eow
\   222 i>m 0 g-feed





\ another example:
\
\   c:begin
\     _ :inlet  0 1 0 wire  0 2 0 wire  eow
\     11 i>m :print  eow
\     22 i>m :print  eow
\   c:end  eow
\   789 i>m 0 g-feed





\ need to think about location of all data: flash vs ram
\ make sure this can run from power-up, i.e. all flash-based
