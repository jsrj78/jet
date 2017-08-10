// type definitions

typedef int Message;

typedef struct Wire_t {
    uint8_t out :4; /* source outlet */
    uint8_t gid;    /* gadget ID of destination */
    uint8_t in :4;  /* destination inlet */
} Wire;

typedef struct Gadget_t {
    uint8_t inlets, outlets;
    uint16_t extra;
    void (*handler)(struct Gadget_t*,int,Message);
    struct Circuit_t *parent;
    const Wire* wires;
    Message arg; // TODO should treat this as extra data
    // extra data stored here
} Gadget;

typedef struct Circuit_t {
    Gadget _;
    Gadget* child[1]; // entries extend past end
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
extern Gadget* NewGadget (uint8_t i, uint8_t o, uint16_t x,
                          void (*h)(Gadget*,int,Message));
extern Circuit* NewCircuit (uint8_t i, uint8_t o, uint8_t g);
extern void Add (Circuit* cp, int pos, Gadget* gp);
extern void AddWires (Gadget* gp, const Wire* w);
extern void Feed (Gadget* gp, int inlet, Message msg);
extern void Emit (Gadget* gp, int outlet, Message msg);

// fake print output, saved in a buffer instead

#define NMSGS 10
extern Message g_PrintBuffer[];
extern void ResetPrint (void);
