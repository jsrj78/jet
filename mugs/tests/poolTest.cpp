#include "mobj.h"
#include "CppUTest/TestHarness.h"

Chunk pool [100];

TEST_GROUP(Pool) {
  TEST_SETUP() { Pool::init(sizeof pool); }
};

TEST(Pool, PoolSize) {
  CHECK_EQUAL(100, poolSize)
}
