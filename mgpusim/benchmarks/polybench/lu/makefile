all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg lu -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/lu.cl
	clang-ocl -mcpu=fiji -o $@ $^
