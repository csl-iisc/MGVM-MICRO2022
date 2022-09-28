all: bindata.go kernels.disasm

bindata.go: kernels.hsaco
	esc -o $@ -pkg gups -private $^

kernels.hsaco: gups_kernel.cl
	clang-ocl -mcpu=fiji -o $@ $^
	clang-ocl -mcpu=fiji -S -o kernels.disasm gups_kernel.cl

kernels.disasm:
	clang-ocl -mcpu=fiji -S -o kernels.disasm gups_kernel.cl
