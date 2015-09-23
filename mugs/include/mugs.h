#pragma once

class Value {
  uintptr_t value;

public:
  typedef enum { STR, INT } Types;

  Value () : value (STR) {}
  Value (int v) : value (((uintptr_t) v << 2) | INT) {}
  Value (const char* s) : value (((uintptr_t) s << 2) | STR) {}

  uintptr_t Raw () const { return value; }
  Types Type () const { return (Types) (value & 3); }
  bool isNil () const { return value == 0; }

  operator int () const { return (int) value >> 2; }
  operator const char* () const { return (const char*) (value >> 2); }
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
