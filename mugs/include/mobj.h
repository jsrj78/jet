// Micro objects.
#pragma once
#include <assert.h>
#include <string.h>

class Chunk {
 public:
  enum { SLACK = sizeof (void*) / 2 - 1, MAXDATA = 2 * sizeof (void*) - 2 };
  typedef enum {
    ANY, INT, FLT, CHN, // scalar values, might also use chains
    TXT, BYT, VEC, MAP, SET, // all the chained collection types
    PAD0, PAD1, PAD2, PAD3, PAD4, PAD5, PAD6 // tuple padding
  } SubTyp;

  union {
    int         num;
    void*       ptr;
    const char* str;
    int8_t      i1 [4];
    uint8_t     u1 [4];
    int16_t     i2 [2];
    uint16_t    u2 [2];
    int32_t     i4 [1];
    uint32_t    u4 [1];
    float       f4 [1];
  }        val;
  uint16_t aux [SLACK];
  uint16_t nxt;

  SubTyp type () const { return (SubTyp) (aux[SLACK-1] & 15); }
  int refs () const { return aux[SLACK-1] >> 4; }
  void incRef () { aux[SLACK-1] += 1<<4; assert(refs() != 0); }
  void decRef () { assert(refs() != 0); aux[SLACK-1] -= 1<<4; }

  void init (SubTyp t) { aux[SLACK-1] = (uint16_t) t; }
};

class Pool {
 public:
  static Chunk mem [];
  static uint16_t& size () { return mem[0].aux[Chunk::SLACK-1]; }
  static int numAllocs () { return mem[0].val.num; }

  static void init(size_t bytes) {
    size() = (uint16_t) (bytes / sizeof (Chunk));
    for (int i = 0; i < size() - 1; ++i)
      mem[i].nxt = (uint16_t) (i+1);
    // this code assumes that the memory pool starts out as all zeroes
  }
  static Chunk* allocate (int cnt =1) {
    ++mem[0].val.num; // count the number of times we've been called
    int free = mem[0].nxt;
    while (--cnt >= 0) {
      int next = mem[0].nxt;
      mem[0].nxt = mem[next].nxt;
      if (cnt == 0)
        mem[next].nxt = 0; // terminate the returned chain
    }
    return mem + free;
  }
  static void release (Chunk* p) {
    Chunk* head = p;
    while (true) {
      memset(p, 0, Chunk::MAXDATA);
      int next = p->nxt;
      if (next == 0) {
        p->nxt = mem[0].nxt;
        mem[0].nxt = (uint16_t) (head - mem);
        return;
      }
      p = mem + next;
    }
  }
};

class Val {
 public:
  uint16_t v;

  typedef enum { REF, TUP, INT, FIX } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 2) | INT) {}

  Val (const char* s) {
    size_t len = strlen(s);
    int cnt = (int) len / Chunk::MAXDATA + 2;
    Chunk* p = Pool::allocate(cnt);
    v = (uint16_t) ((p - Pool::mem) << 2) | REF;
    p->val.u2[0] = (uint16_t) len;
    while (true) {
      p = Pool::mem + p->nxt;
      if (len < Chunk::MAXDATA) {
        memcpy(p, s, len);
        p->val.u1[len] = 0;
        return;
      }
      printf("hhh\n");
      memcpy(p, s, Chunk::MAXDATA);
      s += Chunk::MAXDATA;
      len -= Chunk::MAXDATA;
    }
  }

  bool isNil () const { return v == 0; }
  Typ type () const { return (Typ) (v & 3); }
  unsigned chunk () const { return (unsigned) (v >> 2); }

  int size () const {
    if (type() != REF)
      return 0;
    Chunk* p = Pool::mem + chunk();
    return p->val.u2[0];
  }

  operator int () const { return (int16_t) v >> 2; }
  operator const char* () const {
    if (type() != REF)
      return 0;
    Chunk* p = Pool::mem + chunk();
    return (const char*) (Pool::mem + p->nxt);
  }
};
