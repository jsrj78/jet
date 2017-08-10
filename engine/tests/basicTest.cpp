#include "CppUTest/TestHarness.h"

extern "C" {
#include "engine.h"
}

static Gadget *gp1, *gp2, *cp1;

TEST_GROUP(Basic) {};

TEST(Basic, LookupNonexistent) {
    CHECK_EQUAL(0, LookupGadget("blah", 0));
}

TEST_GROUP(Printing)
{
    void setup() {
        ResetPrint();
        gp1 = gp2 = 0;
        cp1 = 0;
    }
    void teardown() {
        FreeGadget(gp1);
        FreeGadget(gp2);
        FreeGadget(cp1);
    }
};

TEST(Printing, PrintGadgetExists) {
    gp1 = LookupGadget("print", 0);
    CHECK(gp1 != 0);
}

TEST(Printing, PrintGadget) {
    gp1 = LookupGadget("print", 0);
    Feed(gp1, 0, 10);
    static Message result[] = { 10 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PrintGadgetArg) {
    gp1 = LookupGadget("print", 123);
    Feed(gp1, 0, 11);
    static Message result[] = { 123, 11 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PassGadgetExists) {
    gp1 = LookupGadget("pass", 0);
    CHECK(gp1 != 0);
}

TEST(Printing, PassAndPrintGadget) {
    cp1 = NewCircuit(0, 0, 2);
    CHECK_EQUAL(3 * sizeof(Gadget*), cp1->extra);
    static Wire w0[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    Gadget* g;
    Add(cp1, g = LookupGadget("pass", 0), w0);
    Add(cp1, LookupGadget("print", 0), 0);

    Feed(g, 0, 12);

    static Message result[] = { 12 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PassPrintTwiceGadget) {
    cp1 = NewCircuit(0, 0, 3);
    static Wire w0[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 2, 0 },    /* g0.0 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };
    Gadget* g;
    Add(cp1, g = LookupGadget("pass", 0), w0);
    Add(cp1, LookupGadget("print", 1), 0);
    Add(cp1, LookupGadget("print", 2), 0);

    Feed(g, 0, 13);

    static Message result[] = { 1, 13, 2, 13 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}
