#define LOG_LEVEL 5

#include "lib.h"

struct PACKED payload {
  uint64_t entry;
  uint32_t size;
  uint32_t cksum;
  uint8_t data[];
};

extern const struct payload payload;

void kmain() {
  {
    TRACE("uart init\n");
    uart_init(2000000);

    gpio_set_function(14, GPIO_UART);
    gpio_set_function(15, GPIO_UART);

    uartmux_configure(uartmux_signal_number(14), 0, UARTMUX_TX);
    uartmux_configure(uartmux_signal_number(15), 0, UARTMUX_RX);
  }
  INFO("bootloader init\n");
  {
    TRACE("psram init\n");
    psram_init();
    TRACE("psram check\n");
    volatile u32 *psram = (volatile u32 *)0x50000000;
    put32(psram, 0xdeadbeef);
    u32 x = get32(psram);
    DEBUG("read back 0x%x from psram\n", x);
    if (x != 0xdeadbeef) {
      ERROR("psram init failed!\n");
      while(1);
    }
    INFO("psram init success\n");
  }
  TRACE("timer init\n");
  timer_init(1000000);

  INFO("addr = 0x%lx\n", payload.entry);
  INFO("size = 0x%x\n", payload.size);
  INFO("cksum = 0x%x\n", payload.cksum);

  DEBUG("copying %d bytes of code to address %p\n", payload.size, payload.entry);
  memcpy((void *)payload.entry, payload.data, payload.size);

  uint32_t cksum = crc32((void *)payload.entry, payload.size);
  if (payload.cksum != cksum) {
    ERROR("invalid checksum\n");
    ERROR("refusing to boot\n");
    while(1);
  } else {
    INFO("valid checksum: 0x%x\n", cksum);
  }

  TRACE("branching to code\n");
  void (*fn)(void) = (void (*)(void))payload.entry;
  fn();
  WARN("code returned to bootloader\n");
}
