#pragma once

class Value {
  uintptr_t value;

public:
  typedef enum { VEC, STR, INT } Types;

  Value () : value (VEC) {}
  Value (int v) : value (((uintptr_t) v << 2) | INT) {}
  Value (const char* s) : value (((uintptr_t) s << 2) | STR) {}

  ~Value () {
    if (type() == VEC)
      free((void*) value);
  }

  //uintptr_t Raw () const { return value; }
  Types type () const { return (Types) (value & 3); }
  bool isNil () const { return value == 0; }

  int len () const {
    return !isNil() && type() == VEC ? (int) (((const Value*) value)[0]) : 0;
  }

  operator int () const { return (int) value >> 2; }
  operator const char* () const { return (const char*) (value >> 2); }

  Value& operator[] (int i) { return ((Value*) value)[i]; }

  Value& operator<< (int v) { return *this << Value (v); }
  Value& operator<< (const char* s) { return *this << Value (s); }

  Value& operator<< (const Value& newVal) {
    Value* v = (Value*) value;
    if (isNil()) {
      v = (Value*) malloc(sizeof (Value));
      v[0] = 0;
    }
    int n = v[0];
    v = (Value*) realloc(v, ((unsigned) n + 2) * sizeof (Value));
    v[0] = ++n;
    v[n] = newVal;
    value = (uintptr_t) v;
    return *this;
  }
};

class MugBase {
protected:
  virtual ~MugBase () {}

  virtual int inputs () const =0;
  virtual int outputs () const =0;
  virtual void trigger (int /*idx*/, const Value& /*slt*/) {}
};

template < int I =0, int O =0 >
  class Mug : MugBase {
public:
  int inputs() const { return I; }
  int outputs() const { return O; }
};
