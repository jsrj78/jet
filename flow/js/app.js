"use strict";
let p = (...args) => console.log("p:", ...args)
p("hello", "world")

let events = []

let timer = (gi, hv) => {
    let id = setInterval(() => {
        p("ping!")
        events.push([gi,null])
    }, hv)
    // don't return, this creates a circular structure, can't JSON.stringify
    //return id
}

let show_h = (gi, hv) => { console.log("show:", hv) }

let show_t = { name: "show", handler: show_h, outCount: 0, }
let show_i = { inputs: [null], outVec: null, type: show_t, }

let metro_h = (gi, hv) => { clearInterval(gi.l_id); gi.l_id = timer(gi, hv) }

let metro_t = { name: "metro", handler: metro_h, outCount: 1, }
let metro_i = { inputs: [null], outVec: null, type: metro_t, }

let top_t = { name: "top", handler: null, outcount: 0, }
let top_i = {
    inputs: [null],
    outVec: null,
    type: top_t,
    gadgets: [show_i, metro_i],
    wires: [
        [], // show
        [ [0,0] ], // metro
    ],
}

let top_h = (gi, hv) => {
    p("top", hv)
}
top_t.handler = top_h

p("show_i", show_i)
show_h(show_i, 123)
metro_h(metro_i, 1000)
top_h(top_i, null)
//console.log(JSON.stringify(top_i, null, 2))
