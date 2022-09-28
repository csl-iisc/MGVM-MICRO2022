all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg syr2k -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/syr2k.cl
	clang-ocl -mcpu=fiji -o $@ $^
