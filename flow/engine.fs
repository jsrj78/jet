<<<core>>>
cr compiletoflash

0 variable cg

: ... cr ." NOT-YET!" quit ;

: g-extra ( g -- a )  3 cells + inline ;
: extra ( -- a )  cg @ g-extra inline ;
: parent ( g -- g )  2 cells + @ inline ;

: feed ( msg in gadget -- )  cg @ >r  dup cg !  @ execute  r> cg ! ;
: emit ( msg out gadget -- )
  ... ;

: new-gadget ( h x -- g )
  ... ;

cornerstone <<<engine>>>
include gadgets.fs
