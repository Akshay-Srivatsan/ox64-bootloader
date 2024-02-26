#pragma once
#include <stdbool.h>
#include <stdint.h>

struct uart;
extern volatile struct uart *const UART0;
extern volatile struct uart *const UART1;
extern volatile struct uart *const UART2;

bool uart_can_getc(volatile struct uart *uart);
char uart_getc(volatile struct uart *uart);
bool uart_can_putc(volatile struct uart *uart);
void uart_putc(volatile struct uart *uart, char c);
void uart_puts(volatile struct uart *uart, const char *c);
void uart_init(volatile struct uart *uart, unsigned baud);
