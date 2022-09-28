all:
	clang-ocl -mcpu=fiji -S -o kernels.disasm native/mergesort.cl
	clang-ocl -mcpu=fiji -o kernels.hsaco native/mergesort.cl
	go generate
