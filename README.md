# Ox64 Baremetal Bootloader

The BL808 CPU on the [Ox64](https://wiki.pine64.org/wiki/Ox64) has a rather
convoluted boot process.  The vendor-provided instructions involve linking
against the chip manufacturer's poorly-documented SDK, which requires a custom
fork of GCC to even build.  Once a binary is built, it must be flashed with a
custom flashing tool, which again assumes that a binary was built against the
manufacturer SDK.  Even if you manage to get a binary on the board without
using the SDK, critical board-specific initialization functions are left undone
by the bootrom, so basic things like reading the system timer and accessing the
64MB of RAM require additional setup steps (which, as always, are
undocumented). Rather than bothering with all that, this repo contains a
minimal flashing tool and bootloader which are completely standalone.

## The Ox64

The Ox64 has an unusual architecture; it has three CPUs with different ISAs.

* The "D0" (or "MM", or "DSP") CPU is the primary one; it is an RV64IMAFC
  T-Head C906 core running at 480 MHz.  It has an MMU and is capable of running
  Linux.
* The "M0" (or "MCU") CPU is a device coprocessor; it is an RV32IMAFC T-Head
  E906 core running at 320 MHz.  Annoyingly, many peripherals only send
  interrupts to this core, so it needs to relay them to the D0.
* The "LP" CPU is a low-power coprocessor; it is an RV32EMC T-Head E902 core
  running at 150 MHz.  Not much is known about this CPU.

## Flashing Tool
```
Usage of ox-flash:
  -addr uint
    	the address to flash to (default: 0x0 for firmware, 0x10000 for a payload)
  -baud uint
    	baud rate (default 2000000)
  -bootloader
    	flash a bootloader
  -echo
    	echo bytes received from the serial port after flashing
  -echo-baud uint
    	baud rate for echoing (default 115200)
  -payload
    	flash a payload
  -port string
    	serial port
  -run
    	exit the bootrom and run the program from flash
```

The flashing tool is set up for a use case where there's a "bootloader" (e.g.,
the bootloader provided here), and a "payload" (an arbitrary RISC-V binary). It
tells the bootrom to run the provided binary on the 64-bit "D0" core.

Note that the addresses provided here are actually offsets within flash, which
itself is at `0x58000000`.  Payloads should instead be linked to absolute
physical addresses within PSRAM (e.g., `0x50000000`); the bootloader will take
care of copying the binary to the requested address.

## Bootloader

The default SDK does many different setup steps.  For now, the bootloader only does the bare minimum necessary to have a useable computing environment; a few more steps are necessary for a truly useful setup.

- [x] Initializes UART0 at 115200 baud.
- [x] Configures pins 14 and 15 as UART pins.
- [x] Routes TX0 and RX0 to pins 14 and 15.
- [x] Initializes the 64MB PSRAM.
- [x] Starts the RISC-V `TIME` system timer at 1 MHz.
- [x] Copies the payload out of flash into PSRAM.
- [x] Starts running the payload on the D0.
- [ ] Initalizes JTAG so the payload can be debugged.
- [ ] Starts an M0 helper program to forward peripheral interrupts to the D0.

## Payload

The payload is a 64-bit RISC-V binary with a simple header prefixed:
```C
struct payload {
  uint64_t entry; // the load address of the binary
  uint32_t size;  // the size of the binary
  uint32_t cksum; // a CRC-32 (IEEE 802.3) of the binary
  uint8_t data[]; // the binary itself
};
```

The payload will be executed in *m-mode*, with full access to physical memory.
The bootloader **does not** provide an SBI implementation; it is assumed that
the payload will completely take over once control is transferred.

## References
* [Reference Manual](https://files.pine64.org/doc/datasheet/ox64/BL808_RM_en_1.0(open).pdf)
* [Datasheet](https://files.pine64.org/doc/datasheet/ox64/BL808_DS_en_1.1(open).pdf)
* [Product Page](https://wiki.pine64.org/wiki/Ox64)
* [SDK](https://github.com/bouffalolab/bl808_linux)
* [OpenC906 Manual](https://occ-intl-prod.oss-ap-southeast-1.aliyuncs.com/resource/XuanTie-OpenC906-UserManual.pdf)
* [OpenC906 Source](https://github.com/T-head-Semi/openc906)
