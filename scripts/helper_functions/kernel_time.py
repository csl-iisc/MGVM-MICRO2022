import re


def getKernelTime(lines):
    regex = re.compile("^.*?GPU.\.CommandProcessor, kernel_time.*?$", re.MULTILINE)
    foundLines = regex.findall(lines)
    # assert(len(foundLines) == 1)
    val = float(foundLines[0].split(',')[-1][:-1].strip())
    toReturn = {'kernel_time': val}
    for i, line in enumerate(foundLines[1:]):
        toReturn['kernel_time_{}'.format(i)] = float(line.split(',')[-1].strip())
    # if len(toReturn) < 10:
    #     print(toReturn)

    ## NEWW
    if val == 0.0 :
        regex = re.compile("^.*?GPU.\.CommandProcessor, kernel_time \(force stop\).*?$", re.MULTILINE)
        foundLines = regex.findall(lines)
        print(foundLines)
        val = float(foundLines[0].split(',')[-1][:-1].strip())
        toReturn['kernel_time'] = val

    return toReturn
