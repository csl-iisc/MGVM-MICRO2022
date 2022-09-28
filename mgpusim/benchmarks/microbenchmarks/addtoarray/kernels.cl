/**********************************************************************
 ********************************************************************/

__kernel void MicroBenchmark(__global uint* array, const uint index)
{

uint threadId = get_global_id(0);

uint element = array[threadId];
element = element + index;
array[threadId] = element;
}
