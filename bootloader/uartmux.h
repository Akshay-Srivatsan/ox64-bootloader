#pragma once
#include <stdint.h>
#include <stddef.h>
#include "gpio.h"

enum uartmux_function {
  UARTMUX_RTS,
  UARTMUX_CTS,
  UARTMUX_TX,
  UARTMUX_RX,
};

size_t uartmux_signal_number(gpio_t pin);
void uartmux_configure(size_t signal, size_t uart, enum uartmux_function function);
