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
  int16_t  aux [SLACK];
  uint16_t nxt;

  SubTyp type () const { return (SubTyp) (aux[SLACK-1] & 15); }
  int refs () const { return aux[SLACK-1] >> 4; }
#if 0
  void incRef () {
    aux[SLACK-1] = aux[SLACK-1] + (1<<4);
    assert(refs() != 0);
  }
  void decRef () {
    assert(refs() != 0);
    aux[SLACK-1] = aux[SLACK-1] - (1U<<4);
  }
#endif
  void incRef () { aux[SLACK-1] += 1<<4; assert(refs() != 0); }
  void decRef () { assert(refs() != 0); aux[SLACK-1] -= 1<<4; }

  void init (SubTyp t) { aux[SLACK-1] = (int16_t) t; }
};

class Pool {
 public:
  static Chunk mem [];
  static int16_t& size () { return mem[0].aux[Chunk::SLACK-1]; }

  static void init(size_t bytes) {
    size() = (int16_t) (bytes / sizeof (Chunk));
    for (int i = 0; i < size(); ++i)
      mem[i].nxt = (uint16_t) (i+1);
    mem[size()-1].nxt = 0; // last chunk is end of free chain
  }

  static Chunk* alloc (int cnt =1) {
    int free = mem[0].nxt;
    while (--cnt >= 0) {
      int next = mem[0].nxt;
      mem[0].nxt = mem[next].nxt;
      if (cnt == 0)
        mem[next].nxt = 0; // terminate the returned chain
    }
    return mem + free;
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
