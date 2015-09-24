#include "mugs.h"
#include "CppUTest/TestHarness.h"

static Value lastValue;

class TraceMug : public Mug<1,0> {
  void trigger (int, const Value& val) { lastValue = val; }
};

class PassMug : public Mug<1,1> {
  void trigger (int idx, const Value& val) { send(idx, val); }
};

static TraceMug tm;
static PassMug pm;

MugBase* const mugs [] = { &pm, &tm };

const uint8_t flows [] = {
  2, // mug count
  // mug 0:
  1, 1, 1, // out 0.1 -> in 1.1
  // mug 1:
};

TEST_GROUP(Flow)
{
  TEST_SETUP() {
    lastValue = 0;
  }
  //TEST_TEARDOWN() {}
};

TEST(Flow, TraceMug)
{
  tm.feed(1, 2);
  CHECK_EQUAL(2, (int) lastValue);
}

TEST(Flow, MinimalFlow)
{
  pm.feed(1, 2);
  CHECK_EQUAL(2, (int) lastValue);
}
