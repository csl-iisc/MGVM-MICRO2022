all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg color -private $^

kernels.disasm: kernels.hsaco
	/opt/rocm/hcc/bin/llvm-objdump -mcpu=fiji -disassemble $^ > $@

kernels.hsaco: native/kernel/kernel_max.cl
	clang-ocl -mcpu=fiji -o $@ $^
