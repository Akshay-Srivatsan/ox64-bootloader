#define LOG_LEVEL 5

#include "lib.h"

struct PACKED payload {
  uint64_t entry;
  uint32_t size;
  uint32_t cksum;
  uint8_t data[];
};

extern const struct payload payload;

#define CLOCK 1000000
#define CKSUM 1

void kmain() {
  {
    TRACE("uart init\n");
    uart_init(115200);

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
  timer_init(CLOCK);

  INFO("addr = 0x%lx\n", payload.entry);
  INFO("size = 0x%x\n", payload.size);
  INFO("cksum = 0x%x\n", payload.cksum);

  DEBUG("copying %d bytes of code to address %p\n", payload.size, payload.entry);
  u64 start = timer_read();
  memcpy((void *)payload.entry, payload.data, payload.size);
  u64 end = timer_read();
  u64 ms = (end - start) * 1000 / CLOCK;
  DEBUG("memcpy took %d ms\n", ms);

#if CKSUM
  start = timer_read();
  uint32_t cksum = crc32((void *)payload.entry, payload.size);
  end = timer_read();
  ms = (end - start) * 1000 / CLOCK;
  DEBUG("crc32 took %d ms\n", ms);
  if (payload.cksum != cksum) {
    ERROR("invalid checksum: got 0x%x, expected 0x%x\n", cksum, payload.cksum);
    ERROR("refusing to boot\n");
    while(1);
  } else {
    INFO("valid checksum: 0x%x\n", cksum);
  }
#else
#warning "Checksum verification disabled."
  WARN("not verifying checksum\n");
#endif

  TRACE("branching to code\n");
  void (*fn)(void) = (void (*)(void))payload.entry;
  fn();
  WARN("code returned to bootloader\n");
}
