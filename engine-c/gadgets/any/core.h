#include <jet.h>

#define CORE_NAME_LIST \
    "jet", \
    "init", \
    "pass", \
    "print", \

#define CORE_CONFIG_LIST \
    core_jetConfig, \
    core_initConfig, \
    core_passConfig, \
    core_printConfig, \

#define CORE_DISPATCH_LIST \
    core_jetDispatch, \
    core_initDispatch, \
    core_passDispatch, \
    core_printDispatch, \

extern void core_jetConfig (Config_t*,Value_t);
extern void core_jetDispatch (Gadget_t*,int,Value_t);

extern void core_initConfig (Config_t*,Value_t);
extern void core_initDispatch (Gadget_t*,int,Value_t);

extern void core_passConfig (Config_t*,Value_t);
extern void core_passDispatch (Gadget_t*,int,Value_t);

extern void core_printConfig (Config_t*,Value_t);
extern void core_printDispatch (Gadget_t*,int,Value_t);
