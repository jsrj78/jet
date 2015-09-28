// Micro objects.
#pragma once
#include <assert.h>

class Chunk {
 public:
  enum { SLACK = sizeof (void*) / 2 - 1, MAXDATA = 2 * sizeof (void*) - 2 };
  typedef enum {
    STRING, INTEGER,
  } SubTyp;

  union {
    int         i;
    void*       p;
    const char* s;
    int8_t      c [4];
    uint8_t     b [4];
    int16_t     h [2];
    uint16_t    w [2];
    int32_t     l [1];
    uint32_t    q [1];
    float       f [1];
  }        val;
  uint16_t aux [SLACK];
  uint16_t nxt;

  SubTyp type () const { return (SubTyp) (aux[SLACK-1] & 15); }
  int refs () const { return aux[SLACK-1] >> 4; }
  void incRef () { aux[SLACK-1] += 1<<4; assert(refs() != 0); }
  void decRef () { assert(refs() != 0); aux[SLACK-1] -= 1<<4; }
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

  typedef enum { REF, TUP, INT, FIX } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 2) | INT) {}
  Val (const char* s) : v (1 << 2) { (void) s; }

  bool isNil () const { return v == 0; }
  Typ type () const { return (Typ) (v & 3); }
  unsigned chunk () const { return (unsigned) (v >> 2); }

  operator int () const { return (int16_t) v >> 2; }
  operator const char* () const { return "abc"; }
};
