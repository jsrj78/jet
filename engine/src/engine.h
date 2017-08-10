// type definitions

typedef int Message;

typedef struct Wire_t {
    uint8_t gnum, gin;
} Wire;

typedef struct Gadget_t {
    uint8_t inlets, outlets;
    void (*handler)(struct Gadget_t*,int,Message);
    struct Circuit_t *parent;
    Message arg;
    // Wire*[] array stored here, one per outlet
} Gadget;

typedef struct Circuit_t {
    Gadget _;
    // Gadget*[] array stored here
} Circuit;

typedef Gadget* (*Constructor)(Message msg);

typedef struct Lookup_t {
    const char* s;
    Constructor c;
} Lookup;

// data structures

extern Lookup g_Gadgets[];

// public API

extern Gadget* LookupGadget (const char *name, Message msg);
extern Gadget* NewGadget (uint8_t i, uint8_t o, size_t x,
                          void (*h)(Gadget*,int,Message));
extern Circuit* NewCircuit (uint8_t i, uint8_t o, uint8_t g);
extern void Add (Circuit* cp, int pos, Gadget* gp);
extern void Feed (Gadget* gp, int inlet, Message msg);
extern void Emit (Gadget* gp, int outlet, Message msg);

// inlined API

static inline const Wire* NthOutletWires (Gadget* gp, int num)
    { return ((const Wire**)(gp + 1))[num]; }
static inline void* GadgetEnd (Gadget* gp)
    { return (const Wire**)(gp + 1) + gp->outlets + 1; }

static inline Gadget* NthGadget (Circuit* cp, int num)
    { return ((Gadget**)(GadgetEnd(&cp->_)))[num]; }
static inline void SetNthGadget (Circuit* cp, int num, Gadget* gp)
    { ((Gadget**)(GadgetEnd(&cp->_)))[num] = gp; }

// fake print output, saved in a buffer instead

#define NMSGS 10
extern Message g_PrintBuffer[];
extern void ResetPrint (void);
