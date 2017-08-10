// Swix configuration options.

// chunks are 8/16/32/etc bytes, this is the Log2
#define CHUNK_BITS  4

// exactly one of the following reference models must be enabled
#define REF_SMALL   1   // 14-bit references
#define REF_LARGE   0   // 30-bit references
