PAYLOADS=$(wildcard *.c)
TARGETS=$(patsubst %.c,%.bin,$(PAYLOADS))

ASM=$(wildcard lib/*.S)
SRC=$(wildcard lib/*.c)

OBJ = $(ASM:.S=.o) $(SRC:.c=.o)

PREFIX=riscv64-unknown-elf
CC=$(PREFIX)-gcc
AS=$(PREFIX)-gcc
LD=$(PREFIX)-ld
OBJDUMP=$(PREFIX)-objdump
OBJCOPY=$(PREFIX)-objcopy

OPT ?= 3

COMMON_FLAGS=-MMD -ffreestanding -nostdlib -nostartfiles -march=rv64imafc -mabi=lp64
CFLAGS=$(COMMON_FLAGS) -O$(OPT) -Ilib
ASFLAGS=$(COMMON_FLAGS) -Ilib
LDFLAGS=-nostdlib -flto

all: tags $(TARGETS)

clean:
	rm -f **/*.o **/*.elf **/*.bin **/*.d tags

tags: $(wildcard *.[chS])
	ctags *.[chS]

ifneq ($(MAKECMDGOALS),clean)
-include $(wildcard *.d)
endif

%.bin: %.elf %.list
	$(OBJCOPY) $< -O binary $@
%.list: %.elf
	$(OBJDUMP) $< -D > $@

%.elf: memmap.ld %.o $(OBJ)
	$(LD) $(LDFLAGS) -T $^ -o $@

%.o: %.S
	$(AS) -c $< $(ASFLAGS) -o $@
%.o: %.c
	$(CC) -c $< $(CFLAGS) -o $@

.PHONY: all clean
.PRECIOUS: %.list
