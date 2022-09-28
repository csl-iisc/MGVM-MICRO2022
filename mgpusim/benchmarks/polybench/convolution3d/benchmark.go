// Package convolution3d implements the convolution3d benchmark from Polybench.
package convolution3d

import (
	"log"
	"math"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	NI                  int32
	NJ                  int32
	NK                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpus    []int
	queues  []*driver.CommandQueue
	kernel1 *insts.HsaCo

	a, b            []float32
	NI, NJ, NK      int
	da, db          driver.GPUPtr
	b_outputFromGPU []float32

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	return b
}

// SelectGPU selects GPU
func (b *Benchmark) SelectGPU(gpus []int) {
	b.gpus = gpus
}

// SetUnifiedMemory uses Unified Memory
func (b *Benchmark) SetUnifiedMemory() {
	b.useUnifiedMemory = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "Convolution3D_kernel")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
}

// Run runs
func (b *Benchmark) Run() {
	for _, gpu := range b.gpus {
		b.driver.SelectGPU(b.context, gpu)
		b.queues = append(b.queues, b.driver.CreateCommandQueue(b.context))
	}

	b.initMem()
	b.exec()
}

func (b *Benchmark) initMem() {
	rand.Seed(1)
	b.a = make([]float32, b.NI*b.NJ*b.NK)
	b.b = make([]float32, b.NI*b.NJ*b.NK)
	b.b_outputFromGPU = make([]float32, b.NI*b.NJ*b.NK)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = float32(i*j) / (float32(b.NI * b.NJ))
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*b.NK*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*b.NK*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*b.NK*4), "test")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*b.NK*4), "test")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*b.NK*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*b.NK*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.NJ-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.NK-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		b.db,
		int32(b.NI),
		int32(b.NJ),
		int32(b.NK),
		1,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	b.driver.MemCopyD2H(b.context, b.b_outputFromGPU, b.db)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuconvolution3d()

	for i := 1; i < b.NI-1; i++ {
		for j := 1; j < b.NJ-1; j++ {
			for k := 1; k < b.NK-1; k++ {
				if math.Abs(float64(b.b_outputFromGPU[i*b.NJ*b.NK+j*b.NK+k]-b.b[i*b.NJ*b.NK+j*b.NK+k])) > 0.001 {
					log.Panicf("Mismatch at %d, %d, %d, expected %f, but get %f",
						i, j, k, // i*b.NJ*b.NK+j*b.NK+k,
						b.b[i*b.NJ*b.NK+j*b.NK+k],
						b.b_outputFromGPU[i*b.NJ*b.NK+j*b.NK+k])
				}
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuconvolution3d() {
	c11 := float32(+2)
	c21 := float32(+5)
	c31 := float32(-8)
	c12 := float32(-3)
	c22 := float32(+6)
	c32 := float32(-9)
	c13 := float32(+4)
	c23 := float32(+7)
	c33 := float32(+10)

	// 	B[i*(nk * nj) + j*nk + k] =
	// 	c11 * A[(i - 1)*(nk * nj) + (j - 1)*nk + (k - 1)]  +  c13 * A[(i + 1)*(nk * nj) + (j - 1)*nk + (k - 1)]
	// +   c21 * A[(i - 1)*(nk * nj) + (j - 1)*nk + (k - 1)]  +  c23 * A[(i + 1)*(nk * nj) + (j - 1)*nk + (k - 1)]
	// +   c31 * A[(i - 1)*(nk * nj) + (j - 1)*nk + (k - 1)]  +  c33 * A[(i + 1)*(nk * nj) + (j - 1)*nk + (k - 1)]
	// +   c12 * A[(i + 0)*(nk * nj) + (j - 1)*nk + (k + 0)]  +  c22 * A[(i + 0)*(nk * nj) + (j + 0)*nk + (k + 0)]
	// +   c32 * A[(i + 0)*(nk * nj) + (j + 1)*nk + (k + 0)]  +  c11 * A[(i - 1)*(nk * nj) + (j - 1)*nk + (k + 1)]
	// +   c13 * A[(i + 1)*(nk * nj) + (j - 1)*nk + (k + 1)]  +  c21 * A[(i - 1)*(nk * nj) + (j + 0)*nk + (k + 1)]
	// +   c23 * A[(i + 1)*(nk * nj) + (j + 0)*nk + (k + 1)]  +  c31 * A[(i - 1)*(nk * nj) + (j + 1)*nk + (k + 1)]
	// +   c33 * A[(i + 1)*(nk * nj) + (j + 1)*nk + (k + 1)];
	for i := 1; i < b.NI-1; i++ {
		for j := 1; j < b.NJ-1; j++ {
			for k := 1; k < b.NK-1; k++ {
				b.b[i*b.NJ*b.NK+j*b.NK+k] =
					c11*b.a[(i-1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c13*b.a[(i+1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c21*b.a[(i-1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c23*b.a[(i+1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c31*b.a[(i-1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c33*b.a[(i+1)*b.NJ*b.NK+(j-1)*b.NK+(k-1)] +
						c12*b.a[(i+0)*b.NJ*b.NK+(j-1)*b.NK+(k+0)] +
						c22*b.a[(i+0)*b.NJ*b.NK+(j+0)*b.NK+(k+0)] +
						c32*b.a[(i+0)*b.NJ*b.NK+(j+1)*b.NK+(k+0)] +
						c11*b.a[(i-1)*b.NJ*b.NK+(j-1)*b.NK+(k+1)] +
						c13*b.a[(i+1)*b.NJ*b.NK+(j-1)*b.NK+(k+1)] +
						c21*b.a[(i-1)*b.NJ*b.NK+(j+0)*b.NK+(k+1)] +
						c23*b.a[(i+1)*b.NJ*b.NK+(j+0)*b.NK+(k+1)] +
						c31*b.a[(i-1)*b.NJ*b.NK+(j+1)*b.NK+(k+1)] +
						c33*b.a[(i+1)*b.NJ*b.NK+(j+1)*b.NK+(k+1)]
			}
		}
	}
}
