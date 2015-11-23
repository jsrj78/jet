// Tight Data structures
#pragma once

#define TdCHUNKSIZE 8
#define TdPOOLSIZE  1000

union Td_Val {
    struct {
        uint16_t x :1;
         int16_t v :14;
        uint16_t y :1;
    } _;

     int16_t s;
    uint16_t u;
};

union Td_Tag {
    struct {
        uint16_t x :1;
        uint16_t r :6;
        uint16_t l :1;
        uint16_t f :4;
        uint16_t d :1;
        uint16_t t :3;
    } _;

     int16_t s;
    uint16_t u;
};

union Td_Chunk {
      int8_t c[TdCHUNKSIZE];
     uint8_t b[TdCHUNKSIZE];
     int16_t s[TdCHUNKSIZE/2];
    uint16_t u[TdCHUNKSIZE/2];
     int32_t l[TdCHUNKSIZE/4];
    uint32_t q[TdCHUNKSIZE/4];
     int64_t w[TdCHUNKSIZE/8];
    uint64_t h[TdCHUNKSIZE/8];
      Td_Tag t[TdCHUNKSIZE/2];
};

extern Td_Chunk tdChunks [];
extern   Td_Tag tdTags [];

extern      void tdInitPool ();
extern uint16_t* tdFreeP ();
extern  uint16_t tdAlloc ();
extern    Td_Val tdNewInt (int32_t v);
extern   int32_t tdAsInt (Td_Val v);

//------------------------------------------------------------------------------

Td_Chunk tdChunks [TdPOOLSIZE];
Td_Tag tdTags [TdPOOLSIZE+1];

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
        tdTagP(i<<1)->u = (uint16_t) ((i-1) << 1);
}

uint16_t* tdFreeP () {
    return &tdTags[TdPOOLSIZE].u;
}

uint16_t tdAlloc () {
    uint16_t cid = *tdFreeP();
    *tdFreeP() = tdTagP(cid)->u;
    return cid;
}

Td_Val tdNewInt (int32_t n) {
    if (-8192 <= n && n < 8192)
        return {{1, (int16_t) n, 1}};
    uint16_t cid = tdAlloc();
    tdChunkP(cid)->l[0] = n;
    return {.u = cid};
}

int32_t tdAsInt (Td_Val v) {
    if (v._.y)
        return v._.v;
    uint16_t cid = v.u;
    return tdChunkP(cid)->l[0];
}
