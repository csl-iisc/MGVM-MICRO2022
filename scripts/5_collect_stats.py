#!/usr/bin/python3

import subprocess
import os
import sys
import argparse
import re
import csv
import itertools
import functools
import traceback

from helper_functions.kernel_time import *
from helper_functions.mmu import *
from helper_functions.pagetable import *
from helper_functions.l1_tlb import *
from helper_functions.l2_tlb import *
from helper_functions.benchmarks import *
from helper_functions.util import *

metrics = [
    'kernel_time',

    'l2-tlb-mpki',

    'mmu-pw-lat',

    'pw_local_num',
    'pw_remote_num',

]

locTypes = ['local', 'remote']
# accTypes = ['TLBHit', 'TLBMiss', 'TLBMshrHit']
accTypes = ['TLBHit']
#measureTypes = ['latency', 'num']
measureTypes = ['num']
for locType in locTypes:
    for accType in accTypes:
        for measureType in measureTypes:
            metrics.append(locType + '-' + accType + '-' + measureType+ '-' + 'total')

inputFilepath = sys.argv[1]
if "samples" in inputFilepath:
    print("Need the folder in which samples is placed!")
    exit()
if inputFilepath[-1] == "/":
    inputFilepath = inputFilepath[:-1]
benchmarksPath = inputFilepath + "/samples/"

allStats = dict()
configs = list()

def collectStats(inputFolders):
    for f in inputFolders:
        inputFilePath = extractPaths(f)
        print(inputFilePath)
        oneConfigStats = dict()
        splitFolderPath = f.split('/')
        print(splitFolderPath)
        config = "{}".format(splitFolderPath[0])  
        configs.append(config)
        for b in benchmarks:
            try:
                print(b)
                metricsFile = open(inputFilePath + b + "/metrics.csv", "r")
                lines = metricsFile.read()
                oneConfigStats[b] = dict()
                oneConfigStats[b].update(getKernelTime(lines))
                oneConfigStats[b].update(getL2TlbStats(lines))
                oneConfigStats[b].update(getMMUStats(lines))
                oneConfigStats[b].update(getPageWalkStats(lines))
                oneConfigStats[b].update(getPerChipletL1MissStats(lines))

                # oneConfigStats[b].update(getRTUCoalescingStats(lines))
                # # oneConfigStats[b].update(getTLBCoalescingStats(lines))
                # oneConfigStats[b].update(getTLBMSHRStats(lines))
                # oneConfigStats[b].update(getImbalanceStats(lines))
                # # oneConfigStats[b].update(getL2CacheStats(lines))
                # oneConfigStats[b].update(getTLBCoalescingStats(lines))
                # oneConfigStats[b].update(getPipelineStats(lines))
                # oneConfigStats[b].update(getDataAccessToChipletNumbers(lines))
                # oneConfigStats[b].update(getPageAccessToChipletNumbers(lines))
                # oneConfigStats[b].update(getAllL2CacheNumbers(lines))
                # oneConfigStats[b].update(getLocalPageWalksPerChiplet(lines))
                # # oneConfigStats[b].update(getNetworkLat(lines))
                # oneConfigStats[b].update(getSMStall(lines))
            except Exception as e:
                print(e)
                traceback.print_exc()
                continue
        allStats[config] = oneConfigStats
    outFile = open('results.csv', 'w')
    writeRow(outFile, ['benchmark'] + functools.reduce(lambda x, y: x+y, map(lambda x: ['', x] + configs,
                                          metrics)))
    for b in benchmarks:
        row = [b]
        for metric in metrics:
            row += ['', '']
            for config in configs:
                configStats = allStats[config]
                row += [getVal(configStats[b], metric) if b in configStats else 0]
        writeRow(outFile, row)


if __name__ == '__main__':
    collectStats(sys.argv[1:])

