#define LOG_LEVEL 3
#include "lib.h"

void kmain(void) {
  while (true) {
    char c = uart_getc();
    uart_putc(c);
  }
}
