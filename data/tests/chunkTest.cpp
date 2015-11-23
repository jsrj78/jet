extern "C" {
#include "data.h"
}

#include "CppUTest/TestHarness.h"

TEST_GROUP(Chunk) {
    TEST_SETUP() {
        tdInitPool();
    }
    //TEST_TEARDOWN() {}
};

TEST(Chunk, HasFree) {
    CHECK(tdChain() != 0);
}

TEST(Chunk, SmallInt) {
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

TEST(Chunk, LargeInt) {
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

TEST(Chunk, ShortStr) {
    Td_Val v = tdNewStr("abcde");
    CHECK_EQUAL(5, tdSize(v));
    STRCMP_EQUAL("abcde", (const char*) tdPeek(v));
}

TEST(Chunk, ShortVec) {
    Td_Val v = tdNewVec(3);
    CHECK_EQUAL(3, tdSize(v));
}
