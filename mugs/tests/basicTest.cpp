#include "mugs.h"

#include "CppUTest/TestHarness.h"

TEST_GROUP(Basic)
{
  //TEST_SETUP() {}
  //TEST_TEARDOWN() {}
};

TEST(Basic, TrivialEquality)
{
  CHECK_EQUAL(3, 1+2); // expected 1st, actual 2nd
}

TEST(Basic, ValueSize)
{
  CHECK_EQUAL(sizeof (void*), sizeof (Value));
}

TEST(Basic, EmptyMugSize)
{
  Mug<> m;
  CHECK_EQUAL(sizeof (void*), sizeof m);
}

TEST(Basic, MugsHaveInputsAndOutputs)
{
  Mug<1,2> m;
  CHECK_EQUAL(1, m.Inputs());
  CHECK_EQUAL(2, m.Outputs());
}

TEST(Basic, DerivedMug)
{
  static int lastIndex = 0;
  static Value lastValue;

  class MyMug : Mug<1> {
  public:
    void Trigger (int idx, const Value& slt) {
      lastIndex = idx;
      lastValue = slt;
    }
  } m;

  m.Trigger(1, 2);
  CHECK_EQUAL(1, lastIndex);
  CHECK_EQUAL(2, (int) lastValue);
}
