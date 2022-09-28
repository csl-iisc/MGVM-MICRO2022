/**
 * adi.c: This file is part of the PolyBench/GPU 1.0 test suite.
 *
 *
 * Contact: Scott Grauer-Gray <sgrauerg@gmail.com>
 * Will Killian <killian@udel.edu>
 * Louis-Noel Pouchet <pouchet@cse.ohio-state.edu>
 * Web address: http://www.cse.ohio-state.edu/~pouchet/software/polybench/GPU
 */

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>
#include <math.h>

#ifdef __APPLE__
#include <OpenCL/opencl.h>
#else
#include <CL/cl.h>
#endif

#define POLYBENCH_TIME 1

//select the OpenCL device to use (can be GPU, CPU, or Accelerator such as Intel Xeon Phi)
#define OPENCL_DEVICE_SELECTION CL_DEVICE_TYPE_GPU

#include "adi.h"
#include <polybench.h>
#include <polybenchUtilFuncts.h>

//define the error threshold for the results "not matching"
#define PERCENT_DIFF_ERROR_THRESHOLD 10.05

#define GPU_DEVICE 0

#define MAX_SOURCE_SIZE (0x10000000)

#if defined(cl_khr_fp64)  // Khronos extension available?
#pragma OPENCL EXTENSION cl_khr_fp64 : enable
#elif defined(cl_amd_fp64)  // AMD extension available?
#pragma OPENCL EXTENSION cl_amd_fp64 : enable
#endif

char str_temp[1024];

cl_platform_id platform_id;
cl_device_id device_id;   
cl_uint num_devices;
cl_uint num_platforms;
cl_int errcode;
cl_context clGPUContext;
cl_kernel clKernel1;
cl_kernel clKernel2;
cl_kernel clKernel3;
cl_kernel clKernel4;
cl_kernel clKernel5;
cl_kernel clKernel6;
cl_command_queue clCommandQue;
cl_program clProgram;
cl_mem a_mem_obj;
cl_mem b_mem_obj;
cl_mem x_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;
unsigned int mem_size_A;
unsigned int mem_size_B;
unsigned int mem_size_X;

#define RUN_ON_CPU


void init_array(int n, DATA_TYPE POLYBENCH_2D(A,N,N,n,n), DATA_TYPE POLYBENCH_2D(B,N,N,n,n), DATA_TYPE POLYBENCH_2D(X,N,N,n,n))
{
  	int i, j;

  	for (i = 0; i < n; i++)
	{
    		for (j = 0; j < n; j++)
      		{
			X[i][j] = ((DATA_TYPE) i*(j+1) + 1) / N;
			A[i][j] = ((DATA_TYPE) (i-1)*(j+4) + 2) / N;
			B[i][j] = ((DATA_TYPE) (i+3)*(j+7) + 3) / N;
      		}
	}
}


