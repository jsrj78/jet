// Swix internal definitions.

#include "swix.h"
#include "config.h"

#define CHUNK_SIZE  (1 << CHUNK_BITS)
#define REF_SIZE    (REF_SMALL ? 2 : 4)

typedef struct {
#if REF_SMALL
    int16_t _;
#else
    int32_t _;
#endif
} Obj;

// private definitions
extern uint8_t gcCount;

Obj specialObj (int n);
int boxedType (Obj o);
Obj boxedNewInt (int n);
int boxedAsInt (Obj o);

// in mem.c
Obj Init (void);
Obj NewInt (int n);
Obj NewVec (void);
Obj NewStr (const char* s);
Obj NewStrN (const char* s, size_t n);
int Size (Obj o);
const char* AsStr (Obj o, char* buf, size_t len);
Obj Append (Obj v, Obj o);
Obj At (Obj v, int n);
void Drop (Obj v);
Obj Pack (Obj v, int n);

// in json.c
typedef struct {
    uint8_t state;
    uint8_t level;
    uint8_t minus;
    int8_t decimals;
    int32_t value;
    int16_t exponent;
    Obj result;
} JsonState;

void JsonEmit (Obj o, int(*)(char,void*), void*);
void JsonInit (JsonState* j);
char JsonFeed (JsonState* j, int c);
Obj JsonDone (JsonState* j);

#define NilVal() specialObj(0)
//inline Obj NilVal (void) {
//    return specialObj(0);
//}

#define IsNil(o) ((o)._ == 0)
//inline int IsNil (Obj o) {
//    return o._ == 0;
//}

#define TagVal() specialObj(1)
//inline Obj TagVal (void) {
//    return specialObj(1);
//}

#define IsTag(o) ((o)._ == 2)
//inline int IsTag (Obj o) {
//    return o._ == 2;
//}

#define BoolVal(v) specialObj((v)?3:2)
//inline Obj BoolVal (int v) {
//    return specialObj(v ? 3 : 2);
//}

#define IsBool(o) ((o)._ == 4 || (o)._ == 6)
//inline int IsBool (Obj o) {
//    return o._ == 4 || o._ == 6;
//}

#define IsRef(o) (((o)._ & 1) == 0)
//inline int IsRef (Obj o) {
//    return (o._ & 1) == 0;
//}

#define IsInt(o) (((o)._ & 3) == 1 || boxedType(o) == 0)
//inline int IsInt (Obj o) {
//    return (o._ & 3) == 1 || boxedType(o) == 0;
//}

#define IsStr(o) (boxedType(o) == 1)
//inline int IsStr (Obj o) {
//    return boxedType(o) == 1;
//}

#define IsVec(o) (boxedType(o) == 2)
//inline int IsVec (Obj o) {
//    return boxedType(o) == 2;
//}

#define AsInt(o) (((o)._ & 3) == 1 ? (int) ((o)._ >> 2) : boxedAsInt(o))
//inline int AsInt (Obj o) {
//    return (o._ & 3) == 1 ? (int) (o._ >> 2) : boxedAsInt(o);
//}
