#define LOG_LEVEL 3
#include "lib.h"

void kmain(void) {
  uart_init(UART0, 115200);
  uart_init(UART1, 115200);
  uart_init(UART2, 115200);

  gpio_set_function(14, GPIO_UART);
  gpio_set_function(15, GPIO_UART);
  gpio_set_output(14, true);
  gpio_set_input(15, true);
  uartmux_configure(uartmux_signal_number(14), 0, UARTMUX_TX);
  uartmux_configure(uartmux_signal_number(15), 0, UARTMUX_RX);

  gpio_set_function(32, GPIO_UART);
  gpio_set_function(33, GPIO_UART);
  gpio_set_output(32, true);
  gpio_set_input(33, true);
  uartmux_configure(uartmux_signal_number(32), 1, UARTMUX_TX);
  uartmux_configure(uartmux_signal_number(33), 1, UARTMUX_RX);

  while (true) {
    uart_puts(UART0, "hello, UART0!\n");
    uart_puts(UART1, "hello, UART1!\n");
    unsigned long start = timer_read();
    delay_ms(1000);
  }
}
