#include "mobj.h"
#include "CppUTest/TestHarness.h"

TEST_GROUP(Object) {
  TEST_SETUP() {
    // see poolTest.cpp
    const int n = 100;
    memset(Pool::mem, 0, n * sizeof Pool::mem[0]);
    Pool::init(n);
  }
};

TEST(Object, ValueSize) {
  CHECK_EQUAL(2, sizeof (Val));
}

TEST(Object, DefaultIsNil) {
  Val v;
  CHECK_TRUE(v.isNil());
  CHECK_EQUAL(Val::REF, v.type());
}

TEST(Object, IntVal) {
  Val v = -123;
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Val::INT, v.type());
  CHECK_EQUAL(-123, (int) v);
}

TEST(Object, StrVal) {
  Val v = "abc";
  CHECK_FALSE(v.isNil());
  CHECK_EQUAL(Val::REF, v.type());
  CHECK_EQUAL(3, v.size());
  STRCMP_EQUAL("abc", v);
}
