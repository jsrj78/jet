// JET engine C header, used by gadgets.

#include <stdint.h>

typedef struct {
    uint16_t ref :1;
    int16_t  val :15;
} JValue;

typedef struct {
} JGadget;

typedef struct {
    int inlets;
    int outlets;
} JConfig;

extern const char* jet (); // TODO remove
extern void jEmit (JGadget*,int,JValue);
