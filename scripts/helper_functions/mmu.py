from helper_functions.util import *
# pwReqCols = map(lambda x: 'pw-req-to-chiplet-{}'.format(x), range(4))
# dataReqCols = map(lambda x: 'data-req-to-chiplet-{}'.format(x), range(4))
# pwcHitCols = ['pwc-miss'] + map(lambda x: 'pwc-hit-{}'.format(x), (1, 2, 3))
# cols = ['benchmark'] + map(lambda x: 'chiplet-{}-{}'.format(*x), itertools.product(range(4), pwReqCols +
#                            dataReqCols + pwcHitCols + ['active-walkers', 'pw-latency', 'rdma-outgoing-lat', 'rdma-incoming-lat']))
# cols += ['avg-pw-lat', 'avg-l2-miss-lat', 'avg-active-walkers']


def getMMUStats(lines):
    numPws = [0]*4
    avgPwLat = [0.0]*4
    l2MissLat = [0.0]*4
    localPwMemReqs = [0]*4
    remotePwMemReqs = [0]*4
    level3RemotePwMemReqs = [0]*4
    pwcHits = [[0]*4 for _ in range(4)]
    pwTraffic = [[0]*4 for _ in range(4)]
    dataTraffic = [[0]*4 for _ in range(4)]
    activeWalkers = [0.0]*4
    rdmaOutgoingLat = [0.0]*4
    rdmaIncomingLat = [0.0]*4
    rdmaOutgoingReqs = [0]*4

    for line in lines.split('\n'):
        splitLine = line.split(',')
        try:
            componentName = splitLine[1].strip()
            chiplet = getChipletFromComponent(componentName)
            val = splitLine[-1][:-1].strip()
        except Exception as e:
            continue
        if 'L2TLB' in line:
            if 'tlb-miss' in line:
                numPws[chiplet] = asInt(val)
            elif 'down_req_average_latency' in line:
                l2MissLat[chiplet] = float(val)
        if 'MMU' in line:
            if 'req_average_latency' in line:
                avgPwLat[chiplet] = float(val)
            if 'pwc-miss' in line:
                pwcHits[chiplet][0] = asInt(val)
            elif 'pwc-hit' in line:
                ptLevel = int(splitLine[2].split('-')[-1][-1])
                pwcHits[chiplet][ptLevel] = asInt(val)
            elif 'page_walk_req_local' in line:
                localPwMemReqs[chiplet] = asInt(val)
            elif 'page_walk_req_remote' in line:
                remotePwMemReqs[chiplet] = asInt(val)
                # print("remote mem accesses per walk:", asInt(val), numPws[chiplet], float(asInt(val))/numPws[chiplet])
            elif 'pw-level-3-remote-reqs' in line:
                level3RemotePwMemReqs[chiplet] = asInt(val)
            elif 'average active walkers' in line:
                activeWalkers[chiplet] = float(val)
        elif 'ChipRDMA' in line:
            if 'inter_chiplet_traffic' in line:
                sourceString, _, destString = splitLine[2].split(':')[1:]
                source = int(sourceString.split('_')[1])
                assert source == chiplet
                dest = int(destString.split('_')[1])
                (pwTraffic if 'PageAccess' in line else dataTraffic)[
                    source][dest] = asInt(val)
            elif 'outgoing_trans_latency' in line:
                rdmaOutgoingLat[chiplet] = float(val)
            elif 'incoming_trans_latency' in line:
                rdmaIncomingLat[chiplet] = float(val)
            elif 'outgoing_trans_count' in line:
                rdmaOutgoingReqs[chiplet] = asInt(val)
    toReturn = dict()
    pwLatSum = 0.0
    l2MissLatSum = 0.0
    activeWalkersSum = 0.0
    totPwMemReqs = 0
    pwRemoteMemReqs = 0
    pwLevel3RemoteMemReqs = 0
    for chiplet in range(4):
        outgoingReqs = sum(pwTraffic[chiplet]) + sum(dataTraffic[chiplet])
        # assert rdmaOutgoingReqs[chiplet] == outgoingReqs
        # assert numPws[chiplet] == sum(pwcHits[chiplet])
        # assert l2MissLat[chiplet] >= avgPwLat[chiplet]
        # assert pwTraffic[chiplet][chiplet] == 0
        # assert dataTraffic[chiplet][chiplet] == 0
        # assert remotePwMemReqs[chiplet] == sum(pwTraffic[chiplet])
        
        pwRemoteMemReqs += remotePwMemReqs[chiplet]
        pwLevel3RemoteMemReqs += level3RemotePwMemReqs[chiplet]
        
        pwMemReqs = localPwMemReqs[chiplet] + remotePwMemReqs[chiplet]
        # assert pwMemReqs == 4*pwcHits[chiplet][0] + 3*pwcHits[chiplet][1] + \
            # 2*pwcHits[chiplet][2] + 1*pwcHits[chiplet][3]
        totPwMemReqs += pwMemReqs
        # pwTraffic[chiplet][chiplet] = localPwMemReqs[chiplet]
        toReturn.update({'rdma-pw-req-{}-to-{}'.format(chiplet, i)                        : pwTraffic[chiplet][i] for i in range(4)})
        toReturn.update({'rdma-data-req-{}-to-{}'.format(chiplet, i)                        : dataTraffic[chiplet][i] for i in range(4)})
        # row += dataTraffic[chiplet]
        toReturn.update({'pwc-miss-chiplet-{}'.format(chiplet): pwcHits[chiplet][0]})
        toReturn.update(
            {'pwc-miss-chiplet-{}-level-{}'.format(chiplet, i): pwcHits[chiplet][i] for i in range(1, 4)})

        # row += pwcHits[chiplet]
        toReturn.update(
            {'mmu-active-walkers-{}'.format(chiplet): activeWalkers[chiplet]})
        toReturn.update({'mmu-pw-lat-{}'.format(chiplet): avgPwLat[0]})
        toReturn.update(
            {'rdma-outgoing-lat-{}'.format(chiplet): rdmaOutgoingLat[0]})
        toReturn.update(
            {'rdma-incoming-lat-{}'.format(chiplet): rdmaIncomingLat[0]})
        # row += [activeWalkers[chiplet], avgPwLat[chiplet],
        # rdmaOutgoingLat[chiplet], rdmaIncomingLat[chiplet]]
        pwLatSum += (avgPwLat[chiplet]*numPws[chiplet])
        l2MissLatSum += (l2MissLat[chiplet]*numPws[chiplet])
        activeWalkersSum += activeWalkers[chiplet]
    avgPwLat = div(pwLatSum, sum(numPws))
    avgL2MissLat = div(l2MissLatSum, sum(numPws))
    assert avgPwLat < avgL2MissLat or avgPwLat == avgL2MissLat == 0
    # row += [avgPwLat, avgL2MissLat, activeWalkersSum]
    toReturn.update({'mmu-pw-lat': asCycles(avgPwLat)})
    # toReturn.update({'pwc-miss-chiplet-{}'.format(chiplet): })
    toReturn.update({'mmu-active-walkers': activeWalkersSum})
    toReturn.update({'mmu-mem-reqs-per-pw': div(totPwMemReqs,float(sum(numPws)))})
    toReturn.update({'mmu-remote-mem-reqs-per-pw': div(pwRemoteMemReqs,float(sum(numPws)))})
    toReturn.update({'mmu-remote-mem-reqs-level-3': pwLevel3RemoteMemReqs})
    toReturn.update({'mmu-frac-remote-mem-reqs-level-3': div(pwLevel3RemoteMemReqs, float(pwRemoteMemReqs))})
    print('********************mmu-frac-remote-mem-reqs-level-3', toReturn['mmu-frac-remote-mem-reqs-level-3'])
    return toReturn
