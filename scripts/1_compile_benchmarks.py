#!/usr/bin/python3

import subprocess
import os
import sys
import argparse
import re
import csv

path = "../mgpusim/"
benchmarks_path = path + "samples/"


class Test(object):
    """ define a benchmark to test """

    def __init__(self, path):
        self.path = path

    def compile(self):
        fp = open(os.devnull, 'w')
        p = subprocess.Popen('go build', shell=True,
                             cwd=self.path, stdout=fp, stderr=fp)
        p.wait()
        if p.returncode == 0:
            print("Compiled " + self.path, 'green')
            return False
        else:
            print("Compile failed " + self.path, 'red')
            return True



def main():

    adi = Test(benchmarks_path + 'adi')
    aes = Test(benchmarks_path + 'aes')
    atax = Test(benchmarks_path + 'atax')
    bfs = Test(benchmarks_path + 'bfs')
    bicg = Test(benchmarks_path + 'bicg')
    bitonicsort = Test(benchmarks_path + 'bitonicsort')
    color = Test(benchmarks_path + 'color')
    convolution2d = Test(benchmarks_path + 'convolution2d')
    convolution3d = Test(benchmarks_path + 'convolution3d')
    correlation = Test(benchmarks_path + 'correlation')
    covariance = Test(benchmarks_path + 'covariance')
    doitgen = Test(benchmarks_path + 'doitgen')
    fastwalshtransform = Test(benchmarks_path + 'fastwalshtransform')
    fdtd2d = Test(benchmarks_path + 'fdtd2d')
    fft = Test(benchmarks_path + 'fft')
    fir = Test(benchmarks_path + 'fir')
    floydwarshall = Test(benchmarks_path + 'floydwarshall')
    gemm = Test(benchmarks_path + 'gemm')
    gemver = Test(benchmarks_path + 'gemver')
    gesummv = Test(benchmarks_path + 'gesummv')
    gramschmidt = Test(benchmarks_path + 'gramschmidt')
    gups = Test(benchmarks_path + 'gups')
    im2col = Test(benchmarks_path + 'im2col')
    jacobi1d = Test(benchmarks_path + 'jacobi1d')
    jacobi2d = Test(benchmarks_path + 'jacobi2d')
    kmeans = Test(benchmarks_path + 'kmeans')
    lenet = Test(benchmarks_path + 'lenet')
    lu = Test(benchmarks_path + 'lu')
    matrixmultiplication = Test(benchmarks_path + 'matrixmultiplication')
    matrixtranspose = Test(benchmarks_path + 'matrixtranspose')
    maxpooling = Test(benchmarks_path + 'maxpooling')
    mineva = Test(benchmarks_path + 'mineva')
    mis = Test(benchmarks_path + 'mis')
    mm2 = Test(benchmarks_path + 'mm2')
    mm3 = Test(benchmarks_path + 'mm3')
    mvt = Test(benchmarks_path + 'mvt')
    nbody = Test(benchmarks_path + 'nbody')
    pagerank = Test(benchmarks_path + 'pagerank')
    relu = Test(benchmarks_path + 'relu')
    simpleconvolution = Test(benchmarks_path + 'simpleconvolution')
    shocreduction = Test(benchmarks_path + 'shoc-reduction')
    spmv = Test(benchmarks_path + 'spmv')
    sssp = Test(benchmarks_path + 'sssp')
    stencil2d = Test(benchmarks_path + 'stencil2d')
    syrk = Test(benchmarks_path + 'syrk')
    syr2k = Test(benchmarks_path + 'syr2k')
    vgg16 = Test(benchmarks_path + 'vgg16')
    xor = Test(benchmarks_path + 'xor')

    err = False

    # err |= adi.compile()
    # err |= aes.compile()
    # err |= atax.compile()
    # err |= bfs.compile()
    # err |= bicg.compile()
    # err |= bitonicsort.compile()
    # err |= color.compile()
    err |= convolution2d.compile()
    # err |= convolution3d.compile()
    # err |= correlation.compile()
    # err |= covariance.compile()
    # err |= doitgen.compile()
    err |= fastwalshtransform.compile()
    # err |= fdtd2d.compile()
    # err |= fir.compile()
    # err |= fft.compile()
    # err |= floydwarshall.compile()
    # err |= gemm.compile()
    # err |= gesummv.compile()
    # err |= gemver.compile()
    # err |= gramschmidt.compile()
    err |= gups.compile()
    # err |= im2col.compile()
    err |= jacobi1d.compile()
    err |= jacobi2d.compile()
    err |= kmeans.compile()
    # err |= lu.compile()
    # err |= lenet.compile()
    # err |= matrixmultiplication.compile()
    err |= matrixtranspose.compile()
    # err |= maxpooling.compile()
    # err |= mineva.compile()
    err |= mis.compile()
    # err |= mm2.compile()
    # err |= mm3.compile()
    # err |= mvt.compile()
    # err |= nbody.compile()
    err |= pagerank.compile()
    # err |= relu.compile()
    err |= simpleconvolution.compile()
    err |= shocreduction.compile()
    err |= spmv.compile()
    # err |= sssp.compile()
    err |= stencil2d.compile()
    err |= syrk.compile()
    err |= syr2k.compile()
    # err |= vgg16.compile()
    # err |= xor.compile()

    print(err)

if __name__ == '__main__':
    main()
