import math
import numpy as np

def gaussian(x, mu, sigma):
    return (1 / (sigma * np.sqrt(2 * np.pi))) * np.exp(-0.5 * ((x - mu) / sigma) ** 2)

def find_weighted_dist(p, q):
    weights = [1, 1, 5]
    return math.sqrt(sum(weights[i] * (p[i] - q[i]) ** 2 for i in range(len(p))))