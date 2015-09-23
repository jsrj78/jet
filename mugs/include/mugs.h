#pragma once

class Slot {
  uintptr_t value;

 public:
  Slot (uintptr_t v =0) : value (v) {}
  uintptr_t Value () const { return value; }
};

class MugBase {
protected:
  virtual ~MugBase () {}
  virtual int Inputs () const =0;
  virtual int Outputs () const =0;
  virtual void Trigger (int /*idx*/, const Slot& /*slt*/) {}
};

template < int I =0, int O =0 >
  class Mug : MugBase {
public:
  int Inputs() const { return I; }
  int Outputs() const { return O; }
};
