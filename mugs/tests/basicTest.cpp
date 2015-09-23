#include "mugs.h"

#include "CppUTest/TestHarness.h"

TEST_GROUP(Basic)
{
    void setup () {
    }
};

TEST(Basic, TrivialEquality)
{
    CHECK_EQUAL(1, 1);
}
