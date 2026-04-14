import csv
import matplotlib.pyplot as plt
from collections import defaultdict

data = defaultdict(lambda: {"x": [], "y": []})
print(data)

# Здесь можно запунить файл на тот, который вы хотите посмотреть
with open("program_model1.csv", "r") as f:
    reader = csv.DictReader(f)
    
    for i, row in enumerate(reader):
        body = int(row["body"])
        x = float(row["x"])
        y = float(row["y"])
        
        data[body]["x"].append(x)
        data[body]["y"].append(y)

        print(i)

print(data)

plt.figure(figsize=(8, 8))



for body, coords in data.items():
    plt.plot(coords["x"], coords["y"], label=f"Body {body}")

# ===== финальные позиции =====
for body, coords in data.items():
    plt.scatter(coords["x"][-1], coords["y"][-1])

plt.xlabel("X")
plt.ylabel("Y")
plt.title("Orbit simulation")
plt.legend()
plt.axis("equal")
plt.grid()

plt.show()

