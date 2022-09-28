all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg jacobi1d -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/jacobi1D.cl
	clang-ocl -mcpu=fiji -o $@ $^
