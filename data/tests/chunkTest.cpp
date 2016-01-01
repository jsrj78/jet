extern "C" {
#include "data.h"
}

#include "CppUTest/TestHarness.h"

TEST_GROUP(Chunk) {
    TEST_SETUP() {
        tdInitPool();
    }
    TEST_TEARDOWN() {
        // TODO: verify that all chunks have been released
    }
};

TEST(Chunk, HasFree) {
    CHECK(tdChain() != 0);
}

TEST(Chunk, SmallInt) {
    int oldChain = tdChain();
    Td_Val v1 = tdNewInt(1234);
    Td_Val v2 = tdNewInt(-1234);
    CHECK(tdChain() == oldChain);
    CHECK_EQUAL(1234, tdAsInt(v1));
    CHECK_EQUAL(-1234, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdChain() == oldChain);
}

TEST(Chunk, LargeInt) {
    int oldChain = tdChain();
    Td_Val v1 = tdNewInt(123456789);
    Td_Val v2 = tdNewInt(-123456789);
    CHECK(tdChain() != oldChain);
    CHECK_EQUAL(123456789, tdAsInt(v1));
    CHECK_EQUAL(-123456789, tdAsInt(v2));
    tdDelRef(v2);
    tdDelRef(v1);
    CHECK(tdChain() == oldChain);
}

TEST(Chunk, ShortStr) {
    Td_Val v = tdNewStr("abcde");
    CHECK_EQUAL(5, tdSize(v));
    STRCMP_EQUAL("abcde", (const char*) tdPeek(v));
}

TEST(Chunk, ShortVec) {
    Td_Val v = tdNewVec(3);
    CHECK_EQUAL(3, tdSize(v));
    CHECK(tdIsUndef(tdAt(v, 0)));
    CHECK(tdIsUndef(tdAt(v, 1)));
    CHECK(tdIsUndef(tdAt(v, 2)));
}

TEST(Chunk, SetShortVec) {
    Td_Val v = tdNewVec(2);
    tdSetAt(v, 0, tdNewInt(11));
    tdSetAt(v, 1, tdNewInt(22));
    CHECK_EQUAL(2, tdSize(v));
    CHECK_EQUAL(11, tdAsInt(tdAt(v, 0)));
    CHECK_EQUAL(22, tdAsInt(tdAt(v, 1)));
}
