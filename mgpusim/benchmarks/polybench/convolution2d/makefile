all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg convolution2d -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/2DConvolution.cl
	clang-ocl -mcpu=fiji -o $@ $^
