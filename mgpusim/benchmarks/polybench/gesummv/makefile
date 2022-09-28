all: esc.go

esc.go: kernels.hsaco
	esc -o $@ -pkg gesummv -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/gesummv.cl
	clang-ocl -mcpu=fiji -o $@ $^
