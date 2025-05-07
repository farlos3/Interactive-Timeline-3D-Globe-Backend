from .distance import find_weighted_dist

def brute_force(arr):
    if len(arr) < 2:
        return float('inf'), None, None
    min_dist = find_weighted_dist(arr[0], arr[1])
    pair = (arr[0], arr[1])
    for i in range(len(arr)):
        for j in range(i+1, len(arr)):
            dist = find_weighted_dist(arr[i], arr[j])
            if dist < min_dist:
                min_dist = dist
                pair = (arr[i], arr[j])
    return min_dist, pair[0], pair[1]

def merge_sort(arr, coord):
    if len(arr) <= 1:
        return arr
    mid = len(arr) // 2
    left = merge_sort(arr[:mid], coord)
    right = merge_sort(arr[mid:], coord)
    return merge(left, right, coord)

def merge(A, B, coord):
    result = []
    i = j = 0
    while i < len(A) and j < len(B):
        a_val = float(A[i][0][coord]) if isinstance(A[i], tuple) else float(A[i][coord])
        b_val = float(B[j][0][coord]) if isinstance(B[j], tuple) else float(B[j][coord])
        if a_val <= b_val:
            result.append(A[i])
            i += 1
        else:
            result.append(B[j])
            j += 1
    result.extend(A[i:])
    result.extend(B[j:])
    return result

def closest_pair(Px, Py):
    n = len(Px)
    if n <= 3:
        return brute_force(Px)
    mid = n // 2
    Qx, Rx = Px[:mid], Px[mid:]
    mid_x = Px[mid][0]
    Qy = [p for p in Py if p[0] < mid_x]
    Ry = [p for p in Py if p[0] >= mid_x]
    d1, a1, b1 = closest_pair(Qx, Qy)
    d2, a2, b2 = closest_pair(Rx, Ry)
    d, A, B = (d1, a1, b1) if d1 < d2 else (d2, a2, b2)
    strip = [p for p in Py if abs(p[0] - mid_x) < d]
    for i in range(len(strip)):
        for j in range(i+1, min(i+7, len(strip))):
            d_strip = find_weighted_dist(strip[i], strip[j])
            if d_strip < d:
                d, A, B = d_strip, strip[i], strip[j]
    return d, A, B
