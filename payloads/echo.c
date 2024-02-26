#define LOG_LEVEL 3
#include "lib.h"

void kmain(void) {
  uart_init(UART0, 115200);
  while (true) {
    char c = uart_getc(UART0);
    uart_putc(UART0, c);
  }
}
