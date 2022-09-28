all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg reduction -private $^

kernels.disasm: reduction.cl
	clang-ocl -mcpu=fiji -S -o $@ reduction.cl

kernels.hsaco: reduction.cl
	clang-ocl -mcpu=fiji -o $@ reduction.cl
