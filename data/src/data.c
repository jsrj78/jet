// Tight Data structures, implementation

#include "data.h"
#include <string.h>

#define TdCHUNKSIZE 8
#define TdPOOLSIZE  1000

typedef int16_t Td_Tag;

typedef union {
      int8_t c[TdCHUNKSIZE];
     uint8_t b[TdCHUNKSIZE];
     int16_t s[TdCHUNKSIZE/2];
    uint16_t u[TdCHUNKSIZE/2];
     int32_t l[TdCHUNKSIZE/4];
    uint32_t q[TdCHUNKSIZE/4];
     int64_t w[TdCHUNKSIZE/8];
    uint64_t h[TdCHUNKSIZE/8];
       void* p[TdCHUNKSIZE/sizeof(void*)];
} Td_Chunk;

static Td_Chunk tdChunks [TdPOOLSIZE];
static Td_Tag tdTags [TdPOOLSIZE+1];

static Td_Tag* tdTagP (int cid) {
    // use a convoluted way to avoid bit shifting
    return (Td_Tag*)(void*) ((uint8_t*) tdTags + cid);
}

static Td_Chunk* tdChunkP (int cid) {
    // use a convoluted way to shift bits once instead of twice
    return (Td_Chunk*)(void*) ((uint32_t*) tdChunks + cid);
}

void tdInitPool () {
    for (uint16_t i = 1; i < TdPOOLSIZE+1; ++i)
        *tdTagP(i<<1) = (int16_t) ((i-1) << 1);
}

int16_t tdChain () {
    return tdTags[TdPOOLSIZE];
}

int16_t tdAlloc () {
    int16_t cid = tdTags[TdPOOLSIZE];
    tdTags[TdPOOLSIZE] = *tdTagP(cid);
    return cid;
}

void tdDelRef (Td_Val val) {
    if ((val._ & 1) == 0) {
        int16_t cid = val._;
        *tdTagP(cid) = tdTags[TdPOOLSIZE];
        tdTags[TdPOOLSIZE] = cid;
    }
}

static void tdSetTag (int cid, int type, int len) {
    *tdTagP(cid) = (int16_t) ((1<<15) | (len<<12) | type);
}

Td_Val tdNewInt (int32_t num) {
    if (-8192 <= num && num < 8192)
        return (Td_Val) {(int16_t) ((num<<2) | 1)};
    int16_t cid = tdAlloc();
    tdSetTag(cid, 1, 0);
    tdChunkP(cid)->l[0] = num;
    return (Td_Val) {cid};
}

extern  Td_Val tdNewStr (const char* str) {
    int16_t cid = tdAlloc();
    tdSetTag(cid, 2, (int) strlen(str));
    // avoid compiler warnings about casts and consts
    tdChunkP(cid)->p[0] = (void*)(intptr_t) str;
    return (Td_Val) {cid};
}

extern  Td_Val tdNewVec (int len) {
    int16_t cid = tdAlloc();
    tdSetTag(cid, 3, len);
    // avoid compiler warnings about casts and consts
    tdChunkP(cid)->u[0] = (uint16_t) len;
    return (Td_Val) {cid};
}

extern int16_t tdSize (Td_Val val) {
    if (val._ & 1)
        return 0;
    int16_t cid = val._;
    return (*tdTagP(cid)>>12) & 0x7;
}

int32_t tdAsInt (Td_Val val) {
    if (val._ & 1)
        return val._ >> 2;
    int16_t cid = val._;
    return tdChunkP(cid)->l[0];
}
