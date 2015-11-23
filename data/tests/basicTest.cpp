extern "C" {
#include "data.h"
}

#include "CppUTest/TestHarness.h"

TEST_GROUP(Basic) {};

TEST(Basic, TrivialEquality) {
    CHECK_EQUAL(3, 1+2); // expected 1st, actual 2nd
}

TEST(Basic, ValueSize) {
    CHECK_EQUAL(2, sizeof (Td_Val));
}

TEST(Basic, UndefinedValue) {
    Td_Val v = {0};
    CHECK(tdIsUndef(v));
}
