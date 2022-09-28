all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg jacobi2d -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/jacobi2D.cl
	clang-ocl -mcpu=fiji -o $@ $^
