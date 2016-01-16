// Swix memory management.

#include "defs.h"
#include <string.h>

// needed for cpputest-3.7.1 with gcc 4.8.2 on mohse
//#pragma GCC diagnostic ignored "-Wpedantic"

uint8_t gcCount;
static uint8_t gcLock;

enum { T_INT, T_STR, T_VEC };

#define MAX_PAYLOAD (CHUNK_SIZE - 2 * sizeof (Obj))

typedef struct {
    Obj dummy [CHUNK_SIZE / sizeof (Obj) - 2];
    uint16_t size :12;
    uint16_t type :4;
    Obj next;
} Head;

typedef union Chunk {
    Head     x;                             // header/extend
    char     c [CHUNK_SIZE];                // char
    uint8_t  b [CHUNK_SIZE];                // byte
    int16_t  i [CHUNK_SIZE/2];              // int
    uint16_t u [CHUNK_SIZE/2];              // unsigned
    int32_t  l [CHUNK_SIZE/4];              // long
    uint32_t q [CHUNK_SIZE/4];              // quad
    float    f [CHUNK_SIZE/4];              // float
#if CHUNK_SIZE > 8
    int64_t  w [CHUNK_SIZE/8];              // wide
    uint64_t h [CHUNK_SIZE/8];              // huge
    double   d [CHUNK_SIZE/8];              // double
#endif
    void*    p [CHUNK_SIZE/sizeof(void*)];  // pointer
    char*    s [CHUNK_SIZE/sizeof(char*)];  // string
    Obj      o [CHUNK_SIZE/sizeof(Obj)];    // object
} Chunk;

#define CHUNK(n) ((Chunk*) swixPool)[n]
#define OBJ(o) CHUNK((o)._>>1)

static Obj freeSlot;

Obj specialObj (int n) {
    Obj r = { (int16_t) (n<<1) };
    return r;
}

static void markObj (Obj o) {
    if (!IsRef(o))
        return;
    Obj p = OBJ(o).x.next;
    if (p._ & 1)
        return;
    ++OBJ(o).x.next._;
    // traverse and mark next chain
    while (!IsNil(p) && (OBJ(p).x.next._ & 1) == 0) {
        Obj q = OBJ(p).x.next;
        ++OBJ(p).x.next._;
        p = q;
    }
    // travere and mark all referenced vector objects
    if (IsVec(o)) {
        int n = Size(o);
        for (int i = 0; i < n; ++i)
            markObj(At(o, i)); // recurse
    }
}

static void sweepMem (int start) {
    freeSlot = NilVal();
    for (int n = (int) swixSize / CHUNK_SIZE; --n >= start; ) {
        if (--CHUNK(n).x.next._ & 1) {
            CHUNK(n).x.next = freeSlot;
            freeSlot = specialObj(n);
        }
    }
}

static Obj newChunk (void) {
    if (IsNil(freeSlot) && !gcLock) {
        for (int i = 1; i < 5; ++i)
            markObj(specialObj(i));
        sweepMem(5);
        ++gcCount;
    }
    Obj r = freeSlot;
    freeSlot = OBJ(freeSlot).x.next;
    memset(OBJ(r).c, 0, sizeof (Chunk));
    return r; // this is nil when memory is exhausted!
}

Obj Init (void) {
    gcCount = 0;
    memset(swixPool, 0, swixSize);
    sweepMem(1);        // ends with freeSlot being 1
    NewStrN(0, 0);      // 1: tag
    boxedNewInt(0);     // 2: false
    boxedNewInt(1);     // 3: true
    return NewVec();    // 4: root vector
}

Obj boxedNewInt (int n) {
    Obj r = newChunk();
    OBJ(r).x.type = T_INT;
    OBJ(r).l[0] = n;
    return r;
}

int boxedAsInt (Obj o) {
    return OBJ(o).l[0];
}

Obj NewVec (void) {
    Obj r = newChunk();
    OBJ(r).x.type = T_VEC;
    return r;
}

int boxedType (Obj o) {
    return IsRef(o) ? OBJ(o).x.type : -1;
}

static void* offsetInChunk (Obj o, size_t off, int extend) {
    void* result = 0;
    // lock down so we don't GC and risk re-using a just-allocated chunk
    gcLock = 1;
    Chunk* p = &OBJ(o);
    for (;;) {
        if (off < MAX_PAYLOAD) {
            result = p->c + off;
            break;
        }
        off -= MAX_PAYLOAD;
        if (IsNil(p->x.next)) {
            if (!extend)
                break;
            p->x.next = newChunk();
            if (IsNil(p->x.next))
                break; // memory exhausted
        }
        p = &OBJ(p->x.next);
    }
    gcLock = 0;
    return result;
}

Obj NewStrN (const char* s, size_t n) {
    Obj r = newChunk();
    if (!IsNil(r)) {
        OBJ(r).x.type = T_STR;
        OBJ(r).x.size = 0xFFF & n; // FIXME: depends on size's definition
        if (offsetInChunk(r, n, 1) == 0) // extend and zero-fill all required space
            r = NilVal(); // memory exhausted
        else
            for (Obj q = r; n > 0; q = OBJ(q).x.next) {
                size_t m = n;
                if (m > MAX_PAYLOAD)
                    m = MAX_PAYLOAD;
                memcpy(OBJ(q).c, s, m);
                s += m;
                n -= m;
            }
    }
    return r;
}

Obj NewStr (const char* s) {
    return NewStrN(s, strlen(s));
}

int Size (Obj o) {
    return IsRef(o) ? OBJ(o).x.size : -1;
}

const char* AsStr (Obj o, char* buf, size_t len) {
    if (!IsStr(o))
        return 0;
    int n = Size(o);
    if (buf == 0)
        return (size_t) n < MAX_PAYLOAD ? OBJ(o).c : 0;
    ++n; // add trailing null byte
    if ((size_t) n > len)
        n = (int) len;
    for (int i = 0; i < n; ++i)
        buf[i] = (char) AsInt(At(o, i));
    return buf;
}

Obj Append (Obj v, Obj o) {
    size_t pos = OBJ(v).x.size;
    if (IsVec(v)) {
        Obj* p = offsetInChunk(v, pos++ * sizeof (Obj), 1);
        if (p == 0)
            return NilVal();
        *p = o;
    } else if (IsStr(v)) {
        if (IsStr(o)) {
            int n = Size(o);
            for (int i = 0; i < n; ++i)
                Append(v, At(o, i)); // bail out, use recursion
            return v;
        }
        char* p = offsetInChunk(v, pos++, 1);
        if (p == 0)
            return NilVal();
        *p = (char) AsInt(o);
        if (offsetInChunk(v, pos, 1) == 0) // make sure there's a trailing zero
            return NilVal();
    }
    OBJ(v).x.size = 0xFFF & pos; // FIXME: depends on size's definition
    return v;
}

Obj At (Obj o, int n) {
    if (n < 0)
        n += Size(o);
    if (IsStr(o)) {
        const char* p = offsetInChunk(o, (size_t) n, 0);
        if (p != 0)
            return NewInt(*p);
    } else if (IsVec(o)) {
        Obj* p = offsetInChunk(o, (unsigned) n * sizeof (Obj), 0);
        if (p != 0)
            return *p;
    }
    return NilVal();
}

void Drop (Obj o) {
    // TODO: check type, and release chunks as needed
    if (OBJ(o).x.size > 0)
        --OBJ(o).x.size;
}

Obj Pack (Obj v, int n) {
    Obj r = NewVec();
    for (int i = 0; i < n; ++i)
        Append(r, At(v, i-n));
    while (--n >= 0)
        Drop(v);
    return r;
}
