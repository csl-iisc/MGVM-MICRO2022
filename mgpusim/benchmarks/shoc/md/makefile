all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg md -private $^

kernels.disasm: md.cl
	clang-ocl -mcpu=fiji -S -o $@ md.cl

kernels.hsaco: md.cl
	clang-ocl -mcpu=fiji -o $@ md.cl
