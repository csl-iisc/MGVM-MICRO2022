all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg blackscholes -private $^

kernels.hsaco: native/BlackScholes_Kernels.cl
	clang-ocl -mcpu=fiji -o $@ $^

kernels.disasm:
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/BlackScholes_Kernels.cl
