all: bindata.go kernels.disasm

bindata.go: kernels.hsaco
	esc -o $@ -pkg streamcluster -private $^

kernels.hsaco: native/Kernels.cl
	clang-ocl -mcpu=fiji -o $@ $^
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/Kernels.cl

kernels.disasm:
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/Kernels.cl
