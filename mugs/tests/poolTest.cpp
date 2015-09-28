#include "mobj.h"
#include "CppUTest/TestHarness.h"

Chunk Pool::mem [100];

TEST_GROUP(Pool) {
  TEST_SETUP() {
    memset(Pool::mem, 0, sizeof Pool::mem);
    Pool::init(sizeof Pool::mem);
  }
};

TEST(Pool, ChunkPoolSize) {
  CHECK_EQUAL(100, Pool::size());
  CHECK_EQUAL(2 * sizeof (void*), sizeof (Chunk));
  CHECK_EQUAL(sizeof (Chunk) - 2, Chunk::MAXDATA);
}

TEST(Pool, ChunkAlignment) {
  Chunk c;
  // chn should be in the last two bytes, i.e. at offset 6 or 14
  CHECK_EQUAL(sizeof c - 2, (uint8_t*) &c.nxt - (uint8_t*) &c);
}

TEST(Pool, Alloc) {
  Chunk* p = Pool::allocate();
  CHECK_EQUAL(Pool::mem + 1, p);
  Chunk* q = Pool::allocate(2);
  CHECK_EQUAL(Pool::mem + 2, q);
  Chunk* r = Pool::allocate();
  CHECK_EQUAL(Pool::mem + 4, r);
  CHECK_EQUAL(Pool::mem + 5, Pool::allocate(0));
  CHECK_EQUAL(Pool::mem + 5, Pool::allocate(0));
  CHECK_EQUAL(5, Pool::numAllocs());
}

TEST(Pool, RefCounts) {
  Chunk* p = Pool::allocate();
  CHECK_EQUAL(0, p->refs());
  p->incRef();
  CHECK_EQUAL(1, p->refs());
  p->incRef();
  CHECK_EQUAL(2, p->refs());
  p->decRef();
  p->decRef();
  CHECK_EQUAL(0, p->refs());
}

TEST(Pool, FreeInReverseOrder) {
  Chunk* p = Pool::allocate(2); // 1 2
  Chunk* q = Pool::allocate(3); // 3 4 5
  Pool::release(q);
  CHECK_EQUAL(Pool::mem + 3, Pool::allocate(0));
  Pool::release(p);
  Chunk* r = Pool::allocate(6); // 1 2 3 4 5 6
  CHECK_EQUAL(Pool::mem + 1, r);
  CHECK_EQUAL(Pool::mem + 7, Pool::allocate(0));
  CHECK_EQUAL(5, Pool::numAllocs());
}

TEST(Pool, FreeInSameOrder) {
  Chunk* p = Pool::allocate(2); // 1 2
  Chunk* q = Pool::allocate(3); // 3 4 5
  Pool::release(p);
  CHECK_EQUAL(Pool::mem + 1, Pool::allocate(0));
  Pool::release(q);
  Chunk* r = Pool::allocate(6); // 3 4 5 1 2 6
  CHECK_EQUAL(Pool::mem + 3, r);
  CHECK_EQUAL(Pool::mem + 7, Pool::allocate(0));
  CHECK_EQUAL(5, Pool::numAllocs());
}
