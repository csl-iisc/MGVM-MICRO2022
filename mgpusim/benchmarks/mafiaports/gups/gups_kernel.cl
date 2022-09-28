/**********************************************************************
bleh
********************************************************************/

#pragma OPENCL EXTENSION cl_khr_int64_extended_atomics : enable

#define POLY 0x0000000000000007UL

#define NUPDATE (32 * TableSize)
#define THREADBLOCKSIZE 32
#define NTHREADSBLOCKS 128

__kernel void RandomAccessUpdate(ulong  TableSize, __global ulong* Table,
                                        __global ulong* Starts) {

  uint threadIdx =  get_global_id(0);
  ulong i, ran;

  ran = Starts[threadIdx];
  int numAccess = 1 * TableSize;
  numAccess = numAccess / THREADBLOCKSIZE;
  numAccess = numAccess / NTHREADSBLOCKS;

  for (i = 0; i < numAccess; i++) {
    ran = (ran << 1) ^ ((long)ran < 0 ? POLY : 0);
    ulong index = ran & (TableSize - 1);
    // atom_xor(&Table[index], ran);
    Table[index] ^= ran;
  }
}
