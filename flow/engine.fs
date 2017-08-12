\ <<<core>>>
\ cr compiletoflash

: ...  cr cr h.s cr ." NOT-YET!" ct quit ;

0 variable cg
0 variable gpos

0 constant _

1000 buffer: mem   \ memory pool for gadget state
mem variable memp  \ next free pointer
mem 1000 $FF fill  \ for high-water mark debugging

: alloc ( b - p )  \ allocate and clear specified number of bytes from mem pool
  3 + 3 bic  \ round up
  memp @ tuck over 0 fill
  memp +! ;

: i>m ( n -- m ) ;

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
\ cg @ feed  \ TODO test code
  wires >r begin ( msg out r: wptr )
    r@ c@ 255 <> while
    r@ h@ 12 rshift ( msg out wout r: wptr )
    over = if
      over r@ 1+ c@ $F and ( msg out msg in r: wptr )
      parent g-extra r@ c@ cells + @ feed
    then
    r> 2+ >r
  repeat rdrop 2drop ;

: new-gadget ( h x -- g )  \ construct a new gadget instance
  3 cells + alloc tuck !  dup cg !
  here over 2 cells + ! ;  \ ptr to wiring (set up after new-gadget returns) 

: circuit-h ( msg in -- )
  cells extra + @  cg @ >r cg ! 0 g-emit r> cg ! ;

: c:begin ( -- ogpos )
  gpos @  sp@ gpos ! ;

: c:count ( -- n )
  gpos @ sp@ - 4 / 2- ;

: c:end ( ogpos g* -- g )  \ construct a new circuit instance
  ['] circuit-h  c:count cells new-gadget ( ogpos g* g )
  c:count 1- 0 swap do
    2dup swap cell+ !  \ fill in parent
    tuck g-extra i cells + !
  -1 +loop
  swap gpos ! ;

: wire ( o g i -- )  \ encode a wire as 16-bit int: (o:4,i:4,g:8)
  8 lshift or swap 12 lshift or h, ;

: eow ( -- )  \ special end-of-wiring code
  $FFFF h, align ;

\ cornerstone <<<engine>>>

include gadgets.fs
