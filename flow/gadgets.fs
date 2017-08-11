<<<engine>>>
cr compiletoflash

: m. ( m -- ) . ;

: :print ( arg -- g )
  [: ( msg in -- ) drop  extra @ ?dup if m. then m. ;]
  1 cell new-gadget ( arg g )
  tuck g-extra m! ;

: :pass ( arg -- g )
  drop
  [: ( msg in -- ) cg @ emit ;]
  0 new-gadget ;

: :inlet ( arg -- g )
  drop  0 0 new-gadget ;

: outlet-h ( msg in -- )
  ... ;
: :outlet ( arg -- g )
  drop  ['] outlet-h 0 new-gadget ;

: swap-h ( msg in -- )
  ... ;
: :swap ( arg -- g )
  ['] swap-h 1 cell new-gadget ( arg g )
  tuck g-extra m! ;

