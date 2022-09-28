from helper_functions.util import *
import traceback
import re


def getL2TlbStats(lines):
    accessResults = getL2TlbAccessResults(lines)
    accessResultRates = getAccessResultRates(accessResults)
    accessResults.update(accessResultRates)
    mpki = getL2TlbMpki(lines, accessResults)
    accessResults.update(mpki)
    accessLat = getAccessLat(lines, accessResults)
    accessResults.update(accessLat)
    missLat = getMissLat(lines, accessResults)
    accessResults.update(missLat)
    qLen = getTLBQLen(lines, accessResults)
    accessResults.update(qLen)
    mshrOccupancy = getTLBMSHRStats(lines)
    accessResults.update(mshrOccupancy)
    accessResults.update(getL2ImbalanceMetric(lines, accessResults))
    accessResults.update(getL2TLBStalledReqs(lines))
    return {'l2-{}'.format(k): v for k, v in accessResults.items()}


def getL2TlbAccessResults(lines):
    l2TlbStats = list(map(
        lambda x: {'tlb-hit': 0, 'tlb-miss': 0, 'tlb-mshr-hit': 0}, range(4)))
    rtuStats = [[0]*4 for _ in range(4)]
    for line in lines.split('\n'):
        try:
            splitLine = line.split(',')
            componentName = splitLine[1].strip()
            val = int(splitLine[-1][:-1].strip().split('.')[0])
            chiplet = getChipletFromComponent(componentName)
            if 'L2TLB' in line:
                if 'tlb-mshr-hit' in line:
                    addToArray(l2TlbStats, chiplet, 'tlb-mshr-hit', val)
                elif 'tlb-hit' in line:
                    addToArray(l2TlbStats, chiplet, 'tlb-hit', val)
                elif 'tlb-miss' in line:
                    addToArray(l2TlbStats, chiplet, 'tlb-miss', val)
            elif 'RTU, inter_chiplet_tlb_traffic' in line:
                sourceString, _, destString = splitLine[2].split(':')[1:]
                source = int(sourceString.split('_')[1])
                assert source == chiplet
                dest = int(destString.split('_')[1])
                val = int(splitLine[-1][:-1].strip().split('.')[0])
                rtuStats[source][dest] = val
        except Exception as e:
            # traceback.print_exc()
            continue
    toReturn = {'tlb-hit': 0, 'tlb-miss': 0, 'tlb-mshr-hit': 0, 'tlb-access': 0}
    for i in range(4):
        assert rtuStats[i][i] == 0
        outgoingReqs = sum(rtuStats[i])
        incomingReqs = sum(map(lambda x: rtuStats[x][i], range(4)))
        for m in ('hit', 'miss', 'mshr-hit'):
            mKey = 'tlb-{}'.format(m)
            mVal = l2TlbStats[i][mKey]
            toReturn['tlb-{}-{}'.format(i, m)] = mVal
            toReturn['tlb-{}-access'.format(i)] = toReturn.pop('tlb-{}-access'.format(i), 0) + mVal
            toReturn[mKey] += mVal
            toReturn['tlb-access'] += mVal
        localReqs = sum(map(
            lambda x: l2TlbStats[i][x], ('tlb-hit', 'tlb-miss', 'tlb-mshr-hit'))) - incomingReqs
        rtuStats[i][i] = localReqs
    return toReturn


def getAccessResultRates(accessResults):
    for i in range(4):
        print(accessResults['tlb-{}-miss'.format(i)], accessResults['tlb-{}-access'.format(i)])
        print("tlb miss rate:", div(float(accessResults['tlb-{}-miss'.format(i)]), accessResults['tlb-{}-access'.format(i)])*100)
    totalAccesses = float(
        sum(accessResults['tlb-{}'.format(m)] for m in ('hit', 'miss', 'mshr-hit')))
    return {'tlb-{}-rate'.format(m): div(accessResults['tlb-{}'.format(m)], totalAccesses)*100
            for m in ('hit', 'miss', 'mshr-hit')}


def getL2TlbMpki(lines, accessResults):
    regex = re.compile("^.*?inst_count.*?$", re.MULTILINE)
    instCountLines = regex.findall(lines)
    instCount = sum(map(lambda x: float(x.split(',')[-1]), instCountLines))
    return {'tlb-mpki': accessResults['tlb-miss']/instCount*1000}


def getAccessLat(lines, accessResults):
    latSum = 0.0
    totReqs = 0
    regex = re.compile(
            "^.*?chiplet_0.*?L2TLB, req_average_latency, .*?$", re.MULTILINE)
    avgLatLines = regex.findall(lines)
    avgLat = [0.0]*4
    for l in avgLatLines:
        splitLine = l.split(',')
        l2Tlb = getChipletFromComponent(splitLine[1].strip())
        val = float(splitLine[-1]) 
        avgLat[l2Tlb] = val
    for i in range(4):
        reqs = sum(accessResults['tlb-{}-{}'.format(i, m)]
                   for m in ('hit', 'miss', 'mshr-hit'))
        latSum += (avgLat[i]*reqs)
        totReqs += reqs
    return {'tlb-access-lat': asCycles(div(latSum, totReqs))}


