// Basic datatype tests.

extern "C" {
#include "defs.h"
}

#include "CppUTest/TestHarness.h"

static char tmpBuf[50];

TEST_GROUP(Basic)
{
};

TEST(Basic, DataTypeSizes) {
#if REF_SMALL
    CHECK_EQUAL(2, sizeof (Obj));
#else
    CHECK_EQUAL(4, sizeof (Obj));
#endif
}

TEST(Basic, Nil) {
    CHECK(IsNil(NilVal()));
}

TEST(Basic, CreateSmallInt) {
    Obj v = NewInt(123);
    CHECK(IsInt(v));
    CHECK(!IsRef(v));
    CHECK(!IsNil(v));
    CHECK_EQUAL(123, AsInt(v));
}

TEST_GROUP(BasicPool)
{
    Obj root;

    TEST_SETUP() {
        root = Init();
    }
};

TEST(BasicPool, CreateLargeInt) {
    Obj v = NewInt(123456789);
    CHECK(IsInt(v));
    CHECK(IsRef(v));
    CHECK(!IsNil(v));
    CHECK_EQUAL(123456789, AsInt(v));
}

TEST(BasicPool, CreateSmallString) {
    Obj v = NewStr("abc");
    CHECK(IsStr(v));
    CHECK(IsRef(v));
    CHECK(!IsNil(v));
    CHECK_EQUAL(3, Size(v));
    STRCMP_EQUAL("abc", AsStr(v, tmpBuf, sizeof tmpBuf));
}

TEST(BasicPool, CreateSizedString) {
    Obj v = NewStrN("def", 2);
    CHECK(IsStr(v));
    CHECK(IsRef(v));
    CHECK_EQUAL(2, Size(v));
    STRCMP_EQUAL("de", AsStr(v, tmpBuf, sizeof tmpBuf));
}

TEST(BasicPool, StringIndexing) {
    Obj v = NewStr("xy");
    CHECK_EQUAL('x', AsInt(At(v, 0)));
    CHECK_EQUAL('y', AsInt(At(v, 1)));
    CHECK_EQUAL('x', AsInt(At(v, -2)));
    CHECK_EQUAL('y', AsInt(At(v, -1)));
}

TEST(BasicPool, AppendIntToStr) {
    Obj v = NewStr("abc");
    Append(v, NewInt('d'));
    CHECK_EQUAL(4, Size(v));
    STRCMP_EQUAL("abcd", AsStr(v, tmpBuf, sizeof tmpBuf));
}

TEST(BasicPool, AppendStrToStr) {
    Obj v = NewStr("abc");
    Append(v, NewStr("de"));
    CHECK_EQUAL(5, Size(v));
    STRCMP_EQUAL("abcde", AsStr(v, tmpBuf, sizeof tmpBuf));
    Append(v, NewStr("fghijklmnopqrstuvwxyz"));
    CHECK_EQUAL(26, Size(v));
    STRCMP_EQUAL("abcdefghijklmnopqrstuvwxyz", AsStr(v, tmpBuf, sizeof tmpBuf));
    Append(v, NewStr("_0123456789"));
    CHECK_EQUAL(37, Size(v));
    STRCMP_EQUAL("abcdefghijklmnopqrstuvwxyz_0123456789",
                    AsStr(v, tmpBuf, sizeof tmpBuf));
}

TEST(BasicPool, CreateLongStr) {
    Obj o = NewStr("abcdefghijklmnopqrstuvwxyz_0123456789");
    CHECK_EQUAL(37, Size(o));
    for (int i = 0; i < 20; ++i)
        CHECK_EQUAL('a'+i, AsInt(At(o, i)));
}

TEST(BasicPool, CreateVector) {
    Obj v = NewVec();
    CHECK(IsVec(v));
    CHECK(IsRef(v));
    CHECK_EQUAL(0, Size(v));
}

TEST(BasicPool, RecycleMem) {
    // allocate new objects until there have been two garbage collector runs
    // note: can lead to an infinite loop if the GC is not working properly
    CHECK_EQUAL(0, gcCount);
    while (gcCount < 2)
        CHECK(!IsNil(NewVec()));
}

TEST(BasicPool, ExhaustMem) {
    // append new objects to root vector until we run out of chunks
    // note: can lead to an infinite loop if the GC is not working properly
    CHECK_EQUAL(0, gcCount);
    for (;;) {
        Obj o = NewVec();
        if (IsNil(o))
            break;
        Append(root, o);
    }
    CHECK_EQUAL(1, gcCount);
    // the number of allocations depends greatly on chunk and pool sizes
    int limit = -1; // will always fail
    if (swixSize / CHUNK_SIZE == 100) {
#if CHUNK_SIZE == 8
        limit = 64;
#elif CHUNK_SIZE == 16
        limit = 82;
#elif CHUNK_SIZE == 32
        limit = 89;
#elif CHUNK_SIZE == 64
        limit = 92;
#elif CHUNK_SIZE == 128
        limit = 94;
#endif
    }
    // the rest is for some special objects and the vector's chunk chain
    CHECK_EQUAL(limit, Size(root));
}

TEST_GROUP(BasicVector)
{
    Obj v;

    TEST_SETUP() {
        Init();
        v = NewVec();
    }
};

TEST(BasicVector, Empty) {
    CHECK_EQUAL(0, Size(v));
}

TEST(BasicVector, SmallAppend) {
    Append(v, NewInt(456));
    Append(v, NewStr("def"));
    CHECK_EQUAL(2, Size(v));
}

TEST(BasicVector, SmallAccess) {
    Append(v, NewInt(-123));
    Append(v, NewStr("abcde"));

    CHECK_EQUAL(2, Size(v));
    CHECK(IsInt(At(v, 0)));
    CHECK(IsStr(At(v, 1)));

    CHECK_EQUAL(-123, AsInt(At(v, 0)));
    CHECK_EQUAL(5, Size(At(v, 1)));
    STRCMP_EQUAL("abcde", AsStr(At(v, 1), tmpBuf, sizeof tmpBuf));
}

TEST(BasicVector, AccessFromEnd) {
    Append(v, NewInt(11));
    Append(v, NewInt(22));

    CHECK_EQUAL(11, AsInt(At(v, -2)));
    CHECK_EQUAL(22, AsInt(At(v, -1)));
}

TEST(BasicVector, Drop) {
    Append(v, NewInt(1));
    Append(v, NewInt(2));

    CHECK_EQUAL(2, Size(v));
    Drop(v);
    CHECK_EQUAL(1, Size(v));
    Drop(v);
    CHECK_EQUAL(0, Size(v));
}

TEST(BasicVector, Pack) {
    Append(v, NewInt(1));
    Append(v, NewInt(2));
    Append(v, NewInt(3));
    Append(v, NewInt(4));
    Obj o = Pack(v, 3);

    CHECK_EQUAL(1, Size(v));
    CHECK_EQUAL(3, Size(o));

    CHECK_EQUAL(2, AsInt(At(o, 0)));
    CHECK_EQUAL(3, AsInt(At(o, 1)));
    CHECK_EQUAL(4, AsInt(At(o, 2)));
}

TEST(BasicVector, LongVec) {
    for (int i = 0; i < 20; ++i)
        Append(v, NewInt(2*i));
    CHECK_EQUAL(20, Size(v));
    for (int i = 0; i < 20; ++i)
        CHECK_EQUAL(2*i, AsInt(At(v, i)));
}
