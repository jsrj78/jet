#pragma once

class Slot {
  uint16_t value;

 public:
  Slot () : value (0) {}
  uint16_t Value () const { return value; }
};

class MugBase {
protected:
  virtual ~MugBase () {}
  virtual int Inputs() const =0;
  virtual int Outputs() const =0;
};

template < int I =0, int O =0 >
  class Mug : MugBase {
public:
  int Inputs() const { return I; }
  int Outputs() const { return O; }
};
