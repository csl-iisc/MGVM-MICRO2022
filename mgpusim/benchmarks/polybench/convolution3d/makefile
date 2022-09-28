all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg convolution3d -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/3DConvolution.cl
	clang-ocl -mcpu=fiji -o $@ $^
