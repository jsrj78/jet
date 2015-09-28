// Micro objects.
#pragma once

struct Chunk {
  union {
    int i;
    long l;
    float f;
    void* p;
    const char* s;
    uint16_t u[2];
  } val;
  uint16_t aux[sizeof (void*) == 4 ? 1 : 3];
  uint16_t chn;

  typedef enum { STRING, ARRAY, OBJECT, BYTECODE } CTyp;

  CTyp ctype () const { return (CTyp) (aux[0] & 3); }
  bool isMarked () const { return (aux[0] & 4) != 0; }
  int len () const { return aux[0] >> 3; }

  void mark () { aux[0] |= 4; }
  void unmark () { aux[0] &= (uint16_t) ~4; }
};

extern Chunk pool [];
static int poolSize;

class Pool {
 public:
  static void init(size_t bytes) {
    poolSize = (int) (bytes / sizeof (Chunk));
    for (int i = 0; i < poolSize; ++i)
      pool[i].chn = (uint16_t) (i+1);
    pool[poolSize-1].chn = 0; // last chunk is end of free chain
  }
  static Chunk* alloc (int cnt =1) {
    int free = pool[0].chn;
    while (--cnt >= 0) {
      int next = pool[0].chn;
      pool[0].chn = pool[next].chn;
      if (cnt == 0)
        pool[next].chn = 0; // terminate the returned chain
    }
    return pool + free;
  }
};

class Val {
 public:
  uint16_t v;

  typedef enum { PTR, SYM, FUNC, LONG, INT, AUX } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 3) | INT) {}
  Val (const char* s) : v (1 << 3) { (void) s; }

  bool isNil () const { return v == 0; }
  Typ type () const { return (Typ) (v & 7); }
  uint16_t chunk () const { return (uint16_t) (v >> 3); }

  operator int () const { return (int16_t) v >> 3; }
  operator const char* () const { return "abc"; }
};
