// JET engine C core.

#include "jet.h"

#include <assert.h>
#include <stdio.h>
#include <string.h>

#define CHUNK_BITS  3
#define CHUNK_SIZE  (1 << CHUNK_BITS)
#define STORE_BITS  12
#define STORE_SIZE  (1 << STORE_BITS)

static int nthRef (int cnum, int index);

//------------------------------------------------------------------------------

#define TYPE(cnum)  tags[cnum].type
#define NEXT(cnum)  tags[cnum].next
#define MARK(cnum)  tags[cnum].mark

#define FREE        NEXT(0)

struct {
    uint16_t type :3;
    uint16_t next :STORE_BITS;
    uint16_t mark :1;
} tags [STORE_SIZE];

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

//------------------------------------------------------------------------------

enum { T_OBJ, T_STR, T_INT, T_FLT };

union {
    char      m [CHUNK_SIZE];
    int32_t   l;
    double    d;
    Value_t   v;
} chunks [STORE_SIZE];

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

//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------

void jEmit (Gadget_t* gptr, int inlet, Value_t msg) {
}
