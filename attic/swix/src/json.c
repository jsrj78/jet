// Swix JSON input and output conversion.

#include "defs.h"
#include <string.h>

static void strOut (const char* s, int(*f)(char,void*), void* p) {
    while (*s)
        f(*s++, p);
}

static void intOut (uint32_t v, int(*f)(char,void*), void* p) {
    if (v > 10)
        intOut(v/10, f, p);
    f((char)('0' + v%10), p);
}

void JsonEmit(Obj o, int(*f)(char,void*), void* p) {
    if (IsNil(o)) {
        strOut("null", f, p);
        return;
    }
    if (IsBool(o)) {
        strOut(AsInt(o) ? "true" : "false", f, p);
        return;
    }
    if (IsInt(o)) {
        int v = AsInt(o);
        if (v < 0) {
            v = -v;
            f('-', p);
        }
        intOut((uint32_t) v, f, p);
        return;
    }
    if (IsVec(o)) {
        int n = Size(o);
        int map = n >= 2 && IsTag(At(o, 0)) && IsNil(At(o, 1));
        char sep = map ? '{' : '[';
        if (n == 2*map)
            f(sep, p);
        int i;
        for (i = 2*map; i < n; ++i) {
            f(sep, p);
            JsonEmit(At(o, i), f, p);
            sep = map && i % 2 == 0 ? ':' : ',';
        }
        f(map ? '}' : ']', p);
    }
    if (IsStr(o)) {
        f('"', p);
        int i, n = Size(o);
        for (i = 0; i < n; ++i) {
            char esc = 0, c = (char) AsInt(At(o, i));
            switch (c) {
                case '\n':  esc = 'n'; break;
                case '\r':  esc = 'r'; break;
                case '\t':  esc = 't'; break;
                case '"':   esc = '"'; break;
                case '\\':  esc = c; break;
                default:    if ((uint8_t) c <= 0x1F) {
                                f('\\', p);
                                f('u', p);
                                f('0', p);
                                f('0', p);
                                f((char) ('0'+(c>>4)), p);
                                c = "0123456789abcdef"[c&0xF];
                            }
            }
            if (esc)  {
                c = esc;
                f('\\', p);
            }
            f(c, p);
        }
        f('"', p);
    }
}

enum { Sini, Snum, Sexp, Skey, Sstr, Sesc, Suni };

void JsonInit (JsonState* j) {
    memset(j, 0, sizeof *j);
    j->result = NewVec();
}

char JsonFeed (JsonState* j, int c) {
    char r = 0;
    switch (j->state++) {
        case Sini:
       ini: j->state = Sini;
            j->value = 0;
            j->decimals = -1;
            j->exponent = 0;
            if ('a' <= c && c <= 'z') {
                j->value = c;
                j->state = Skey;
                break;
            }
            if ('0' <= c && c <= '9') {
                j->value = c - '0';
                j->minus = 0;
                j->state = Snum;
                break;
            }
            if (c == '-') {
                j->minus = 1;
                j->state = Snum;
            } else if (c == '"') {
                Append(j->result, NewStr(""));
                j->state = Sstr;
            } else if (c == '[' || c == '{') {
                Append(j->result, NewInt(j->level));
                j->level = (uint8_t) Size(j->result);
                if (c == '{') {
                    Append(j->result, TagVal());
                    Append(j->result, NilVal());
                }
            } else if (c == ']' || c == '}') {
                Obj o = Pack(j->result, Size(j->result) - j->level);
                j->level = (uint8_t) AsInt(At(j->result, -1));
                Drop(j->result);
                Append(j->result, o);
                r = c == ']' ? 'v' : 'm';
            } else if (c == ',' || c == ':') {
                r = 0;
            } else if (r == 0 && c < ' ')
                return '?';
            break;
        case Snum:
            j->state = Snum;
            if ('0' <= c && c <= '9') {
                if ((uint32_t) j->value >= (1U<<31)/10)
                    break; // drop decimals when precision is exceeded
                if (j->decimals >= 0)
                    ++j->decimals; // FIXME: wrong if j->decimals < 0
                j->value = 10 * j->value + (c - '0');
                break;
            }
            if (c == '.') {
                j->decimals = 0;
                break;
            }
            if (j->minus)
                j->value = -j->value;
            j->minus = 0;
            if (c == 'e' || c == 'E') {
                j->state = Sexp;
                break;
            }
            // fall through
        case Sexp:
            j->state = Sexp;
            if ('0' <= c && c <= '9') {
                j->exponent = (int16_t) (10 * j->exponent + (c - '0'));
                break;
            }
            if (c == '-') {
                j->minus = !j->minus;
                break;
            }
            if (j->minus)
                j->exponent = (int16_t) -j->exponent;
            if (j->decimals > 0)
                j->exponent = (int16_t) (j->exponent - j->decimals);
            Append(j->result, NewInt(j->value));
            r = 'n';
            goto ini;
        case Skey:
            if ('a' <= c && c <= 'z') {
                j->state = Skey;
                break;
            }
            switch (j->value) {
                case 'n': Append(j->result, NilVal()); break;
                case 'f': Append(j->result, BoolVal(0)); break;
                case 't': Append(j->result, BoolVal(1)); break;
                default:  ;
            }
            r = 'k';
            goto ini;
        case Sstr:
            if (c == '\\')
                break;
            if (c == '"') {
                r = 's';
                j->state = Sini;
                break;
            }
            Append(At(j->result, -1), NewInt(c));
            j->state = Sstr;
            break;
        case Sesc:
            if (c == 'u') {
                j->value = 0;
                j->decimals = 0;
                break;
            }
            switch (c) {
                case 'n': c = '\n'; break;
                case 'r': c = '\r'; break;
                case 't': c = '\t'; break;
                default:  ;
            }
            Append(At(j->result, -1), NewInt(c));
            j->state = Sstr;
            break;
        case Suni:
            if (c > '9')
                c = (uint8_t) (c - 7);
            j->value = (j->value << 4) + (c & 0x0F);
            if (++j->decimals < 4) {
                j->state = Suni;
                break;
            }
            Append(At(j->result, -1), NewInt((uint8_t) j->value));
            j->state = Sstr;
            break;
        default:
            return '?';
    }
    return j->level ? 0 : r;
}

Obj JsonDone (JsonState* j) {
    Obj r = j->result;
    JsonInit(j);
    return Size(r) > 0 ? At(r, -1) : NilVal();
}
