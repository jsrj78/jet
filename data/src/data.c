// Tight Data structures, implementation

#include "data.h"

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

void tdDelRef (Td_Val v) {
    if ((v._ & 1) == 0) {
        int16_t cid = v._;
        *tdTagP(cid) = tdTags[TdPOOLSIZE];
        tdTags[TdPOOLSIZE] = cid;
    }
}

Td_Val tdNewInt (int32_t n) {
    if (-8192 <= n && n < 8192)
        return (Td_Val) {(int16_t) ((n<<2) | 1)};
    int16_t cid = tdAlloc();
    tdChunkP(cid)->l[0] = n;
    return (Td_Val) {cid};
}

int32_t tdAsInt (Td_Val v) {
    if (v._ & 1)
        return v._ >> 2;
    int16_t cid = v._;
    return tdChunkP(cid)->l[0];
}
