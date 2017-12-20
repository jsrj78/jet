// JET engine C header, used by gadgets.

#include <stdint.h>

typedef struct {
    uint16_t f :1;
    int16_t  i :15;
} Value_t;

typedef struct {
    int8_t   inlets;
    int8_t   outlets;
    void*    state;
    uint16_t stateLen;
} Config_t;

typedef struct {
    uint16_t handler :8;
    uint16_t state   :14;
    uint16_t outlets :12;
    uint16_t chunk   :12;
} Gadget_t;

extern const char* jNameTable [];
extern void (*jConfigTable[])(Config_t*,Value_t);
extern void (*jDispatchTable[])(Gadget_t*,int,Value_t);

extern void jPrint (Value_t);
extern void jEmit (Gadget_t*,int,Value_t);
