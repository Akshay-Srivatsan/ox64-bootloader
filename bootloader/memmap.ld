ENTRY(_start)

MEMORY
{
  OCRAM (rwx) : ORIGIN = 0x22020000, LENGTH = 64K
  WRAM (rwx) : ORIGIN = 0x22030000, LENGTH = 160K
  DRAM (rwx) : ORIGIN = 0x3ef80000, LENGTH = 16K
  VRAM (rwx) : ORIGIN = 0x3f000000, LENGTH = 16K
  XRAM (rwx) : ORIGIN = 0x40000000, LENGTH = 16K

  PSRAM (rwx) : ORIGIN = 0x50000000, LENGTH = 64M

  /*
   * This is actually something like this:
   *
   * FLASH (rx) : ORIGIN = 0x58000000, LENGTH = 64M
   *
   * But we use the linker script to enforce that the bootloader is less than
   * 64K (so the payload consistently starts at 0x58010000).
   */
  FLASH (rx) : ORIGIN = 0x58000000, LENGTH = 64K
  PAYLOAD (rx) : ORIGIN = 0x58010000, LENGTH = ORIGIN(FLASH)+64M-ORIGIN(PAYLOAD)

  ROM (rx) : ORIGIN = 0x90000000, LENGTH = 128K
}
ASSERT((ORIGIN(FLASH) + LENGTH(FLASH) == ORIGIN(PAYLOAD)), "Error: Memory configured incorrectly. (FLASH and PAYLOAD are not contiguous)");
ASSERT((ORIGIN(PAYLOAD) + LENGTH(PAYLOAD) == ORIGIN(FLASH) + 64M), "Error: Memory configured incorrectly. (FLASH + PAYLOAD != 64M)");

SECTIONS
{
    . = ORIGIN(FLASH);
    __stack_size = 2K;
    .text : {
        _kcode_start = .;
        KEEP(*(.text.boot))  
        *(.text*) 
        _kcode_end = .;
        . = ALIGN(8);
    } > FLASH
    .rodata : {
        . = ALIGN(8);
        _krodata_start = .;
        *(.rodata*)
        *(.srodata*)
        . = ALIGN(8);
        _krodata_end = .;
    } > FLASH
    .data : { 
        . = ALIGN(8);
        _kdata_start = .;
        _kdata_start_load = LOADADDR(.data);
        __global_pointer$ = . + 0x800;
        *(.sdata*)
        *(.data*)
        . = ALIGN(8);
        _kdata_end = .;
    } > OCRAM AT> FLASH
    .bss (NOLOAD) : {
        . = ALIGN(8);
        _kbss_start = .;
        *(.sbss*)
        *(.bss*)
        *(COMMON)
        . = ALIGN(8);
        _kbss_end = .;
    } > OCRAM

    .stack (NOLOAD) : {
      __stack_bottom__ = .;
      . += __stack_size;
      __stack_top__ = .;
    } > OCRAM

    /DISCARD/ : {
      *(.comment)
      *(.riscv.attributes)
      *(.note)
      *(.eh_frame)
    }
}

PROVIDE(payload = ORIGIN(PAYLOAD));
