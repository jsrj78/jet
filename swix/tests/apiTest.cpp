// Public API tests.

extern "C" {
#include "swix.h"
}

#include "CppUTest/TestHarness.h"

TEST_GROUP(Api)
{
    TEST_SETUP() {
        Swix("^");
    }
};

TEST(Api, Version) {
    CHECK_EQUAL(SWIX_VERSION, Swix("^"));
}

TEST(Api, DummyCall) {
    CHECK_EQUAL(0, Swix(""));
}

TEST(Api, EmptyStack) {
    CHECK_EQUAL(0, Swix("#"));
}

TEST(Api, PushInts) {
    CHECK_EQUAL(3, Swix("iii#", 11, 22, 33));
}

TEST(Api, PushMix) {
    CHECK_EQUAL(2, Swix("is#", 123, "abc"));
}

TEST(Api, Minus) {
    CHECK_EQUAL(-123, Swix("123-"));
}

TEST(Api, PushCount) {
    CHECK_EQUAL(3, Swix("9<8<7<#"));
}

TEST(Api, ClearCount) {
    CHECK_EQUAL(1, Swix("1"));
    CHECK_EQUAL(0, Swix("1 "));
}

TEST(Api, PopInts) {
    Swix("11<22<");
    CHECK_EQUAL(2, Swix("#"));
    CHECK_EQUAL(22, Swix(">"));
    CHECK_EQUAL(11, Swix(">"));
    CHECK_EQUAL(0, Swix("#"));
}

TEST(Api, Pack) {
    Swix("8<11<22<33<3|9<");
    CHECK_EQUAL(3, Swix("#"));
    CHECK_EQUAL(9, Swix(">"));
    CHECK_EQUAL(8, Swix(".>"));
    CHECK_EQUAL(0, Swix("#"));
}

TEST(Api, InlineStr) {
    CHECK_EQUAL(3, Swix("8<3'abc9<#"));
    CHECK_EQUAL(9, Swix(">"));
    CHECK_EQUAL(8, Swix(".>"));
}

TEST(Api, ArgAsCount) {
    CHECK_EQUAL(-123, Swix("*", -123));
}

TEST(Api, ArgSizedString) {
    CHECK_EQUAL(3, Swix("i*bi#", 11, 3, "abcd", 22));
    CHECK_EQUAL(22, Swix(">"));
    CHECK_EQUAL(11, Swix(".>"));
}
