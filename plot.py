import matplotlib.pyplot as plt
import numpy as np
from result import result, size

result = np.array(result)

x = [float(x) for x in result[:, 0]]
y = [float(y) for y in result[:, 1]]

plt.scatter(x, y, s=50.0, c=result[:, 2])
plt.axis([0, size, 0, size])
plt.show()
