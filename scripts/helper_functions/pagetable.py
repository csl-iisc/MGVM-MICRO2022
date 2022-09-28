from helper_functions.util import *

def getPageWalkStats(lines):
    # 4 chiplets to 4 chiplets with 3 types (hit, miss, mshr)
    latencies = [[[[0 for l in range(4)] for k in range(3)] for j in range(4)] for i in range(4)]
    num = [[[[0 for l in range(4)] for k in range(3)] for j in range(4)] for i in range(4)]

    # cast (convert?) lines to normal type
    lines  = lines.splitlines()

    for line in lines: 
        splitLine = line.split(',')
        try:
            componentName = splitLine[1].strip()
            chiplet = int(getChipletFromComponent(componentName))
        except Exception as e:
            continue
        if 'MMU' in line and 'read' in line:
            toComponent = int(splitLine[2].split('-')[1])
            level = int(splitLine[2].split('-')[3])
            val = splitLine[-1][:-1].strip()
            if 'read-hit-latency' in line:
                latencies[chiplet][toComponent][0][level] += float(val)
            if 'read-hit-num' in line:
                num[chiplet][toComponent][0][level] += int(float(val))
            if 'read-miss-latency' in line:
                latencies[chiplet][toComponent][1][level] += float(val)
            if 'read-miss-num' in line:
                num[chiplet][toComponent][1][level] += int(float(val))
            if 'read-mshr-hit-latency' in line:
                latencies[chiplet][toComponent][2][level] += float(val)
            if 'read-mshr-hit-num' in line:
                num[chiplet][toComponent][2][level] += int(float(val))

    localLatency = 0
    localReqs = 0
    remoteLatency = 0
    remoteReqs = 0

    cacheHitLatency = 0
    cacheHitNum = 0
    cacheMissLatency = 0
    cacheMissNum = 0
    cacheMSHRLatency = 0
    cacheMSHRNum = 0

    for i in range(4):
        for j in range(4):
            for k in range(3):
                for l in range(4):
                    if k == 0: 
                        cacheHitLatency += latencies[i][j][k][l] * num[i][j][k][l]
                        cacheHitNum += num[i][j][k][l]
                    if k == 1:
                        cacheMissLatency += latencies[i][j][k][l] * num[i][j][k][l]
                        cacheMissNum += num[i][j][k][l]
                    if k == 2:
                        cacheMSHRLatency += latencies[i][j][k][l] * num[i][j][k][l]
                        cacheMSHRNum += num[i][j][k][l]

    for i in range(4):
        for j in range(4):
            for k in range(3):
                for l in range(4):
                    if i != j: 
                        remoteLatency += latencies[i][j][k][l] * num[i][j][k][l]
                        remoteReqs += num[i][j][k][l]
                    else:
                        localLatency += latencies[i][j][k][l] * num[i][j][k][l]
                        localReqs += num[i][j][k][l]

    toReturn = dict()

    # print cacheHitLatency
    # print cacheHitNum
    # print remoteLatency/remoteReqs
    # print localLatency/localReqs

    toReturn.update({'pw_L2cache_hit_lat':div(cacheHitLatency,cacheHitNum)})
    toReturn.update({'pw_L2cache_miss_lat':div(cacheMissLatency,cacheMissNum)})
    toReturn.update({'pw_L2cache_mshr_lat':div(cacheMSHRLatency,cacheMSHRNum)})
    toReturn.update({'pw_L2cache_hit_num':cacheHitNum})
    toReturn.update({'pw_L2cache_miss_num':cacheMissNum})
    toReturn.update({'pw_L2cache_mshr_num':cacheMSHRNum})

    toReturn.update({'pw_remote_lat':div(remoteLatency, remoteReqs)})
    toReturn.update({'pw_remote_num':remoteReqs})
    toReturn.update({'pw_local_lat':div(localLatency, localReqs)})
    toReturn.update({'pw_local_num':localReqs})

    return toReturn

