#pragma once
#include <stdint.h>

extern void put32(volatile uint32_t *addr, uint32_t value);
extern uint32_t get32(volatile uint32_t *addr);
