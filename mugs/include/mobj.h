// Micro objects.
#pragma once

struct Chunk {
  uint16_t h;
  uint16_t d[sizeof (void*) == 4 ? 1 : 3];
  union {
    int i;
    void* p;
  } v;
};

extern Chunk pool [];
static int poolSize;

class Pool {
 public:
  static void init(size_t bytes) {
    poolSize = (int) (bytes / sizeof (Chunk));
    for (int i = 0; i < poolSize; ++i)
      pool[i].h = (uint16_t) ((i+1) << 3);
    pool[poolSize-1].h = 0; // last chunk is end of free chain
  }
  static Chunk* alloc (int cnt =1) {
    int free = pool[0].h >> 3;
    while (--cnt >= 0) {
      int next = pool[0].h >> 3;
      pool[0].h = pool[next].h;
      if (cnt == 0)
        pool[next].h = 0; // terminate the returned chain
    }
    return pool + free;
  }
};

class Val {
 public:
  uint16_t v;
  typedef enum { PTR, INT } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 3) | INT) {}
  Val (const char* s) : v (1 << 3) { (void) s; }

  Typ type () const { return (Typ) (v & 7); }
  bool isNil () const { return v == 0; }

  operator int () const { return (int16_t) v >> 3; }
  operator const char* () const { return "abc"; }
};
