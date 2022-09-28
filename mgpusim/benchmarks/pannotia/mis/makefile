all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg mis -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/kernel/kernel.cl
	clang-ocl -mcpu=fiji -o $@ $^
