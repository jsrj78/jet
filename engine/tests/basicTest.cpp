#include "CppUTest/TestHarness.h"

extern "C" {
#include "engine.h"
}

static Gadget *gp1, *gp2;

TEST_GROUP(Printing)
{
    void setup() {
        ResetPrint();
        gp1 = gp2 = 0;
    }
    void teardown() {
        FreeGadget(gp1);
        FreeGadget(gp2);
    }
};

//------------------------------------------------------------------------------
// add new tests in reverse order, below this comment

TEST(Printing, MosesGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w123[] = {
        { 0, 2, 0 },    /* g1.0 -> g2.0 */
        { 1, 3, 0 },    /* g1.1 -> g3.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(1, 0, 4);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, LookupGadget("moses", 5), w123);
    Add(gp1, LookupGadget("print", 1), 0);
    Add(gp1, LookupGadget("print", 2), 0);

    Feed(gp1, 0, 4);
    Feed(gp1, 0, 5);
    Feed(gp1, 0, 6);

    static Message result[] = { 1, 4, 2, 5, 2, 6 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, ChangeGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w12[] = {
        { 0, 2, 0 },    /* g1.0 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(1, 0, 3);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, LookupGadget("change", 0), w12);
    Add(gp1, LookupGadget("print", 0), 0);

    Feed(gp1, 0, 0);
    Feed(gp1, 0, 1);
    Feed(gp1, 0, 1);
    Feed(gp1, 0, 2);
    Feed(gp1, 0, 2);
    Feed(gp1, 0, 3);
    Feed(gp1, 0, 0);

    static Message result[] = { 0, 1, 2, 3, 0 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, SwapGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w123[] = {
        { 0, 2, 0 },    /* g1.0 -> g2.0 */
        { 1, 3, 0 },    /* g1.1 -> g3.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(1, 0, 4);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, LookupGadget("swap", 123), w123);
    Add(gp1, LookupGadget("print", 1), 0);
    Add(gp1, LookupGadget("print", 2), 0);

    Feed(gp1, 0, 20);

    static Message result[] = { 2, 20, 1, 123 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, TwoOutlets) {
    static Wire w012[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 2, 0 },    /* g0.0 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w12x[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 1, 2, 0 },    /* g0.1 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };

    Gadget* gp = NewCircuit(1, 2, 3);
    Add(gp, LookupGadget("inlet", 0), w012);
    Add(gp, LookupGadget("outlet", 0), 0);
    Add(gp, LookupGadget("outlet", 0), 0);

    gp1 = NewCircuit(0, 0, 3);
    Add(gp1, gp, w12x);
    Add(gp1, LookupGadget("print", 1), 0);
    Add(gp1, LookupGadget("print", 2), 0);

    Feed(gp, 0, 19);

    static Message result[] = { 1, 19, 2, 19 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, TwoInlets) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w23[] = {
        { 0, 3, 0 },    /* g2.0 -> g3.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(2, 0, 4);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, LookupGadget("print", 1), 0);
    Add(gp1, LookupGadget("inlet", 0), w23);
    Add(gp1, LookupGadget("print", 2), 0);

    Feed(gp1, 0, 17);
    Feed(gp1, 1, 18);

    static Message result[] = { 1, 17, 2, 18 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, NestedGadgets) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };
    static Wire w12[] = {
        { 0, 2, 0 },    /* g1.0 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };

    Gadget* gp = NewCircuit(1, 1, 2);
    Add(gp, LookupGadget("inlet", 0), w01);
    Add(gp, LookupGadget("outlet", 0), 0);

    gp1 = NewCircuit(0, 0, 3);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, gp, w12);
    Add(gp1, LookupGadget("print", 0), 0);

    Feed(gp1, 0, 16);

    static Message result[] = { 16 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, OutletPrintGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };

    Gadget* gp = NewCircuit(1, 1, 2);
    Add(gp, LookupGadget("inlet", 0), w01);
    Add(gp, LookupGadget("outlet", 0), 0);

    gp1 = NewCircuit(0, 0, 2);
    Add(gp1, gp, w01);
    Add(gp1, LookupGadget("print", 0), 0);

    Feed(gp, 0, 15);

    static Message result[] = { 15 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, InletPrintGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(1, 0, 2);
    Add(gp1, LookupGadget("inlet", 0), w01);
    Add(gp1, LookupGadget("print", 0), 0);

    Feed(gp1, 0, 14);

    static Message result[] = { 14 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PassPrintTwiceGadget) {
    static Wire w012[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 2, 0 },    /* g0.0 -> g2.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(0, 0, 3);
    Gadget* gp;
    Add(gp1, gp = LookupGadget("pass", 0), w012);
    Add(gp1, LookupGadget("print", 1), 0);
    Add(gp1, LookupGadget("print", 2), 0);

    Feed(gp, 0, 13);

    static Message result[] = { 1, 13, 2, 13 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PassAndPrintGadget) {
    static Wire w01[] = {
        { 0, 1, 0 },    /* g0.0 -> g1.0 */
        { 0, 255, 0 },  /* end marker */
    };

    gp1 = NewCircuit(0, 0, 2);
    CHECK_EQUAL(3 * sizeof(Gadget*), gp1->extra);
    Gadget* gp;
    Add(gp1, gp = LookupGadget("pass", 0), w01);
    Add(gp1, LookupGadget("print", 0), 0);

    Feed(gp, 0, 12);

    static Message result[] = { 12 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PassGadgetExists) {
    gp1 = LookupGadget("pass", 0);
    CHECK(gp1 != 0);
}

TEST(Printing, PrintGadgetArg) {
    gp1 = LookupGadget("print", 123);
    Feed(gp1, 0, 11);
    static Message result[] = { 123, 11 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PrintGadget) {
    gp1 = LookupGadget("print", 0);
    Feed(gp1, 0, 10);
    static Message result[] = { 10 };
    MEMCMP_EQUAL(result, g_PrintBuffer, sizeof result);
}

TEST(Printing, PrintGadgetExists) {
    gp1 = LookupGadget("print", 0);
    CHECK(gp1 != 0);
}

TEST_GROUP(Basic) {};

TEST(Basic, LookupNonexistent) {
    CHECK_EQUAL(0, LookupGadget("blah", 0));
}
