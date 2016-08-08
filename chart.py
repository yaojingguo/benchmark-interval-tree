#!/usr/bin/env python

import matplotlib.pyplot as plt
labels=["2", "4", "8", "16", "32", "64", "128", "256"]
x=[1, 2, 3, 4, 5, 6, 7, 8]

plt.xticks(x, labels, rotation="vertical")

size=8
insert=[2.40,  1.39,  1.12,  1.01,  0.95,  0.93,  0.89,   0.90]
fastInsert=[2.17,  1.26,  1.03,  1.06,  0.89,  0.87,  0.83,   0.83]
delete=[2.43,  1.56,  1.27,  1.18,  1.11,  1.07,  1.06,   1.08]
get=[4.66,  3.58,  3.39,  3.40,  3.66,  5.15,  5.87,  10.08]

all=[]
for i in range(size):
  all.append((insert[i] + fastInsert[i] + delete[i] + get[i]) / 4)

print all

plt.plot(x, insert, color="red", marker="s", label="Insert")
plt.plot(x, fastInsert, color="blue", marker="s", label="FastInsert")
plt.plot(x, delete, color="green", marker="s", label="Delete")
plt.plot(x, get, color="y", marker="s",label="Get")
plt.plot(x, all, color="orange", marker="s",label="Average")
plt.legend(loc="upper center")

# Pad margins so that markers don"t get clipped by the axes
plt.margins(0.2)
# Tweak spacing to prevent clipping of tick-labels
plt.subplots_adjust(bottom=0.15)
plt.show()
