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

TEST(Basic, TagSize)
{
    CHECK_EQUAL(2, sizeof (Td_Tag));
}

TEST(Basic, ChunkSize)
{
    CHECK_EQUAL(TdCHUNKSIZE, sizeof (Td_Chunk));
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
    CHECK(*tdFreeP() != 0);
}

TEST(ChunkPool, SmallInt)
{
    Td_Val v1 = tdNewInt(1234);
    CHECK_EQUAL(1234, tdAsInt(v1));
    Td_Val v2 = tdNewInt(-1234);
    CHECK_EQUAL(-1234, tdAsInt(v2));
}

TEST(ChunkPool, LargeInt)
{
    uint16_t oldFree = *tdFreeP();
    Td_Val v = tdNewInt(123456789);
    CHECK_EQUAL(123456789, tdAsInt(v));
    CHECK(*tdFreeP() != oldFree);
}
