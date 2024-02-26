#include "timer.h"
#include "memory.h"

static const unsigned cpu_freq = 480 * 1000 * 1000;

static void rmw(volatile u32 *reg, u32 val, u32 mask) {
    uint x = get32(reg);
    x &= ~mask;
    x |= val & mask;
    put32(reg, x);
}

void timer_init(unsigned frequency) {
  u32 *mm_misc_base = (u32 *)0x30000000;
  u32 *mm_misc_cpu_d0 = mm_misc_base + 0x6; // offset = 0x18

  rmw(mm_misc_cpu_d0, 0, 1 << 31); // disable timer
  rmw(mm_misc_cpu_d0, cpu_freq / frequency, 0x3ff); // set timer divider
  rmw(mm_misc_cpu_d0, 1 << 31, 1 << 31); // enable timer
}

u64 timer_read(void) {
  u64 time;
  asm volatile("csrr %0, time":"=r"(time));
  return time;
}
