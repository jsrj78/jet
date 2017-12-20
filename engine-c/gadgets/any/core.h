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

extern void core_jetConfig (JConfig*,JValue);
extern void core_jetDispatch (JGadget*,int,JValue);

extern void core_initConfig (JConfig*,JValue);
extern void core_initDispatch (JGadget*,int,JValue);

extern void core_passConfig (JConfig*,JValue);
extern void core_passDispatch (JGadget*,int,JValue);

extern void core_printConfig (JConfig*,JValue);
extern void core_printDispatch (JGadget*,int,JValue);
