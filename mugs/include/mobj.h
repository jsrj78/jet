// Micro objects.
#pragma once

class Val {
 public:
  uint16_t v;
  typedef enum { PTR, INT } Typ;

  Val () : v (0) {}
  Val (int i) : v ((uint16_t) (i << 3) | INT) {}

  Typ type () const { return (Typ) (v & 7); }
  bool isNil () const { return v == 0; }

  operator int () const { return (int16_t) v >> 3; }
};
