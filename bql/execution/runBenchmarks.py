#!/usr/bin/python

import os
import sys
import numpy

results = {}

for i in range(10):
    output = os.popen("go test -run=XYZ -bench=.").read()
    lines = output.strip().split("\n")
    header = lines[0]
    if not header == "PASS":
        print "failed to run tests"
        sys.exit(1)
    for benchmark in lines[1:-1]:
        name = benchmark.split()[0]
        duration = benchmark.split()[-2]
        if name in results:
            results[name].append(int(duration))
        else:
            results[name] = [int(duration)]

for name in sorted(results.keys()):
    values = results[name]
    print name + "\t" + str(numpy.mean(values)) + "\t" + str(numpy.std(values))
