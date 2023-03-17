#include "uart.h"

void _cstart(void) {
  extern char _kdata_start[], _kdata_start_load[], _kdata_end[];
  extern char _kbss_start[], _kbss_end[];

  for (char *s = _kdata_start, *d = _kdata_start_load; s < _kdata_end; s++, d++) {
    *s = *d;
  }

  for (char *s = _kbss_start; s < _kbss_end; s++) {
    *s = 0;
  }

  void kmain(void);
  kmain();
}
