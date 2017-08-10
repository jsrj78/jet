// Smoke test, to quickly test the test setup itself.

#include "CppUTest/TestHarness.h"

TEST_GROUP(Smoke) {};

TEST(Smoke, TrivialEquality) {
    CHECK_EQUAL(3, 1+2); // expected 1st, actual 2nd
}
