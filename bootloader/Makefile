TARGET=bootloader.bin

ASM=$(wildcard *.S)
SRC=$(wildcard *.c)

OBJ = $(ASM:.S=.o) $(SRC:.c=.o)

PREFIX=riscv64-unknown-elf
CC=$(PREFIX)-gcc
AS=$(PREFIX)-gcc
LD=$(PREFIX)-ld
OBJDUMP=$(PREFIX)-objdump
OBJCOPY=$(PREFIX)-objcopy

OPT ?= 3

COMMON_FLAGS=-MMD -ffreestanding -nostdlib -nostartfiles -march=rv64imafc -mabi=lp64
CFLAGS=$(COMMON_FLAGS) -O$(OPT)
ASFLAGS=$(COMMON_FLAGS)
LDFLAGS=-nostdlib -flto

all: tags $(TARGET)

clean:
	rm -f *.o *.elf *.bin *.d tags

tags: $(wildcard *.[chS])
	ctags *.[chS]

ifneq ($(MAKECMDGOALS),clean)
-include $(wildcard *.d)
endif

%.bin: %.elf %.list
	$(OBJCOPY) $< -O binary $@
%.list: %.elf
	$(OBJDUMP) $< -D > $@

bootloader.elf: memmap.ld $(OBJ)
	$(LD) $(LDFLAGS) -T $^ -o $@

%.o: %.S
	$(AS) -c $< $(ASFLAGS) -o $@
%.o: %.c
	$(CC) -c $< $(CFLAGS) -o $@

run: $(TARGET)
	ox-flash --bootloader --run --echo $(TARGET)

.PHONY: all clean run
.PRECIOUS: %.list
