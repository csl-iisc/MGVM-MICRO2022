#!/usr/bin/python3

configs = ['private', 'shared', 'mgvm', 'mgvm-nobalance']

benchmarks = [
        'convolution2d',
        'fastwalshtransform',
        'gups',
        'jacobi1d',
        'jacobi2d',
        'kmeans',
        'matrixtranspose',
        'mis',
        'pagerank',
        'simpleconvolution',
        'shoc-reduction',
        'spmv',
        'stencil2d',
        'syrk',
        'syr2k',
        ]

for config in configs:
    for benchmark in benchmarks:
        print(config, benchmark)
        submit_file_name = config + '/' + benchmark + ".sh"
        submit_file = open(submit_file_name, "w")
        submit_file.write("#!/bin/bash\n")
        submit_file.write("cd samples\n")
        submit_file.write("cd " + benchmark + "\n")
        submit_file.write("./" + benchmark + " ")
        submit_file.write("-timing ")
        submit_file.write("-no-progress-bar ")
        submit_file.write("-report-all ")
        submit_file.write("-scheduling lasp ")

        if config == 'private':
            submit_file.write("-platform-type privatetlb ")
            submit_file.write("-mem-allocator-type lasp ")
            submit_file.write("-use-lasp-mem-alloc ")
        elif config == 'shared':
            submit_file.write("-platform-type xortlb ")
            submit_file.write("-mem-allocator-type lasp ")
            submit_file.write("-use-lasp-mem-alloc ")
            submit_file.write("-l2-tlb-striping 1 ")
        elif config == 'mgvm':
            submit_file.write("-platform-type customtlb ")
            submit_file.write("-use-lasp-hsl-mem-alloc ")
            submit_file.write("-switch-tlb-striping ")
        elif config == 'mgvm-nobalance':
            submit_file.write("-platform-type customtlb ")
            submit_file.write("-use-lasp-hsl-mem-alloc ")

        if benchmark == 'convolution2d':
            submit_file.write("-sched-partition Ydiv ")
        elif benchmark == 'fastwalshtransform':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'gups':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'jacobi1d':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'jacobi2d':
            submit_file.write("-sched-partition Ydiv ")
        elif benchmark == 'kmeans':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'matrixtranspose':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'mis':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'pagerank':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'simpleconvolution':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'shoc-reduction':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'spmv':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'stencil2d':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'syrk':
            submit_file.write("-sched-partition Xdiv ")
        elif benchmark == 'syr2k':
            submit_file.write("-sched-partition Xdiv ")

        # set appropriate HSL values
        if config == 'mgvm' or config == 'mgvm-nobalance':
            if benchmark == 'convolution2d':
                submit_file.write("-custom-hsl 16384 ")
                submit_file.write("-mem-allocator-type hslaware-32 ")
            if benchmark == 'fastwalshtransform':
                submit_file.write("-custom-hsl 2048 ")
                submit_file.write("-mem-allocator-type hslaware-4 ")
            if benchmark == 'gups':
                submit_file.write("-custom-hsl 1024 ")
                submit_file.write("-mem-allocator-type hslaware-2 ")
            if benchmark == 'jacobi1d':
                submit_file.write("-custom-hsl 16384 ")
                submit_file.write("-mem-allocator-type hslaware-32 ")
            if benchmark == 'jacobi2d':
                submit_file.write("-custom-hsl 4096 ")
                submit_file.write("-mem-allocator-type hslaware-8 ")
            if benchmark == 'kmeans':
                submit_file.write("-custom-hsl 4096 ")
                submit_file.write("-mem-allocator-type hslaware-8 ")
            if benchmark == 'matrixtranspose':
                submit_file.write("-custom-hsl 1024 ")
                submit_file.write("-mem-allocator-type hslaware-2 ")
            if benchmark == 'mis':
                submit_file.write("-custom-hsl 512 ")
                submit_file.write("-mem-allocator-type hslaware-1 ")
            if benchmark == 'pagerank':
                submit_file.write("-custom-hsl 8192 ")
                submit_file.write("-mem-allocator-type hslaware-16 ")
            if benchmark == 'simpleconvolution':
                submit_file.write("-custom-hsl 16384 ")
                submit_file.write("-mem-allocator-type hslaware-32 ")
            if benchmark == 'shoc-reduction':
                submit_file.write("-custom-hsl 16384 ")
                submit_file.write("-mem-allocator-type hslaware-32 ")
            if benchmark == 'spmv':
                submit_file.write("-custom-hsl 512 ")
                submit_file.write("-mem-allocator-type hslaware-1 ")
            if benchmark == 'stencil2d':
                submit_file.write("-custom-hsl 1024 ")
                submit_file.write("-mem-allocator-type hslaware-2 ")
            if benchmark == 'syrk':
                submit_file.write("-custom-hsl 1024 ")
                submit_file.write("-mem-allocator-type hslaware-2 ")
            if benchmark == 'syr2k':
                submit_file.write("-custom-hsl 512 ")
                submit_file.write("-mem-allocator-type hslaware-1 ")

        # limit super long benchmarks
        if benchmark == 'syrk':
            submit_file.write("-max-inst 10000000 ")
        if benchmark == 'syr2k':
            submit_file.write("-max-inst 30000000 ")

        # set benchmark specific parameters
        if benchmark == 'convolution2d':
            submit_file.write("-ni=8192 -nj=8192 ")
        if benchmark == 'fastwalshtransform':
            submit_file.write("-length=8388608 ")
        if benchmark == 'gups':
            submit_file.write(" ")
        if benchmark == 'jacobi1d':
            submit_file.write("-n=67108864 -steps=1")
        if benchmark == 'jacobi2d':
            submit_file.write("-n=4096 -steps=1")
        if benchmark == 'kmeans':
            submit_file.write("-points=524288 -features=32 -clusters=20 -max-iter=1 ")
        if benchmark == 'matrixtranspose':
            submit_file.write("-width=2048 ")
        if benchmark == 'mis':
            submit_file.write("-numNodes=524288 -numItems=1048576 ")
        if benchmark == 'pagerank':
            submit_file.write("-node=8192 -sparsity=0.5 -iterations=1 ")
        if benchmark == 'simpleconvolution':
            submit_file.write("-width=8190 -height=8190 ")
        if benchmark == 'shoc-reduction':
            submit_file.write("-Size=67108864 -Iterations=2 ")
        if benchmark == 'spmv':
            submit_file.write("-dim=2097152 -sparsity=0.00001 ")
        if benchmark == 'stencil2d':
            submit_file.write("-row=2048 -col=2048 ")
        if benchmark == 'syrk':
            submit_file.write("-ni=2048 -nj=2048 ")
        if benchmark == 'syr2k':
            submit_file.write("-ni=1024 -nj=1024 ")