def getMissLat(lines, accessResults):
    latSum = 0.0
    totReqs = 0
    regex = re.compile(
            "^.*?chiplet_0.*?L2TLB, down_req_average_latency, .*?$", re.MULTILINE)
    avgLatLines = regex.findall(lines)
    avgLat = [0.0]*4
    for l in avgLatLines:
        splitLine = l.split(',')
        l2Tlb = getChipletFromComponent(splitLine[1].strip())
        val = float(splitLine[-1]) 
        avgLat[l2Tlb] = val
    for i in range(4):
        reqs = accessResults['tlb-{}-miss'.format(i)]
        latSum += (avgLat[i]*reqs)
        totReqs += reqs
    return {'tlb-miss-lat': asCycles(div(latSum, totReqs))}

def getChipletFromComponentAlternate(componentName): 
    chiplet = int(componentName[-1])
    return chiplet


def getTLBMSHRStats(lines):
    metrics = ['average_mshr_len', 'average_mshr_len_g0', 'average_mshr_uniq_len', 'average_mshr_uniq_len_g0']
    tlbMSHRDict = dict()
    for metric in metrics:
        tlbMSHRDict[metric] = dict()
    for line in lines.split('\n'):
        try:
            splitLine = line.split(',')
            componentName = splitLine[1].strip()
            val = int(splitLine[-1][:-1].strip().split('.')[0])
            chiplet = getChipletFromComponentAlternate(componentName)
        except:
            continue
        if 'L2TLB' in line:
            for metric in metrics:
                # the comma is important to avoid a similarly prefixed stat: 
                regex = metric + ','
                if regex in line:
                    val = float(splitLine[-1][:-1].strip())
                    tlbMSHRDict[metric][chiplet]  = val
    returnDict = dict()
    for metric in metrics:
        returnDict['tlb-{}'.format(metric)] = 0.0
        for chiplet in range(4):
            returnDict['tlb-{}-{}'.format(metric, chiplet)] = getVal(tlbMSHRDict[metric], chiplet)
            if 'uniq' in metric and 'g0' in metric:
                print(metric, chiplet, getVal(tlbMSHRDict[metric], chiplet))
            returnDict['tlb-{}'.format(metric)] += getVal(tlbMSHRDict[metric], chiplet)
    return returnDict



def getTLBQLen(lines, accessData):
    qLen = [0]*4
    for line in lines.split('\n'):
        if 'average_buf_len_g0' in line:
            splitLine = line.split(',')
            tlb = int(splitLine[1].strip()[-1])
            val = float(splitLine[-1][:-1].strip())
            qLen[tlb] = val
    totQLen = 0.0
    for i in range(4):
        print("queue length", i, qLen[i])
        # print(i, accessData['tlb-{}-access'.format(i)])
        totQLen += qLen[i]*accessData['tlb-{}-access'.format(i)]
    return {'tlb-q-len': div(totQLen, accessData['tlb-access'])}



def getL2TLBStalledReqs(lines):
    stalledReqs = 0
    totalReqs = 0
    for line in lines.split('\n'):
        if 'stalled-l2-tlb-req-count' in line:
            splitLine = line.split(',')
            # tlb = int(splitLine[1].strip()[-1])
            val = float(splitLine[-1][:-1].strip())
            stalledReqs += val
            # qLen[tlb] = val
        elif 'l2-tlb-req-count' in line:
            splitLine = line.split(',')
            val = float(splitLine[-1][:-1].strip())
            totalReqs += val
    print("######################", div(stalledReqs, totalReqs))
    return {'tlb-frac-reqs-stalled': div(stalledReqs, totalReqs)}


def getL2ImbalanceMetric(lines, accessResults):
    toReturn = dict()
    for line in lines.split('\n'):
        if 'imbalance' in line and 'L2TLB' in line:
            splitLine = line.split(',')
            componentName = splitLine[1].strip()
            val = float(splitLine[-1][:-1].strip())
            chiplet = getChipletFromComponent(componentName)
            toReturn['{}-{}'.format(splitLine[2].strip(), str(chiplet))] = val 
    imbalanceMetric = 0.0
    imbalanceCount = 0
    for i in range(4):
        imbalance = toReturn['imbalance-{}'.format(i)]*100
        count = toReturn['imbalance-count-{}'.format(i)]
        accesses = accessResults['tlb-{}-access'.format(i)]
        try:
            imbalanceMetric += imbalance * float(count) / accesses
        except:
            imbalanceMetric += 0
        imbalanceCount += count
        # print(imbalance, count, accesses, imbalanceMetric)
    toReturn['tlb-imbalance'] = imbalanceMetric
    toReturn['tlb-imbalance-count'] = imbalanceCount
    return toReturn
