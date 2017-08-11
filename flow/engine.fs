\ <<<core>>>
\ cr compiletoflash

0 variable cg

: ... cr ." NOT-YET!" quit ;

: g-extra ( g -- a )  \ convert g ptr to its extra area
  3 cells + ;
: extra ( -- a )  \ get extra area of current g
  cg @ g-extra ;
: parent ( g -- g )  \ get parent of specified g
  2 cells + @ ;

: feed ( msg in gadget -- )  \ feed a message to given gadget inlet
  cg @ >r  dup cg !  @ execute  r> cg ! ;
: emit ( msg out gadget -- )
  ... ;

: new-gadget ( h x -- g )  \ construct a new gadget instance
  ... ;

: circuit-h ( msg in -- )
  ... ;
: new-circuit ( n -- g )  \ construct a new circuit instance
  ['] circuit-h swap  cells new-gadget ;

: add-g ( c g w n -- c )  \ add a gadget to a circuit, including outlet wiring
  ... ;

: wire ( o g i -- w )  \ encode a wire as 16-bit int: (o:4,i:4,g:8)
  8 lshift or swap 12 lshift or ;

$FFFF constant EOW  \ special end-of-wiring code

\ cornerstone <<<engine>>>

include gadgets.fs
