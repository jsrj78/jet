// JET engine C core.

#include "jet.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define CHUNK_BITS   3
#define CHUNK_SIZE   (1 << CHUNK_BITS)
#define STORE_BITS   12
#define STORE_SIZE   (1 << STORE_BITS)
#define HANDLER_BITS 8
#define HANDLER_SIZE (1 << HANDLER_BITS)
#define STATE_BITS   14
#define STATE_SIZE   (1 << STATE_BITS)
#define OUTLET_BITS  12
#define OUTLET_SIZE  (1 << OUTLET_BITS)
#define GADGET_BITS  10
#define GADGET_SIZE  (1 << GADGET_BITS)
#define NETLIST_BITS 14
#define NETLIST_SIZE (1 << NETLIST_BITS)
#define IOLET_BITS   5
#define IOLET_SIZE   (1 << IOLET_BITS)

static int nthRef (int cnum, int index);

//------------------------------------------------------------------------------

#define TYPE(cnum)  tags[cnum].type
#define NEXT(cnum)  tags[cnum].next
#define MARK(cnum)  tags[cnum].mark

#define FREE        NEXT(0)

typedef struct {
    uint16_t type :3;
    uint16_t next :STORE_BITS;
    uint16_t mark :1;
} Tag_t;

enum { T_OBJ, T_STR, T_INT, T_FLT };

typedef union {
    char      m [CHUNK_SIZE];
    int32_t   l;
    double    d;
    Value_t   v;
} Chunk_t;

typedef struct {
    uint16_t flag :1;
    uint16_t ref  :NETLIST_BITS;
} NetRef_t;

typedef struct {
    uint16_t v;
    uint16_t last   :1;
    uint16_t gadget :GADGET_BITS;
    uint16_t iolet  :IOLET_BITS;
} Iolet_t;

typedef union {
    NetRef_t net;
    Iolet_t  dir;
} Outlet_t;

typedef struct {
    Iolet_t from, to;
} Wire_t;

struct Gadget_s {
    uint16_t handler :8;
    uint16_t state   :14;
    uint16_t outlets :12;
    uint16_t chunk   :12;
};

//------------------------------------------------------------------------------

Tag_t    tags     [STORE_SIZE];
Chunk_t  chunks   [STORE_SIZE];

Gadget_t gadgets  [GADGET_SIZE];
Outlet_t outlets  [OUTLET_SIZE];
int32_t  states   [STATE_SIZE];
Iolet_t  netlists [NETLIST_SIZE];

uint16_t nextGadget, nextState, nextOutlet, nextNetlist = 1;

Wire_t wires [1000]; // FIXME arbitrary high bound, not checked
int nextWire;

//------------------------------------------------------------------------------

static int countFreeChunks (void) {
    int n = 0;
    for (int first = FREE; first != 0; first = NEXT(first))
        ++n;
    return n;
}

static int chunkSize (int cnum) {
    int sz = 0;
    while (NEXT(cnum) > CHUNK_SIZE) {
        sz += CHUNK_SIZE;
        cnum = NEXT(cnum);
    }
    return sz + NEXT(cnum);
}

static void markChunk (int cnum) {
    if (cnum <= CHUNK_SIZE || MARK(cnum))
        return;
    MARK(cnum) = 1;
    if (TYPE(cnum) == 0) {
        int n = chunkSize(cnum) >> 1;
        for (int i = 0; i < n; ++i)
            markChunk(nthRef(cnum, i));
    }
    for (;;) {
        cnum = NEXT(cnum);
        if (cnum <= CHUNK_SIZE)
            break;
        markChunk(cnum);
    }
}

static void sweepChunks (void) {
    int f = 0;
    // release in reverse order, so new ones get allocated in forward order
    // for (int i = CHUNK_SIZE + 1; i < STORE_SIZE; ++i)
    for (int i = STORE_SIZE - 1; i > CHUNK_SIZE; --i)
        if (MARK(i))
            MARK(i) = 0;
        else
            NEXT(i) = f, f = i;
    FREE = f;
}

static void* chunkOffset (int cnum, int off) {
    while (off >= CHUNK_SIZE) {
        if (NEXT(cnum) <= CHUNK_SIZE)
            return 0;
        cnum = NEXT(cnum);
        off -= CHUNK_SIZE;
    }
    return chunks[cnum].m + off;
}

static int nthRef (int cnum, int index) {
    Value_t* p = chunkOffset(cnum, index << 1);
    return p != 0 && p->f ? p->i : 0;
}

static Value_t newChunk (int type, int size) {
    Value_t v;
    v.i = FREE;
    v.f = 1;
    for (;;) {
        assert(FREE != 0); // fails when out of memory
        if (size <= CHUNK_SIZE)
            break;
        size -= CHUNK_SIZE;
        FREE = NEXT(FREE);
    }
    int next = NEXT(FREE);
    NEXT(FREE) = size;
    FREE = next;
    TYPE(v.i) = type;
    return v;
}

static Value_t longToVal (long l) {
    Value_t v;
    v.i = l;
    if (v.i == l)
        v.f = 0;
    else {
        v = newChunk(T_INT, sizeof (uint32_t));
        chunks[v.i].l = l;
    }
    return v;
}

static Value_t doubleToVal (double d) {
    Value_t v = newChunk(T_FLT, sizeof (double));
    chunks[v.i].d = d;
    return v;
}

static Value_t stringToVal (const char* s) {
    int n = strlen(s) + 1;
    Value_t v = newChunk(T_STR, n);
    int cnum = v.i;
    while (n > 0) {
        strncpy(chunks[cnum].m, s, CHUNK_SIZE);
        s += CHUNK_SIZE;
        n -= CHUNK_SIZE;
        cnum = NEXT(cnum);
    }
    return v;
}

