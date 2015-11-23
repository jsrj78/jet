// Tight Data structures, public interface
#pragma once

typedef struct { int16_t _; } Td_Val;

extern    void tdInitPool (void);
extern int16_t tdChain (void);
extern int16_t tdAlloc (void);
extern    void tdDelRef (Td_Val v);
extern  Td_Val tdNewInt (int32_t v);
extern int32_t tdAsInt (Td_Val v);
extern  Td_Val tdNewStr (const char* s);
extern int16_t tdSize (Td_Val v);
