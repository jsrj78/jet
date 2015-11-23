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
    CHECK_EQUAL(2, sizeof (Td_Val));
}

TEST_GROUP(ChunkPool)
{
    TEST_SETUP() {
        tdInitPool();
    }
    //TEST_TEARDOWN() {}
};

TEST(ChunkPool, HasFree)
{
    CHECK(tdChain() != 0);
}

TEST(ChunkPool, SmallInt)
{
    int oldFree = tdChain();
    Td_Val v1 = tdNewInt(1234);
    Td_Val v2 = tdNewInt(-1234);
    CHECK(tdChain() == oldFree);
    CHECK_EQUAL(1234, tdAsInt(v1));
    CHECK_EQUAL(-1234, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdChain() == oldFree);
}

TEST(ChunkPool, LargeInt)
{
    int oldFree = tdChain();
    Td_Val v1 = tdNewInt(123456789);
    Td_Val v2 = tdNewInt(-123456789);
    CHECK(tdChain() != oldFree);
    CHECK_EQUAL(123456789, tdAsInt(v1));
    CHECK_EQUAL(-123456789, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdChain() == oldFree);
}
