#include "CppUTest/CommandLineTestRunner.h"

#include "defs.h"

uint64_t swixPool [100*CHUNK_SIZE/8];
size_t swixSize = sizeof swixPool;

int main (int ac, char** av) {
    return RUN_ALL_TESTS(ac, av);
}
