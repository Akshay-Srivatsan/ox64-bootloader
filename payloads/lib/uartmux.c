#include "uartmux.h"
#include "memory.h"

static volatile uint32_t *const uartmux_signal =
    (volatile uint32_t *)0x20000154;

size_t uartmux_signal_number(gpio_t pin) { return pin % 12; }

void uartmux_configure(size_t signal, size_t uart,
                       enum uartmux_function function) {
  unsigned fn = (uart * 4) + function;
  volatile uint32_t *const reg = uartmux_signal + (signal / 8);
  uint32_t value = get32(reg);
  size_t shift = (signal % 8) * 4;
  value &= ~(0b1111 << shift);
  value |= (fn & 0b1111) << shift;
  put32(reg, value);
}
