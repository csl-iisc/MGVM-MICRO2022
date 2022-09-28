all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg mm3 -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/3mm.cl
	clang-ocl -mcpu=fiji -o $@ $^
