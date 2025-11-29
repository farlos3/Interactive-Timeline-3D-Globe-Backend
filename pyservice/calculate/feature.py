import numpy as np
from datetime import datetime
from .distance import gaussian

def normalize_date_with_gaussian(date_str, min_date, war_years=[(1914, 1918), (1939, 1945)], sigma=183):
    if isinstance(date_str, float) or isinstance(date_str, int):
        date_str = str(int(date_str))
    date_obj = datetime.strptime(date_str, "%Y-%m-%d")

    if isinstance(min_date, str):
        min_date = datetime.strptime(min_date, "%Y-%m-%d")

    days_since_min = (date_obj - min_date).days

    war_effect = 0
    for start_year, end_year in war_years:
        war_center_date = datetime((start_year + end_year) // 2, 7, 1)
        war_center_days = (war_center_date - min_date).days
        war_effect += gaussian(days_since_min, war_center_days, sigma)

    return days_since_min + war_effect

def create_feature_vector(lat, lon, date_str, min_date, war_years=[(1914, 1918), (1939, 1945)], sigma=300):
    time_value = normalize_date_with_gaussian(date_str, min_date, war_years, sigma)
    return np.array([lat, lon, time_value])