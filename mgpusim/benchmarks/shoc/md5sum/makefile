all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg md5sum -private $^

kernels.disasm: md5sum.cl
	clang-ocl -mcpu=fiji -S -o $@ md5sum.cl

kernels.hsaco: md5sum.cl
	clang-ocl -mcpu=fiji -o $@ md5sum.cl
