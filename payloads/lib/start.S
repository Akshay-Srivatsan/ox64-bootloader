.section ".text.boot"

.globl _start
_start:
.option push
.option norelax
  la sp, __stack_top__
  la gp, __global_pointer$
.option pop
  call _cstart
halt:
  nop
  j halt

.globl put32
put32:
  fence rw,rw
  sw a1, 0(a0)
  fence rw,rw
  ret

.globl get32
get32:
  fence rw,rw
  lw a0, 0(a0)
  ret
