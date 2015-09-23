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

TEST(Basic, SlotSize)
{
  CHECK_EQUAL(2, sizeof (Slot));
}

TEST(Basic, EmptySlot)
{
  Slot s;
  CHECK_EQUAL(0, s.Value());
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
