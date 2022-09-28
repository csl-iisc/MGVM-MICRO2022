// bleh

__kernel void IdleLoop( __global uint* array,
                       const uint loopCount) 
{
  uint value = loopCount;
  for (int iter = 0; iter < loopCount; iter++)
    {
      value = value + value;
    }
    array[0] = value;
}

// run with one thread
__kernel void NarrowStridedRead(__global uint* array, const uint start, 
                                     const uint end, const uint stride)
{
  uint value = 0;
  for (int iter = start; iter < end; iter+=stride)
    {
    int index = iter;
    uint part = array[index];
    value = value + part;
    }
    array[0] = value;
}

// run with more threads
__kernel void WideStridedRead(__global uint* array, const uint start,
                       const uint end, const uint threads, const uint stride)
{
  uint value = 0;
  int threadID = get_global_id(0);
  for (int iter = start; iter < end; iter+=stride)
    {
    int index = iter+threadID;
    value = value + array[index];
    }
}

// run with one thread
__kernel void NarrowStridedWrite(__global uint* array, const uint start, 
                                     const uint end, const uint stride)
{
  int index = 0;
  for (int iter = start; iter < end; iter+=stride)
  {
    index = iter;
    array[index] = index;
  }
}

// run with more threads
__kernel void WideStridedWrite(__global uint* array, const uint start,
                       const uint end, const uint threads, const uint stride)
{
  int threadID = get_global_id(0);
  for (int iter = start; iter < end; iter+=stride)
  {
    int index = iter+threadID;
    array[index] = 42;
  }
}

// run with more threads
__kernel void NarrowStridedReadRemote(__global uint* array, const uint start,
         const uint end, const uint stride, const uint loopCount,
                       const uint remoteStart)
{
  int groupID = get_group_id(0);
  int localID = get_local_id(0);
  if(localID == 0)
  {
    if (groupID < remoteStart)
    {
      uint value = loopCount;
      for (int iter = 0; iter < loopCount; iter++)
      {
        value = value + value;
      }
      array[0] = value;
    }
    else
    {
      uint value = 0;
      for (int iter = start; iter < end; iter+=stride)
      {
      int index = iter;
      uint part = array[index];
      value = value + value;
      }
      array[0] = value;
    }
  }
}

// single threaded pchase
// run with one thread
__kernel void PChase(__global uint* array, 
                              const uint start, const uint length)
{
  uint index = start;
  int value = 0;
  for (int iter = 0; iter < length; iter+=1)
    {
    uint part = array[index];
    value = value + part;
    index = part;
    }
    array[0] = value;
}

// single threaded pchase
// run with one thread and two blocks
__kernel void TwoBlockOneDelayedPChase(__global uint* array,
                    const uint loopCount, const uint secondblock,
                    const uint start1, const uint start2,
                    const uint length1, const uint length2)
{
  int groupID = get_group_id(0);
  int threadID = get_local_id(0);
  uint index = 0;
  int value = loopCount;
  if (groupID == 0)
  {
    if(threadID == 0)
    {
      index = start1;
      for (int iter = 0; iter < length1; iter+=1)
      {
        uint part = array[index];
        value = value + part;
        index = part;
      }
      array[0] = value;
    }
  }
  else if (groupID < secondblock)
  {
    if(threadID == 0)
    {
      value = loopCount;
      for (int iter = 0; iter < loopCount; iter++)
      {
        value = value + value;
      }
      array[0] = value;
    }
  }
  else
  {
    if(threadID == 0)
    {
    index = start2;
      value = loopCount;
      for (int iter = 0; iter < loopCount; iter++)
        {
          value = value + value;
        }
        array[0] = value;
      value = loopCount;
        for (int iter = 0; iter < length2; iter+=1)
        {
          uint part = array[index];
          value = value + part;
          index = part;
        }
        array[0] = value;
      }
  }
}

