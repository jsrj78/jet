// type definitions

typedef int Message;

typedef struct Wire_t {
    uint8_t out :4; /* source outlet */
    uint8_t gid;    /* gadget ID of destination */
    uint8_t in :4;  /* destination inlet */
} Wire;

typedef struct Gadget_t Gadget;

struct Gadget_t {
    uint16_t extra;
    void (*handler)(Gadget*,int,Message);
    //void (*onAdded)(Gadget*);
    void (*onFree)(Gadget*);
    Gadget *parent;
    const Wire* wires;
    // extra data stored here
};

typedef Gadget* (*Constructor)(Message msg);

struct Lookup_t {
    const char* s;
    Constructor c;
};

// data structures

extern struct Lookup_t g_Gadgets[];

// public API

extern Gadget* LookupGadget (const char *name, Message msg);
extern Gadget* NewGadget (size_t x, void (*h)(Gadget*,int,Message));
extern void* ExtraData(Gadget *cp);
extern void FreeGadget (Gadget* gp);

extern Gadget* NewCircuit (uint8_t g);
extern void Add (Gadget* cp, Gadget* gp, const Wire* w);

extern void Feed (Gadget* gp, int inlet, Message msg);
extern void Emit (Gadget* gp, int outlet, Message msg);

// fake print output, saved in a buffer instead

#define NMSGS 10
extern Message g_PrintBuffer[];
extern void ResetPrint (void);
