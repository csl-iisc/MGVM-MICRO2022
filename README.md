This repository contains the code for our MICRO 2022 paper:

<a href=https://www.csa.iisc.ac.in/~arkapravab/papers/MICRO22_MGvm.pdf> **Designing Virtual Memory System of MCM GPUs** </a> <br>
Pratheek B.\*, Neha Jawalkar\*, Arkaprava Basu <br>
In Proceedings of 55th ACM/IEEE International Symposium on Microarchitecture <br>
\* *Both Pratheek and Neha contributed equally.*


**Requirements:**

- golang 1.16

**Building and Running:**

1. Go to mgpusim/samples/\<benchmarkname\>
2. Run `go build` 
3. Run the executable generated with appropriate options (please see scripts/3\_gen\_runners.py for various options).

Please check the scripts directory for helper scripts to automatically compile, copy, and generate runners, and run benchmarks.

**Copyright**

Copyright (c) 2022 Indian Institute of Science

All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal with the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimers.
Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimers in the documentation and/or other materials provided with the distribution.
Neither the names of Computer Systems Lab, Indian Institute of Science, nor the names of its contributors may be used to endorse or promote products derived from this Software without specific prior written permission.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE CONTRIBUTORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS WITH THE SOFTWARE.
