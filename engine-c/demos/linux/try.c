#include <stdio.h>
#include <jet.h>

extern const char* leenux();

// TODO get rid of these dependencies, may need to use a library
const char* jNameTable [] = { 0 };
void (*jConfigTable[])(Config_t*,Value_t) = { };
void (*jDispatchTable[])(Gadget_t*,int,Value_t) = { };

int main () {
    puts(leenux());
    return 0;
}