static void pushValue (int cnum, Value_t v) {
// FIXME assert(TYPE(cnum) = T_OBJ);
    while (NEXT(cnum) > CHUNK_SIZE)
        cnum = NEXT(cnum);
    if (NEXT(cnum) == CHUNK_SIZE) {
        assert(FREE != 0); // out of memory
        cnum = NEXT(cnum) = FREE;
        FREE = NEXT(cnum);
        NEXT(cnum) = 0;
    }
    *(Value_t*) (chunks[cnum].m + NEXT(cnum)) = v;
    NEXT(cnum) += sizeof (Value_t);
}

static void showObject (int cnum) {
    printf("[");
    int first = 1;
    do {
        const char* p = chunks[cnum].m;
        cnum = NEXT(cnum);
        for (int i = 0; i < cnum && i < CHUNK_SIZE; i += sizeof (Value_t)) {
            if (first)
                first = 0;
            else
                printf(" ");
            jPrint(*(Value_t*) (p + i));
        }
    } while (cnum > CHUNK_SIZE);
    printf("]");
}

static void showString (int cnum) {
    printf("'");
    do {
        printf("%.*s", CHUNK_SIZE, chunks[cnum].m);
        cnum = NEXT(cnum);
    } while (cnum > CHUNK_SIZE);
    printf("'");
}

static int objSize (int cnum) {
    assert(TYPE(cnum) == T_OBJ);
    return chunkSize(cnum) >> 1;
}

static Value_t objAt (int cnum, int index) {
    assert(index < objSize(cnum));
    Value_t* p = chunkOffset(cnum, index << 1);
    assert(p != 0);
    return *p;
}

static const char* asStr (Value_t v) {
    assert(TYPE(v.i) == T_STR);
    return chunks[v.i].m;
}

static void sendTo (Gadget_t* gptr, int inlet, Value_t msg) {
    printf(">>> %s %d:%d\n",
            jNameTable[gptr->handler], (int) (gptr - gadgets), inlet);
    jDispatchTable[gptr->handler](gptr, inlet, msg);
}

static int lookupHandler (const char* name) {
    for (const char** p = jNameTable; *p != 0; ++p)
        if (strcmp(name, *p) == 0)
            return p - jNameTable;
    assert(0); // no such gadget
    return -1;
}

static void defineGadget (Gadget_t* gp, Value_t v) {
    Value_t g = objAt(v.i, objSize(v.i) - 1);
    assert(g.f && objSize(g.i) > 0);
    const char* s = asStr(objAt(g.i, 0));
    assert(s != 0);

    int h = lookupHandler(s);
    gp->handler = h;
    gp->state = nextState;
    gp->outlets = nextOutlet;

    Config_t cfg;
    memset(&cfg, 0, sizeof cfg);
    cfg.state = states + nextState;

    jConfigTable[h](&cfg, g);

    nextOutlet += cfg.outlets;
    nextState += (cfg.stateLen + 3) >> 2; // round to uint32_t sizes
}

static void addWire (int fo, int fp, int ti, int tp) {
    assert(nextWire < sizeof wires / sizeof wires[0]); // TODO fixed bound
    Wire_t* wp = wires + nextWire++;
    wp->from.gadget = fo;
    wp->from.iolet = fp;
    wp->to.gadget = ti;
    wp->to.iolet = tp;
}

static int wireCmp (const void* p, const void* q) {
    const Wire_t *a = p, *b = q;
    return a->from.gadget < b->from.gadget ||
           a->from.iolet < b->from.iolet ||
           a->to.gadget < b->to.gadget ||
           a->to.iolet < b->to.iolet ||
           a < b ? -1 : 1;
}

static void setupOutlet (int low, int high) {
    const Wire_t* wp = wires + low;
    //printf("OUT %d:%d low %d high %d\n",
    //        wp->from.gadget, wp->from.iolet, low, high);

    Gadget_t* gp = gadgets + wp->from.gadget;
    Outlet_t* op = outlets + gp->outlets + wp->from.iolet;
    switch (high - low) {
        case 0: break;
        case 1:
            op->dir = wp->to;
            op->dir.last = 1;
            break;
        default:
            op->net.ref = nextNetlist;
            for (int i = low; i < high; ++i)
                netlists[nextNetlist++] = wires[i].to;
            netlists[nextNetlist-1].last = 1;
            break;
    }
}

static void setupNetlists () {
    qsort(wires, nextWire, sizeof wires[0], wireCmp);

    int i = 0;
    while (i < nextWire) {
        int low = i;
        while (++i < nextWire)
            if (memcmp(&wires[low].from, &wires[i].from, sizeof (Iolet_t)) != 0)
                break;
        setupOutlet(low, i);
    }
}

//------------------------------------------------------------------------------

void jPrint (Value_t v) {
    if (v.f == 0)
        printf("%d", v.i);
    else
        switch (TYPE(v.i)) {
            case T_OBJ: showObject(v.i);              break;
            case T_STR: showString(v.i);              break;
            case T_INT: printf("%d", chunks[v.i].l);  break;
            case T_FLT: printf("%g",  chunks[v.i].d); break;
            default:    printf("?");                  break;
        }
}

void jEmit (Gadget_t* gptr, int outlet, Value_t msg) {
    Outlet_t* op = outlets + gptr->outlets + outlet;
    if (op->dir.last)
        sendTo(gadgets + op->dir.gadget, op->dir.iolet, msg);
    else if (op->net.ref > 0) {
        Iolet_t* ip = netlists + op->net.ref;
        do
            sendTo(gadgets + ip->gadget, ip->iolet, msg);
        while (!(ip++)->last);
    }
}
