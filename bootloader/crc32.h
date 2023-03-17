#pragma once

#include <stdint.h>

uint32_t crc32_inc(const void *buf, unsigned size, uint32_t crc);
uint32_t crc32(const void *buf, unsigned size);
