#include <stdio.h>
#include <jet.h>

extern const char* macos();

int main () {
    puts(jet());
    puts(macos());
    return 0;
}
