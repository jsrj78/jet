extern "C" {
#include "data.h"
}

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
  CHECK_EQUAL(2, sizeof (TdValue));
}

TEST(Basic, TagSize)
{
  CHECK_EQUAL(2, sizeof (TdTag));
}

TEST(Basic, ChunkSize)
{
  CHECK_EQUAL(Td_CHUNKSIZE, sizeof (TdChunk));
}
