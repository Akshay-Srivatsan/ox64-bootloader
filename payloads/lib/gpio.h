#pragma once

#include "types.h"

enum gpio_alternate_function {
    GPIO_SDH,
    GPIO_SPI0,
    GPIO_FLASH,
    GPIO_I2S,
    GPIO_PDM,
    GPIO_I2C0,
    GPIO_I2C1,
    GPIO_UART,
    GPIO_EMAC,
    GPIO_CAM,
    GPIO_ANALOG,
    GPIO_GPIO,
    GPIO_PWM0,
    GPIO_PWM1,
    GPIO_SPI1,
    GPIO_I2C2,
    GPIO_I2C3,
    GPIO_MM_UART,
    GPIO_DBI_B,
    GPIO_DBI_C,
    GPIO_DPI,
    GPIO_JTAG_LP,
    GPIO_JTAG_M0,
    GPIO_JTAG_D0,
    GPIO_CLOCK_OUT,
};

typedef size_t gpio_t;

void gpio_set_input(gpio_t pin, bool enable);
void gpio_set_output(gpio_t pin, bool enable);
void gpio_set_function(gpio_t pin, enum gpio_alternate_function function);
