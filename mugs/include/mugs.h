// Micro gadgets.
#pragma once
#include <assert.h>

class Value {
  uintptr_t value;

 public:
  typedef enum { VEC, STR, INT } Types;

  Value () : value (VEC) {}
  Value (int v) : value (((uintptr_t) v << 2) | INT) {}
  Value (const char* s) : value (((uintptr_t) s << 2) | STR) {}

  // FIXME need to recursively release the vector's items as well
  ~Value () { if (type() == VEC) free((void*) value); }

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
    Value* vecp = (Value*) value;

    int newLen = len() + 1;
    vecp = (Value*) realloc(vecp, ((unsigned) newLen + 1) * sizeof (Value));
    vecp[0] = newLen;
    vecp[newLen] = newVal;

    value = (uintptr_t) vecp;
    return *this;
  }
};

class MugBase;
extern MugBase* const mugs [];
extern const uint8_t flows [];

class MugBase {
  uint16_t offset;

  void initOffsets () {
    const uint8_t* p = flows;
    int mugCnt = *p++;
    for (int mugIdx = 0; mugIdx < mugCnt; ++mugIdx) {
      mugs[mugIdx]->offset = (uint16_t) (p - flows);
    }
  }

 protected:
  MugBase () : offset (0) {}
  virtual ~MugBase () {}

  virtual int inputs () const =0;
  virtual int outputs () const =0;
  virtual void trigger (int /*idx*/, const Value& /*val*/) {}

  void send (int idx, const Value& val) {
    if (offset == 0)
      initOffsets();
    // FIXME hardcoded i.s.o. using flow connections
    assert(offset == 1);
    mugs[1]->feed(idx, val);
  }

 public:
  void feed (int idx, const Value& val) { trigger(idx, val); }
};

template < int I, int O >
class Mug : public MugBase {
 public:
  int inputs() const { return I; }
  int outputs() const { return O; }
};
