#!/usr/bin/env python3

# Converts the output from the benchmark into a html document.
# Usage:
#  ./hashmap/benchmark.sh | tee /tmp/bench.out
#  ./hashmap/chart.py < /tmp/bench.out > /tmp/chart.html
#  firefox /tmp/chart.html

import sys
from collections import defaultdict

mapping = defaultdict(lambda: defaultdict(list))

#
# collect metric values
#
for line in sys.stdin:
    lineRaw = line.split('\t')
    if len(lineRaw) != 8 or not lineRaw[0].startswith('Benchmark'):
        continue
    firstRaw = lineRaw[0].split('/')
    benchName = firstRaw[0]
    annotationList = firstRaw[1].split('-')
    mapName = annotationList[0]
    n = int(annotationList[1])
    time_ms = int(lineRaw[2].strip().split(' ')[0]) / (1000 * 1000)
    load = 'load=' + lineRaw[4].strip().split(' ')[0]
    mapping[benchName][mapName].append((n,time_ms,load))
    if "BenchmarkRandomFullInsertsInsertsU64" in benchName:
        memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
        mapping["MemoryConsumption"][mapName].append((n,memory_bytes,load))


#
# print html document
#

# print html header
print('''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Hashmap Benchmark</title>
    <script src='https://cdn.plot.ly/plotly-2.20.0.min.js'></script>
</head>
<body>
    <h1 align="center">Hashmap Benchmark</h1>
    <hr>
''')
      
for benchmark in sorted(mapping):
    b = mapping[benchmark]
    print("<div id='"+benchmark+"'><script>")
    names = []
    y_naming = "time (ms)"
    if benchmark == "MemoryConsumption":
        y_naming = "memory (MB)"
    for mapName in b:
        points = b[mapName]
        x_values = map(lambda x: x[0], points)
        y_values = map(lambda x: x[1], points)
        load_values = map(lambda x: x[2], points)
        name = benchmark+'_'+mapName
        names.append(name)
        print('var ' + name + ' = {')
        print("name: '" + mapName + "',")
        print('    x: ', list(x_values), ',')
        print('    y: ', list(y_values), ',')
        print('    text: ', list(load_values), ',')
        print('''   mode: 'lines+markers', type: 'scatter'
};''')
    print("var data_" + benchmark, "=", '[%s]' % ', '.join(map(str, names)), ";")
    print("var layout_" + benchmark + " = {title:'" + benchmark + "', xaxis: {title: 'number of entries in hash table'},yaxis: {title: '" + y_naming + "'}};");
    print("Plotly.newPlot('" + benchmark + "', data_"+ benchmark + ", layout_" + benchmark + ");"),
    print("</script></div><hr>")

# print rest of body
print("</body>")
