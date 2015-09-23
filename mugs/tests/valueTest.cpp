#include "mugs.h"

#include "CppUTest/TestHarness.h"

TEST_GROUP(Value)
{
  //TEST_SETUP() {}
  //TEST_TEARDOWN() {}
};

TEST(Value, DefaultIsNil)
{
  Value v;
  CHECK_TRUE(v.isNil());
  CHECK_EQUAL(0, (int) v);
  CHECK_EQUAL(0, (const char*) v);
}

TEST(Value, IntType)
{
  Value v = 123;
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Value::INT, v.Type());
  CHECK_EQUAL(123, (int) v);
}

TEST(Value, StrType)
{
  Value v = "abc";
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Value::STR, v.Type());
  STRCMP_EQUAL("abc", v);
}
