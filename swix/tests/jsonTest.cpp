// JSON input and output tests.

extern "C" {
#include "defs.h"
}

#include "CppUTest/TestHarness.h"

static char buffer [100], *fill;

static int outc (char c, void*) {
    if (buffer <= fill && fill <= buffer + sizeof buffer - 2) {
        *fill++ = c;
        *fill = 0;
    }
    return 0;
}

static int outs (const char* s, void* arg) {
    while (*s)
        outc(*s++, arg);
    return 0;
}

TEST_GROUP(ToJson)
{
    Obj v;

    TEST_SETUP() {
        Init();
        v = NewVec();
        fill = buffer;
        *fill = 0;
    }
};

TEST(ToJson, BufferOut) {
    outc('a', 0);
    outc('b', 0);
    outc('c', 0);
    STRCMP_EQUAL("abc", buffer);
    outs("def", 0);
    STRCMP_EQUAL("abcdef", buffer);
}

TEST(ToJson, EmitNil) {
    JsonEmit(NilVal(), outc, 0);
    STRCMP_EQUAL("null", buffer);
}

TEST(ToJson, EmitInts) {
    JsonEmit(NewInt(0), outc, 0);
    outc(',', 0);
    JsonEmit(NewInt(123), outc, 0);
    outc(',', 0);
    JsonEmit(NewInt(-123), outc, 0);
    outc(',', 0);
    JsonEmit(NewInt(123456789), outc, 0);
    outc(',', 0);
    JsonEmit(NewInt(-123456789), outc, 0);
    STRCMP_EQUAL("0,123,-123,123456789,-123456789", buffer);
}

TEST(ToJson, EmitEmptyVec) {
    JsonEmit(v, outc, 0);
    STRCMP_EQUAL("[]", buffer);
}

TEST(ToJson, EmitVec) {
    Append(v, NewInt(0));
    Append(v, NewInt(123));
    Append(v, NewInt(-123));
    Append(v, NewInt(123456789));
    Append(v, NewInt(-123456789));
    JsonEmit(v, outc, 0);
    STRCMP_EQUAL("[0,123,-123,123456789,-123456789]", buffer);
}

TEST(ToJson, EmitEmptyStr) {
    JsonEmit(NewStr(""), outc, 0);
    STRCMP_EQUAL("\"\"", buffer);
}

TEST(ToJson, EmitStr) {
    JsonEmit(NewStr("a\nb\rc\td\\e\"f"), outc, 0);
    STRCMP_EQUAL("\"a\\nb\\rc\\td\\\\e\\\"f\"", buffer);
}

TEST(ToJson, EmitCtrlChars) {
    JsonEmit(NewStr("a-\x01-\x10-\x1F-b"), outc, 0);
    STRCMP_EQUAL("\"a-\\u0001-\\u0010-\\u001f-b\"", buffer);
}

TEST(ToJson, EmitNullByte) {
    JsonEmit(NewStrN("a-\x00-b", 5), outc, 0);
    STRCMP_EQUAL("\"a-\\u0000-b\"", buffer);
}

TEST(ToJson, EmitEmptyMap) {
    Obj o = NewVec();
    Append(o, TagVal());
    Append(o, NilVal());
    JsonEmit(o, outc, 0);
    STRCMP_EQUAL("{}", buffer);
}

TEST(ToJson, EmitSingleMap) {
    Obj o = NewVec();
    Append(o, TagVal());
    Append(o, NilVal());
    Append(o, NewStr("abc"));
    Append(o, NewInt(123));
    JsonEmit(o, outc, 0);
    STRCMP_EQUAL("{\"abc\":123}", buffer);
}

TEST(ToJson, EmitDoubleMap) {
    Obj o = NewVec();
    Append(o, TagVal());
    Append(o, NilVal());
    Append(o, NewStr("abc"));
    Append(o, NewInt(123));
    Append(o, NewStr("def"));
    Append(o, NewInt(456));
    JsonEmit(o, outc, 0);
    STRCMP_EQUAL("{\"abc\":123,\"def\":456}", buffer);
}

static Obj parse (const char* str) {
    JsonState j;
    JsonInit(&j);
    int r;
    do {
        r = JsonFeed(&j, *str);
        if (*str)
            ++str;
    } while (r == 0);
    return JsonDone(&j);
}

static const char* dump (Obj o) {
    fill = buffer;
    *fill = 0;
    JsonEmit(o, outc, 0);
    return buffer;
}

TEST_GROUP(FromJson)
{
    TEST_SETUP() {
        Init();
    }
};

TEST(FromJson, Empty) {
    STRCMP_EQUAL("null", dump(parse("")));
}

TEST(FromJson, RoundTrips) {
    static const char* tests[] = {
        "null",
        "false",
        "true",
        "123",
        "-123",
        "\"\"",
        "\"abc\"",
        "\"abcdefghijklmnopqrstuvwxyz_0123456789\"",
        "\"<\\\\ \\\" \\n \\r \\t \\u0000 \\u0001 \\u001e \\u001f>\"",
        "\"\u1234\"",
        "[]",
        "[1,2]",
        "[11,[\"abc\",22],[],33]",
        "[1,2,3,4,5]",
        "[1,2,3,4,5,6]",
        "{}",
        "{\"abc\":123}",
        "{\"abc\":123,\"def\":456}",
        0
    };
    for (const char** p = tests; *p != 0; ++p)
        STRCMP_EQUAL(*p, dump(parse(*p)));
}

TEST(FromJson, BadConversions) {
    static const char* tests[] = {
        // first a quick check that these tests work properly
        "123", "123",
        // the following conversions are unintended but unavoidable
        "\"<\\u001F>\"", "\"<\\u001f>\"",
        // any other conversions listed below are wrong and need to be fixed
        0
    };
    for (const char** p = tests; *p != 0; p += 2)
        STRCMP_EQUAL(p[1], dump(parse(p[0])));
}
