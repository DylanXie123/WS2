from datetime import datetime
import matplotlib.pyplot as plt
import numpy as np
from result import result, xlim, ylim

dt = datetime.now()
result = np.array(result)

x = [float(x) for x in result[:, 0]]
y = [float(y) for y in result[:, 1]]

plt.scatter(x, y, s=50.0, c=result[:, 2])
plt.axis([0, xlim, 0, ylim])
plt.savefig(f'./results/{dt.year}-{dt.month}-{dt.day}_{dt.hour}{dt.minute}.eps')
plt.show()
