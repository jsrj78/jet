#include "core.h"

#include <assert.h>
#include <stdio.h>

const char* jNameTable [] = {
    CORE_NAME_LIST
    0
};

void (*jConfigTable[])(Config_t*,Value_t) = {
    CORE_CONFIG_LIST
};

void (*jDispatchTable[])(Gadget_t*,int,Value_t) = {
    CORE_DISPATCH_LIST
};

int main () {
    assert(sizeof (Value_t) == 2);

    int nameLen = sizeof jNameTable / sizeof *jNameTable - 1;
    int configLen = sizeof jConfigTable / sizeof *jConfigTable;
    int dispatchLen = sizeof jDispatchTable / sizeof *jDispatchTable;

    printf("%d names, %d configs, %d handlers:\n",
            nameLen, configLen, dispatchLen);
    for (const char** p = jNameTable; *p != 0; ++p)
        printf("  %s\n", *p);

    assert(nameLen == configLen);
    assert(configLen == dispatchLen);

    return 0;
}
