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
    CHECK(tdFreeP()->_ != 0);
}

TEST(ChunkPool, SmallInt)
{
    int oldFree = tdFreeP()->_;
    Td_Val v1 = tdNewInt(1234);
    Td_Val v2 = tdNewInt(-1234);
    CHECK(tdFreeP()->_ == oldFree);
    CHECK_EQUAL(1234, tdAsInt(v1));
    CHECK_EQUAL(-1234, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdFreeP()->_ == oldFree);
}

TEST(ChunkPool, LargeInt)
{
    int oldFree = tdFreeP()->_;
    Td_Val v1 = tdNewInt(123456789);
    Td_Val v2 = tdNewInt(-123456789);
    CHECK(tdFreeP()->_ != oldFree);
    CHECK_EQUAL(123456789, tdAsInt(v1));
    CHECK_EQUAL(-123456789, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdFreeP()->_ == oldFree);
}
