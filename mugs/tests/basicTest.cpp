#include "mugs.h"
#include "CppUTest/TestHarness.h"

static int lastIndex;
static Value lastValue;

TEST_GROUP(Basic)
{
  class BasicMug : public Mug<1,0> {
    void trigger (int idx, const Value& val) {
      lastIndex = idx;
      lastValue = val;
    }
  };

  TEST_SETUP() {
    lastIndex = 0;
    lastValue = 0;
  }
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
  Mug<0,0> m;
  CHECK_EQUAL(2 * sizeof (void*), sizeof m);
}

TEST(Basic, MugsHaveInputsAndOutputs)
{
  Mug<1,2> m;
  CHECK_EQUAL(1, m.inputs());
  CHECK_EQUAL(2, m.outputs());
}

TEST(Basic, BasicMug)
{
  BasicMug m;
  m.feed(1, 2);
  CHECK_EQUAL(1, lastIndex);
  CHECK_EQUAL(2, (int) lastValue);
}
