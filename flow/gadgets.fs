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
  parent g-extra
  begin ( msg n gpp )
    dup @ cg @ <> while
    swap 1+ swap
    cell+
  repeat drop
  cg @ >r parent cg ! g-emit r> cg ! ;

: :outlet ( arg -- g )
  drop  ['] outlet-h 0 new-gadget ;

: swap-h ( msg in -- )
  0= if
    1 g-emit
    extra @ 0 g-emit
  else
    extra m!
  then ;

: :swap ( arg -- g )
  ['] swap-h 1 cells new-gadget ( arg g )
  tuck g-extra m! ;

: change-h ( msg in -- )
  drop dup extra @ <> if
    dup extra !
    dup 0 g-emit
  then drop ;

: :change ( arg -- g )
  drop  ['] change-h 1 cells new-gadget ( arg g )
  -1 i>m extra m! ;

: moses-h ( msg in -- )
  drop dup extra @ < if 0 else 1 then g-emit ;

: :moses ( arg -- g )
  ['] moses-h 1 cells new-gadget ( arg g )
  tuck g-extra m! ;
