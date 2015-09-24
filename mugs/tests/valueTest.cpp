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
  CHECK_EQUAL(Value::VEC, v.type());
  CHECK_EQUAL(0, v.len());
  CHECK_EQUAL(0, (int) v);
  CHECK_EQUAL(0, (const char*) v);
}

TEST(Value, IntType)
{
  Value v = -123;
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Value::INT, v.type());
  CHECK_EQUAL(-123, (int) v);
}

TEST(Value, StrType)
{
  Value v = "abc";
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Value::STR, v.type());
  STRCMP_EQUAL("abc", v);
}

TEST(Value, VecType)
{
  Value v;
  v << -123 << "abc";
  CHECK_EQUAL(2, v.len());
  CHECK_EQUAL(-123, (int) v[1]);
  STRCMP_EQUAL("abc", v[2]);
}

TEST(Value, ChangeVec)
{
  Value v;
  v << -123 << "abc";
  v[1] = "defg";
  v[2] = 456;
  CHECK_EQUAL(2, v.len());
  STRCMP_EQUAL("defg", v[1]);
  CHECK_EQUAL(456, (int) v[2]);
}
