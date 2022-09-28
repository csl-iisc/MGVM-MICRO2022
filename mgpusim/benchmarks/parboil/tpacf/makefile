all: bindata.go kernels.disasm

bindata.go: kernels.hsaco
	esc -o $@ -pkg tpacf -private $^

kernels.hsaco: native/kernel.cl
	clang-ocl -mcpu=fiji -o $@ $^
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/kernel.cl

kernels.disasm:
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/kernel.cl
