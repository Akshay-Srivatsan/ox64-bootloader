#define LOG_LEVEL 3
#include "lib.h"

void kmain(void) {
  while (true) {
    uart_puts("hello, world!\n");
    unsigned long start = timer_read();
    while (timer_read() - start < 16000000 / 10) {
    }
  }
}
