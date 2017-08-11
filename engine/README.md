# JET/Engine

* this flow engine is implemented in C
* it uses [CppUTest](https://cpputest.github.io) for testing, tests are in C++
* only handles ints, needs major extension to support richer message types

### Running the tests

    cd jet/engine
    make

### Sample output

```
$ make
compiling basicTest.cpp
compiling main.cpp
compiling engine.c
compiling gadgets.c
Building archive lib/libengine.a
a - objs/src/engine.o
a - objs/src/gadgets.o
Linking engine_tests
Running engine_tests
...............
OK (15 tests, 15 ran, 16 checks, 0 ignored, 0 filtered out, 0 ms)
```
