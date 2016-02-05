"use strict";
let p = (...args) => console.log("p:", ...args)
p("hello", "world")

let events = []

let show_h = (gi, hv) => {
    console.log("show:", hv)
}

let show_t = { name: "show", handler: show_h, outCount: 0, }
let show_i = { id: null, ins: [null], outs: null, type: show_t, }

let metro_h = (gi, hv) => {
    if (!hv) {
        p("PONG!")
        gi.outs[0] = true
        return
    }
    //clearInterval(gi.l_tid)
    //gi.l_tid = setInterval(() => {
    setInterval(() => {
        p("ping!")
        events.push([gi,null])
    }, hv)
}

let metro_t = { name: "metro", handler: metro_h, outCount: 1, }
let metro_i = { id: null, ins: [null], outs: null, type: metro_t, }

let top_h = (gi, hv) => {
    p("top", hv)
    while (events.length > 0) {
        let e = events.shift()
        p("shifted", e)
    }
}

let top_t = { name: "top", handler: top_h, outcount: 0, }
let top_i = { id: null, ins: [null], outs: null, type: top_t, }

top_i.gadgets = [
    null, // self
    metro_i,
    show_i,
]
top_i.wires = [
    [],
    [ [2,0] ],
    [],
]

show_i.parent = top_i
metro_i.parent = top_i

top_i.id = 0
metro_i.id = 1
show_i.id = 2

p("show_i", show_i)
show_h(show_i, 123)
metro_h(metro_i, 1000)
top_h(top_i, null)
//console.log(JSON.stringify(top_i, null, 2))
setTimeout(() => {
    top_h(top_i, null)
}, 1500)
