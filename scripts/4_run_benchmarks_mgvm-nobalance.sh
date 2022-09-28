#!/bin/bash

configs=("mgvm-nobalance")

benchmarks=("convolution2d" "fastwalshtransform" "gups" "jacobi1d" "jacobi2d" "kmeans" "matrixtranspose" "mis" "pagerank" "simpleconvolution" "shoc-reduction" "spmv" "stencil2d" "syrk" "syr2k")

for config in ${configs[@]}; 
do
  for benchmark in ${benchmarks[@]}; 
  do
    echo $config $benchmark
    cd $config
    pwd
    bash ${benchmark}.sh > output &
    cd ..
  done
done
        

