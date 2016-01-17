// Public API definitions.

#include <stdint.h>

#ifdef ARDUINO
#include <stdlib.h> // for size_t
typedef unsigned long long uint64_t;
#endif

#define SWIX_VERSION 1

extern uint64_t swixPool [];
extern size_t swixSize;

extern int Swix (const char* desc, ...);
