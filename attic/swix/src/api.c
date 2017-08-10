// Public API implementation.

#include "swix.h"
#include "defs.h"
#include <stdarg.h>

static Obj stack;

int Swix (const char* desc, ...) {
    if (stack._ == 0)
        stack = Init();

    va_list ap;
    va_start(ap, desc);

    int arg = 0;
    for (;;) {
        switch (*desc++) {
            case '^': // Init
                stack = Init();
                arg = SWIX_VERSION;
                continue;

            case '<': // push arg
                Append(stack, NewInt(arg));
            case ' ': // clear arg
                break;
            case '|': // pack arg
                Append(stack, Pack(stack, arg));
                break;
            case '\'': // inline string of size arg
                Append(stack, NewStrN(desc, (size_t) arg));
                desc += arg;
                break;

            case 'i': // push int
                Append(stack, NewInt(va_arg(ap, int)));
                break;
            case 's': // push string
                Append(stack, NewStr(va_arg(ap, const char*)));
                break;
            case 'b': // push arg bytes
                Append(stack, NewStrN(va_arg(ap, const char*), (size_t) arg));
                break;

            case '*': // int arg from list
                arg = va_arg(ap, int);
                continue;
            case '-': // minus arg
                arg = -arg;
                continue;
            case '#': // stack size to arg
                arg = Size(stack);
                continue;
            case '>': // pop to arg
                arg = AsInt(At(stack, -1));
            case '.': // drop
                Drop(stack);
                continue;

            case 0: // done, return current arg
                return arg;

            default: {
                char c = desc[-1];
                if ('0' <= c && c <= '9') {
                    arg = 10 * arg + (c - '0');
                    continue;
                }
                return c; // return unrecognised descriptor
            }
        }
        arg = 0;
    }
}
