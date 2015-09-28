#include "mobj.h"
#include "CppUTest/TestHarness.h"

Chunk pool [100];

TEST_GROUP(Pool) {
  TEST_SETUP() { Pool::init(sizeof pool); }
};

TEST(Pool, PoolSize) {
  CHECK_EQUAL(100, poolSize)
}

TEST(Pool, Alloc) {
  Chunk* p = Pool::alloc();
  CHECK_EQUAL(&pool[1], p)
  Chunk* q = Pool::alloc(2);
  CHECK_EQUAL(&pool[2], q)
  Chunk* r = Pool::alloc();
  CHECK_EQUAL(&pool[4], r)
  CHECK_EQUAL(&pool[5], Pool::alloc(0))
  CHECK_EQUAL(&pool[5], Pool::alloc(0))
}
