#include "gpio.h"
#include "memory.h"

struct PACKED gpio_config {
  bool input_enable : 1;
  bool schmitt_trigger_enable : 1;
  unsigned drive_strength : 2;
  bool pullup_enable : 1;
  bool pulldown_enable : 1;
  bool output_enable : 1;
  unsigned : 1;
  unsigned alternate_function : 5;
  unsigned : 3;
  unsigned interrupt_mode : 4;
  bool interrupt_clear : 1;
  bool interrupt_state : 1;
  bool interrupt_mask : 1;
  unsigned : 1;
  bool output_value : 1;
  bool output_set : 1;
  bool output_clear : 1;
  unsigned : 1;
  bool input_value : 1;
  unsigned : 1;
  unsigned io_mode : 2;
};

static volatile struct gpio_config *const GPIO_CONFIG =
    (volatile struct gpio_config *)0x200008c4;

enum gpio_io_mode {
  IO_MODE_NORMAL,
  IO_MODE_SET_CLEAR,
  IO_MODE_BUFFER,
  IO_MODE_CACHE,
};

void gpio_set_input(gpio_t pin, bool enable) {
  struct gpio_config x = get32_type(struct gpio_config, GPIO_CONFIG + pin);
  x.io_mode = IO_MODE_SET_CLEAR;
  x.input_enable = enable;
  put32_type(struct gpio_config, GPIO_CONFIG + pin, x);
}

void gpio_set_output(gpio_t pin, bool enable) {
  struct gpio_config x = get32_type(struct gpio_config, GPIO_CONFIG + pin);
  x.io_mode = IO_MODE_SET_CLEAR;
  x.output_enable = enable;
  x.drive_strength = 3;
  put32_type(struct gpio_config, GPIO_CONFIG + pin, x);
}

void gpio_set_function(gpio_t pin, enum gpio_alternate_function function) {
  struct gpio_config x = get32_type(struct gpio_config, GPIO_CONFIG + pin);
  x.alternate_function = function;
  put32_type(struct gpio_config, GPIO_CONFIG + pin, x);
}
