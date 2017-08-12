\ <<<engine>>>
\ cr compiletoflash

: m. ( m -- ) . ;
: m! ( m a -- ) ! ;

: :print ( arg -- g )
  [: ( msg in -- ) drop  extra @ ?dup if m. then m. ;]
  1 cells new-gadget ( arg g )
  tuck g-extra m! ;

: :pass ( arg -- g )
  drop
  [: ( msg in -- ) g-emit ;]
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
  ['] swap-h 1 cells new-gadget ( arg g )
  tuck g-extra m! ;
