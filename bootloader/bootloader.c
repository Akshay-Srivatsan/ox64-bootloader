#define LOG_LEVEL 5

#include "lib.h"

struct PACKED payload {
  uint32_t magic;
  uint64_t entry;
  uint32_t size;
  uint32_t cksum;
  uint8_t data[];
};

extern const struct payload payload;

void kmain() {
  {
    INFO("uart init\n");
    uart_init(2000000);

    gpio_set_function(14, GPIO_UART);
    gpio_set_function(15, GPIO_UART);

    uartmux_configure(uartmux_signal_number(14), 0, UARTMUX_TX);
    uartmux_configure(uartmux_signal_number(15), 0, UARTMUX_RX);
  }
  INFO("bootloader init\n");
  {
    INFO("psram init\n");
    psram_init();
    INFO("psram check\n");
    volatile u32 *psram = (volatile u32 *)0x50000000;
    put32(psram, 0xdeadbeef);
    u32 x = get32(psram);
    INFO("read back 0x%x\n", x);
    if (x != 0xdeadbeef) {
      ERROR("psram init failed!\n");
      while(1);
    }
    INFO("psram init success\n");
  }
  INFO("timer init\n");
  timer_init(1000000);
  INFO("have timer = %d\n", timer_read());

  INFO("magic = 0x%x\n", payload.magic);
  INFO("addr = 0x%lx\n", payload.entry);
  INFO("size = 0x%x\n", payload.size);
  INFO("cksum = 0x%x\n", payload.cksum);

  if (payload.magic != 0x12b9b0a1) {
    ERROR("invalid magic number\n");
    ERROR("refusing to boot\n");
    while(1);
  } else {
    INFO("valid magic: 0x%x == %d\n", payload.magic, payload.magic);
  }
  INFO("copying %d bytes of code to address %p\n", payload.size, payload.entry);
  memcpy((void *)payload.entry, payload.data, payload.size);

  uint32_t cksum = crc32((void *)payload.entry, payload.size);
  if (payload.cksum != cksum) {
    ERROR("invalid checksum\n");
    ERROR("refusing to boot\n");
    while(1);
  } else {
    INFO("valid checksum: 0x%x\n", cksum);
  }

  INFO("branching to code\n");
  void (*fn)(void) = (void (*)(void))payload.entry;
  INFO("0x%x\n", *(volatile u32*)fn);
  fn();
  WARN("code returned to bootloader\n");
}
