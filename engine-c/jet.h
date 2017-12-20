// JET engine C header, used by gadgets.

#include <stdint.h>

typedef struct {
    uint16_t f :1;
    int16_t  i :15;
} Value_t;

typedef struct {
    uint16_t v;
} Iolet_t;

typedef struct {
} Gadget_t;

typedef struct {
    int inlets;
    int outlets;
} Config_t;

extern const char* jNameTable [];
extern void (*jConfigTable[])(Config_t*,Value_t);
extern void (*jDispatchTable[])(Gadget_t*,int,Value_t);

extern void jPrint (Value_t);
extern void jEmit (Gadget_t*,int,Value_t);
