// Tight Data structures
#pragma once

#define Td_CHUNKSIZE 8

union TdValue {
     int16_t s;
    uint16_t u;
    struct {
        uint16_t x :1;
        uint16_t v :14;
        uint16_t y :1;
    } _;
};

union TdTag {
     int16_t s;
    uint16_t u;
    struct {
        uint16_t x :1;
        uint16_t r :6;
        uint16_t l :1;
        uint16_t f :4;
        uint16_t d :1;
        uint16_t t :3;
    } _;
};

union TdChunk {
      int8_t c[Td_CHUNKSIZE];
     uint8_t b[Td_CHUNKSIZE];
     int16_t s[Td_CHUNKSIZE/2];
    uint16_t u[Td_CHUNKSIZE/2];
     int32_t l[Td_CHUNKSIZE/4];
    uint32_t q[Td_CHUNKSIZE/4];
     int64_t w[Td_CHUNKSIZE/8];
    uint64_t h[Td_CHUNKSIZE/8];
       TdTag t[Td_CHUNKSIZE/2];
};
