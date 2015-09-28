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
