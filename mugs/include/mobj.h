// Micro objects.
#pragma once

struct Chunk {
  union {
    int i;
    long l;
    float f;
    void* p;
    const char* s;
    uint8_t b[4];
    uint16_t u[2];
  } val;
  uint16_t aux[sizeof (void*) / 2 - 1];
  uint16_t nxt;
};

extern Chunk pool [];
static int poolSize;

class Pool {
 public:
  static void init(size_t bytes) {
    poolSize = (int) (bytes / sizeof (Chunk));
    for (int i = 0; i < poolSize; ++i)
      pool[i].nxt = (uint16_t) (i+1);
    pool[poolSize-1].nxt = 0; // last chunk is end of free chain
  }
  static Chunk* alloc (int cnt =1) {
    int free = pool[0].nxt;
    while (--cnt >= 0) {
      int next = pool[0].nxt;
      pool[0].nxt = pool[next].nxt;
      if (cnt == 0)
        pool[next].nxt = 0; // terminate the returned chain
    }
    return pool + free;
  }
};

class Val {
 public:
  uint16_t v;

  typedef enum { REF, VEC, INT, FIX } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 2) | INT) {}
  Val (const char* s) : v (1 << 2) { (void) s; }

  bool isNil () const { return v == 0; }
  Typ type () const { return (Typ) (v & 3); }
  unsigned chunk () const { return (uint16_t) (v >> 2); }

  operator int () const { return (int16_t) v >> 2; }
  operator const char* () const { return "abc"; }
};
