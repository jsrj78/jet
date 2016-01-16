// Garbage collector tests.

extern "C" {
#include "defs.h"
}

#include "CppUTest/TestHarness.h"
#include <string.h>

static char tmpBuf[50];

TEST_GROUP(Garbage)
{
    Obj root;
    const char* longStr;

    TEST_SETUP() {
        root = Init();
        longStr = "abcdefghijklmnopqrstuvwxyz_0123456789";
    }

    void allocUntilFull (int len) {
        for (;;) {
            Obj o = NewStrN(longStr, (size_t) len);
            if (IsNil(o))
                break;
            Append(root, o);
        }
    }

    int countValid (int len) {
        int ok = 0;
        int n = Size(root);
        for (int i = 0; i < n; ++i) {
            Obj o = At(root, i);
            if (IsStr(o) && Size(o) == len) {
                const char* s = AsStr(o, tmpBuf, sizeof tmpBuf);
                if (memcmp(longStr, s, (size_t) len) == 0)
                    ++ok;
            }
        }
        return ok;
    }
};

#define MAX_PAYLOAD (CHUNK_SIZE - 2 * sizeof (Obj))

static int expectedAllocs (size_t bytes) {
    int numChunks = (int) ((bytes + MAX_PAYLOAD - 1) / MAX_PAYLOAD);
    if (swixSize / CHUNK_SIZE != 100)
        return -1;
    static const int16_t limits[][5] = {
        // TODO: negative entries below have not yet been determined
        // 8, 16, 32, 64, 128 = CHUNK_SIZE
        {192,576,1344,2880,-7},     // no chunks
        { 64, 82, 89, 92, 94 },     // 1 chunk
        { 38, 44, 46,-26,-27 },     // 2 chunks
        { 27, 30,-35,-36,-37 },     // 3 chunks
        { 21, 23,-45,-46,-47 },     // 4 chunks
        {-53,-54,-55,-56,-57 },     // 5 chunks
        {-63,-64,-65,-66,-67 },     // 6 chunks
        { 12,-74,-75,-76,-77 },     // 7 chunks
        {-83,-84,-85,-86,-87 },     // 8 chunks
        {-93,-94,-95,-96,-97 },     // 9 chunks
        {  9, -4, -5, -6, -7 },     // 10 chunks
    };
    return limits[numChunks][CHUNK_BITS-3];
}

TEST(Garbage, ExhaustNil) {
    CHECK_EQUAL(0, gcCount);
    while (!IsNil(Append(root, NilVal())))
        ;
    CHECK_EQUAL(0, gcCount);
    int limit = expectedAllocs(0);
    CHECK_EQUAL(limit, Size(root));
}

TEST(Garbage, ExhaustStr1) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(1);
    int limit = expectedAllocs(2);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(1));
}

TEST(Garbage, ExhaustStr4) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(4);
    int limit = expectedAllocs(5);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(4));
}

TEST(Garbage, ExhaustStr11) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(11);
    int limit = expectedAllocs(12);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(11));
}

TEST(Garbage, ExhaustStr12) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(12);
    int limit = expectedAllocs(13);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(12));
}

TEST(Garbage, ExhaustStr13) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(13);
    int limit = expectedAllocs(14);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(13));
}

TEST(Garbage, ExhaustStr26) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(26);
    int limit = expectedAllocs(27);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(26));
}

TEST(Garbage, ExhaustStr37) {
    CHECK_EQUAL(0, gcCount);
    allocUntilFull(37);
    int limit = expectedAllocs(38);
    CHECK_EQUAL(limit, Size(root));
    CHECK_EQUAL(limit, countValid(37));
}
