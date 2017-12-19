#include <stdio.h>
#include <jet.h>

extern const char* linux();

int main () {
    puts(jet());
    puts(linux());
    return 0;
}
