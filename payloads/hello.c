#define LOG_LEVEL 3
#include "lib.h"

void kmain(void) {
  uart_init(UART0, 115200);
  while (true) {
    uart_puts(UART0, "hello, world!\n");
    unsigned long start = timer_read();
    delay_ms(1000);
  }
}
