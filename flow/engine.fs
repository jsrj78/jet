\ <<<core>>>
\ cr compiletoflash

: ...  cr cr h.s cr ." NOT-YET!" ct quit ;

0 variable cg
0 variable gpos

1000 buffer: mem   \ memory pool for gadget state
mem variable memp  \ next free pointer
mem 1000 $FF fill  \ for high-water mark debugging

: alloc ( b - p )  \ allocate and clear specified number of bytes from mem pool
  3 + 3 bic  \ round up
  3 cells +
  memp @ tuck over 0 fill
  memp +! ;

: g-extra ( g -- a )  \ convert g ptr to its extra area
  3 cells + ;
: extra ( -- a )  \ get extra area of current g
  cg @ g-extra ;
: parent ( -- g )  \ get parent of current g
  cg @ cell+ @ ;
: wires ( -- w )  \ get wires of current g
  cg @ 2 cells + @ ;

: feed ( msg in gadget -- )  \ feed a message to given gadget inlet
  cg @ >r  dup cg !  @ execute  r> cg ! ;

: g-emit ( msg out -- )  \ send a message to an outlet of current g
  cg @ feed ;  \ TODO test code

: new-gadget ( h x -- g )  \ construct a new gadget instance
  alloc tuck ! ;

: circuit-h ( msg in -- )
  cg @ g-extra cell+ @ feed ;  \ TODO hard-wired to send to gadget #1 for now

: c:begin ( -- ogpos )
  gpos @  sp@ gpos ! ;

: c:count ( -- n )
  gpos @ sp@ - 4 / 2- ;

: c:end ( ogpos g* -- g )  \ construct a new circuit instance
  ['] circuit-h  c:count cells new-gadget ( ogpos g* g )
  cr
  c:count 1- 0 swap do
    tuck g-extra i cells + !
  -1 +loop
  swap gpos ! ;

: add-g ( c g w -- c )  \ add a gadget to a circuit, including outlet wiring
  ... ;

: wire ( o g i -- )  \ encode a wire as 16-bit int: (o:4,i:4,g:8)
  8 lshift or swap 12 lshift or h, ;

: eow ( -- )  \ special end-of-wiring code
  $FFFF h, align ;

\ cornerstone <<<engine>>>

include gadgets.fs
