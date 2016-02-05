"use strict"
let p = (...args) => console.log("p:", ...args)
p("hello", "world")

let pending_g = []
let pending_w = []

let show_h = (gInst, hv) => {
    console.log("show:", hv)
}

let show_t = { name: "show", handler: show_h, outCount: 0, }
let show_i = { id: null, ins: [null], outs: null, type: show_t, }

let metro_h = (gInst, hv) => {
    if (!hv) {
        p("PONG!")
        gInst.outs[0] = true
        return
    }
    //clearInterval(gInst.l_tid)
    //gInst.l_tid = setInterval(() => {
    setInterval(() => {
        p("ping!")
        pending_g.unshift(gInst)
        pending_w.unshift(null)
    }, hv)
}

let metro_t = { name: "metro", handler: metro_h, outCount: 1, }
let metro_i = { id: null, ins: [null], outs: null, type: metro_t, }

let activate = (gInst, hv) => {
    p("activate", gInst.type.name+"#"+gInst.id, "val:", hv)
    if (gInst.outs)
        p("re-activating, but already active???")
    gInst.outs = []
    for (let i = 0; i < gInst.type.outCount; ++i)
        gInst.outs.push(undefined)
    gInst.type.handler(gInst, hv)
}

let top_h = (gInst, hv) => {
    while (pending_g.length > 0) {
        let n = pending_g.length - 1
        let pg = pending_g[n]
        let pw = pending_w[n]
        let hv = null
        if (pw === null) {
            pending_g.pop()
            pending_w.pop()
        } else {
            let nout = pg.outs.length - 1
            let circ = pg.parent
            let wnet = circ.wires[nout]
            if (pw >= wnet.length) {
                pg.outs.pop()
                if (pg.outs.length == 0) {
                    pg.outs = null // marks the end of this gadget's activity
                    pending_g.pop()
                    pending_w.pop()
                }
                continue
            }
            gNum = wnet[pw++]
            gPad = wnet[pw++]
            pending_w[n] = pw
            pg = gInst.gadgets[gNum]
            hv = pg.outs[nout]
            if (gPad > 0)
                continue // not a hot input
        }
        activate(pg, hv)
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
setInterval(() => {
    top_h(top_i, null)
}, 100)
