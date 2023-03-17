#include "string.h"

void *memset(void *s, int c, size_t n) {
  for (size_t i = 0; i < n; i++) {
    ((char *)s)[i] = c;
  }
  return s;
}

void *memcpy(void *restrict dest, const void *restrict src, size_t n) {
  for (size_t i = 0; i < n; i++) {
    ((char *)dest)[i] = ((char *)src)[i];
  }
}
