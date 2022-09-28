all: bindata.go

bindata.go: kernels.hsaco
	esc -o $@ -pkg reduction -private $^

kernels.hsaco: native/Reduction_Kernels.cl
	clang-ocl -mcpu=fiji -o $@ $^