void compareResults(int n, DATA_TYPE POLYBENCH_2D(B_cpu,N,N,n,n), DATA_TYPE POLYBENCH_2D(B_fromGpu,N,N,n,n), DATA_TYPE POLYBENCH_2D(X_cpu,N,N,n,n), 
			DATA_TYPE POLYBENCH_2D(X_fromGpu,N,N,n,n))
{
	int i, j, fail;
	fail = 0;
	
	// Compare b and x output on cpu and gpu
	for (i=0; i < n; i++) 
	{
		for (j=0; j < n; j++) 
		{
			if (percentDiff(B_cpu[i][j], B_fromGpu[i][j]) > PERCENT_DIFF_ERROR_THRESHOLD) 
			{
				fail++;
			}
		}
	}
	
	for (i=0; i<n; i++) 
	{
		for (j=0; j<n; j++) 
		{
			if (percentDiff(X_cpu[i][j], X_fromGpu[i][j]) > PERCENT_DIFF_ERROR_THRESHOLD) 
			{
				fail++;
			}
		errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	clFinish(clCommandQue);
}

void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel1);
	errcode = clReleaseKernel(clKernel2);
	errcode = clReleaseKernel(clKernel3);
	errcode = clReleaseKernel(clKernel4);
	errcode = clReleaseKernel(clKernel5);
	errcode = clReleaseKernel(clKernel6);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(a_mem_obj);
	errcode = clReleaseMemObject(b_mem_obj);
	errcode = clReleaseMemObject(x_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}


void adi(int tsteps, int n, DATA_TYPE POLYBENCH_2D(A,N,N,n,n), DATA_TYPE POLYBENCH_2D(B,N,N,n,n), DATA_TYPE POLYBENCH_2D(X,N,N,n,n))
{
	int t, i1, i2;
	for (t = 0; t < _PB_TSTEPS; t++)
    	{
    		for (i1 = 0; i1 < _PB_N; i1++)
		{
			for (i2 = 1; i2 < _PB_N; i2++)
			{
				X[i1][i2] = X[i1][i2] - X[i1][(i2-1)] * A[i1][i2] / B[i1][(i2-1)];
				B[i1][i2] = B[i1][i2] - A[i1][i2] * A[i1][i2] / B[i1][(i2-1)];
			}
		}

	   	for (i1 = 0; i1 < _PB_N; i1++)
		{
			X[i1][(N-1)] = X[i1][(N-1)] / B[i1][(N-1)];
		}

	   	for (i1 = 0; i1 < _PB_N; i1++)
		{
			for (i2 = 0; i2 < _PB_N-2; i2++)
			{
				X[i1][(N-i2-2)] = (X[i1][(N-2-i2)] - X[i1][(N-2-i2-1)] * A[i1][(N-i2-3)]) / B[i1][(N-3-i2)];
			}
		}

	   	for (i1 = 1; i1 < _PB_N; i1++)
		{
			for (i2 = 0; i2 < _PB_N; i2++) 
			{
		  		X[i1][i2] = X[i1][i2] - X[(i1-1)][i2] * A[i1][i2] / B[(i1-1)][i2];
		  		B[i1][i2] = B[i1][i2] - A[i1][i2] * A[i1][i2] / B[(i1-1)][i2];
			}
		}

	   	for (i2 = 0; i2 < _PB_N; i2++)
		{
			X[(N-1)][i2] = X[(N-1)][i2] / B[(N-1)][i2];
		}

	   	for (i1 = 0; i1 < _PB_N-2; i1++)
		{
			for (i2 = 0; i2 < _PB_N; i2++)
			{
		 	 	X[(N-2-i1)][i2] = (X[(N-2-i1)][i2] - X[(N-i1-3)][i2] * A[(N-3-i1)][i2]) / B[(N-2-i1)][i2];
			}
		}
    }
}


/* DCE code. Must scan the entire live-out data.
   Can be used also to check the correctness of the output. */
static
void print_array(int n,
		 DATA_TYPE POLYBENCH_2D(X,N,N,n,n))

{
  int i, j;

  for (i = 0; i < n; i++)
    for (j = 0; j < n; j++) {
      fprintf(stderr, DATA_PRINTF_MODIFIER, X[i][j]);
      if ((i * N + j) % 20 == 0) fprintf(stderr, "\n");
    }
  fprintf(stderr, "\n");
}


int main(int argc, char *argv[])
{
	int tsteps = TSTEPS;
	int n = N;

	POLYBENCH_2D_ARRAY_DECL(A,DATA_TYPE,N,N,n,n);
	POLYBENCH_2D_ARRAY_DECL(B,DATA_TYPE,N,N,n,n);
	POLYBENCH_2D_ARRAY_DECL(B_outputFromGpu,DATA_TYPE,N,N,n,n);
	POLYBENCH_2D_ARRAY_DECL(X,DATA_TYPE,N,N,n,n);
	POLYBENCH_2D_ARRAY_DECL(X_outputFromGpu,DATA_TYPE,N,N,n,n);

	init_array(n, POLYBENCH_ARRAY(A), POLYBENCH_ARRAY(B), POLYBENCH_ARRAY(X));

	read_cl_file();
	cl_initialization();
	cl_mem_init(POLYBENCH_ARRAY(A), POLYBENCH_ARRAY(B), POLYBENCH_ARRAY(X));
	cl_load_prog();
	
	/* Start timer. */
  	polybench_start_instruments;

	int t, i1;

	for (t = 0; t < _PB_TSTEPS; t++)
	{
		cl_launch_kernel1(n);

		cl_launch_kernel2(n);

		cl_launch_kernel3(n);
	
		for (i1 = 1; i1 < _PB_N; i1++)
		{
			cl_launch_kernel4(i1, n);
		}

		cl_launch_kernel5(n);
		
		for (i1 = 0; i1 < _PB_N-2; i1++)
		{
			cl_launch_kernel6(i1, n);
		}
	}	
	
	/* Stop and print timer. */
	printf("GPU Time in seconds:\n");
  	polybench_stop_instruments;
 	polybench_print_instruments;

	errcode = clEnqueueReadBuffer(clCommandQue, b_mem_obj, CL_TRUE, 0, mem_size_B, POLYBENCH_ARRAY(B_outputFromGpu), 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");
	errcode = clEnqueueReadBuffer(clCommandQue, x_mem_obj, CL_TRUE, 0, mem_size_X, POLYBENCH_ARRAY(X_outputFromGpu), 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");
	
	#ifdef RUN_ON_CPU
	
		/* Start timer. */
	  	polybench_start_instruments;

		adi(tsteps, n, POLYBENCH_ARRAY(A), POLYBENCH_ARRAY(B), POLYBENCH_ARRAY(X));
	
		/* Stop and print timer. */
		printf("CPU Time in seconds:\n");
	  	polybench_stop_instruments;
	 	polybench_print_instruments;

		compareResults(n, POLYBENCH_ARRAY(B), POLYBENCH_ARRAY(B_outputFromGpu), POLYBENCH_ARRAY(X), POLYBENCH_ARRAY(X_outputFromGpu));

	#else //prevent dead code elimination

		polybench_prevent_dce(print_array(n, POLYBENCH_ARRAY(X_outputFromGpu)));

	#endif //RUN_ON_CPU

	cl_clean_up();

	POLYBENCH_FREE_ARRAY(A);
	POLYBENCH_FREE_ARRAY(B);
	POLYBENCH_FREE_ARRAY(B_outputFromGpu);
	POLYBENCH_FREE_ARRAY(X);
	POLYBENCH_FREE_ARRAY(X_outputFromGpu);

    return 0;
}

#include <polybench.c>
