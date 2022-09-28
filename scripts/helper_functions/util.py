import os

def asInt(string):
    return int(string.split('.')[0])

def getChipletFromComponent(componentName): 
    chiplet = int(componentName.split('.')[1].split('_')[-1])
    return chiplet


# def addMetric(component, metric, val, statsDict):
#     chiplet = getChipletFromComponent(component)
#     if component not in statsDict[chiplet]:
#         statsDict[chiplet][component] = {'combined-down-latency':0.0}
#     if metric not in statsDict[chiplet][component]:
#         statsDict[chiplet][component][metric] = 0.0
#     statsDict[chiplet][component][metric] += val


def addToArray(statsArray, index, metric, val):
    statsArray[index][metric] = val


def checkError(float1, float2):
    if abs(float1 - float2) > 10**(-10):
        raise Exception('Check error failed!')

def div(float1, float2):
    return 0.0 if not float2 else float1/float2

def writeRow(outFile, row):
    rowAsString = ','.join(map(lambda x: str(x), row))
    outFile.write('{}\n'.format(rowAsString))

def getVal(dictObj, key):
    return dictObj[key] if key in dictObj else 0.0

def addToDict(statsDict, componentName, metric, val):
    if componentName not in statsDict:
        statsDict[componentName] = dict()
    statsDict[componentName][metric] = val


def addToDictNoOverwrite(statsDict, componentName, metric, val):
    if componentName not in statsDict:
        statsDict[componentName] = dict()
    if metric not in statsDict[componentName]:
        statsDict[componentName][metric] = val
    else:
        statsDict[componentName][metric] += val


def extractPaths(inputFilePath):
    if "samples" in inputFilePath:
        print("Need the folder in which samples is placed!")
        exit()
    benchmarksPath = inputFilePath + "/samples/"

    return benchmarksPath

def asCycles(val):
    return val*10**9
