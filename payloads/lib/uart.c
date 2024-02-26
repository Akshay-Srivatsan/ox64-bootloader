#include <stdbool.h>

#include "memory.h"
#include "types.h"
#include "uart.h"

// from the bl808 reference manual
struct uart {
  u32 tx_config;
  u32 rx_config;
  u32 bit_prd;
  u32 data_config;
  u32 tx_ir_position;
  u32 rx_ir_position;
  u32 rx_rto_timer;
  u32 sw_mode;
  u32 int_sts;
  u32 int_mask;
  u32 int_clear;
  u32 int_en;
  u32 status;
  u32 sts_urx_abr_prd;
  u32 rx_abr_prd_b01;
  u32 rx_abr_prd_b23;
  u32 rx_abr_prd_b45;
  u32 rx_abr_prd_b67;
  u32 rx_abr_pw_tol;
  u32 rx_bcr_int_cfg;
  u32 unknown;
  u32 tx_rs485_cfg;
  u32 reserved[10];
  u32 fifo_config_0;
  u32 fifo_config_1;
  u32 fifo_wdata;
  u32 fifo_rdata;
};

volatile struct uart *const UART0 = (volatile struct uart *)0x2000A000;
volatile struct uart *const UART1 = (volatile struct uart *)0x2000A100;
volatile struct uart *const UART2 = (volatile struct uart *)0x2000AA00;

#define UART_CLOCK 40000000UL

bool uart_can_getc(volatile struct uart *uart) {
  return ((get32(&uart->fifo_config_1) >> 8) & 0x3f) != 0;
}

char uart_getc(volatile struct uart *uart) {
  while (!uart_can_getc(uart))
    ;
  return get32(&uart->fifo_rdata);
}

bool uart_can_putc(volatile struct uart *uart) {
  return (get32(&uart->fifo_config_1) & 0x3f) != 0;
}

void uart_putc(volatile struct uart *uart, char c) {
  while (!uart_can_putc(uart))
    ;
  put32(&uart->fifo_wdata, c);
}

void uart_puts(volatile struct uart *uart, const char *c) {
  while (*c) {
    uart_putc(uart, *c++);
  }
}

void uart_init(volatile struct uart *uart, unsigned baud) {
  // p302
  uint32_t tx = 0;
  tx |= 1 << 0;  // enable
  tx |= 1 << 2;  // freerun mode
  tx |= 7 << 8;  // 8 data bits
  tx |= 2 << 11; // 1 stop bit
  put32(&uart->tx_config, tx);

  uint32_t rx = 0;
  rx |= 1 << 0; // enable
  rx |= 7 << 8; // 8 data bits
  put32(&uart->rx_config, rx);

  uint32_t bit_period_value =
      (2 * UART_CLOCK / baud) - 1; // for some reason the uart clock seems to be
                                   // off by a factor of 2 from what we expect
  uint32_t val = bit_period_value << 16 | bit_period_value;
  put32(&uart->bit_prd, val);
}
