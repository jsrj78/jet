#include "mobj.h"
#include "CppUTest/TestHarness.h"

Chunk pool [100];

TEST_GROUP(Pool) {
  TEST_SETUP() { Pool::init(sizeof pool); }
};

TEST(Pool, ChunkPoolSize) {
  CHECK_EQUAL(100, poolSize);
  CHECK_EQUAL(2 * sizeof (void*), sizeof (Chunk));
  CHECK_EQUAL(sizeof (Chunk) - 2, Chunk::MAXDATA);
}

TEST(Pool, ChunkAlignment) {
  Chunk c;
  // chn should be in the last two bytes, i.e. at offset 6 or 14
  CHECK_EQUAL(sizeof c - 2, (uint8_t*) &c.nxt - (uint8_t*) &c);
}

TEST(Pool, Alloc) {
  Chunk* p = Pool::alloc();
  CHECK_EQUAL(&pool[1], p);
  Chunk* q = Pool::alloc(2);
  CHECK_EQUAL(&pool[2], q);
  Chunk* r = Pool::alloc();
  CHECK_EQUAL(&pool[4], r);
  CHECK_EQUAL(&pool[5], Pool::alloc(0));
  CHECK_EQUAL(&pool[5], Pool::alloc(0));
}

TEST(Pool, RefCounts) {
  Chunk c;
  CHECK_EQUAL(0, c.refs());
  c.incRef();
  CHECK_EQUAL(1, c.refs());
  c.incRef();
  CHECK_EQUAL(2, c.refs());
  c.decRef();
  c.decRef();
  CHECK_EQUAL(0, c.refs());
}
