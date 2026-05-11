import csv
from collections import defaultdict
from pathlib import Path

import matplotlib.pyplot as plt


ROOT = Path(__file__).resolve().parent
OUTPUT_IMAGE = ROOT / "data" / "orbits_3d.png"
FILES = [
    ("Mercury", ROOT / "data" / "output_mercury.csv"),
    ("Earth", ROOT / "data" / "output_earth.csv"),
    ("Mars", ROOT / "data" / "output_mars.csv"),
]


def read_body_coords(path, target_body=1):
    data = defaultdict(lambda: {"x": [], "y": [], "z": []})

    with path.open("r", newline="") as f:
        reader = csv.DictReader(f)

        for row in reader:
            body = int(row["body"])
            if body != target_body:
                continue

            data[body]["x"].append(float(row["x"]))
            data[body]["y"].append(float(row["y"]))
            data[body]["z"].append(float(row["z"]))

    return data[target_body]


def main():
    fig = plt.figure(figsize=(9, 8))
    ax = fig.add_subplot(111, projection="3d")

    ax.scatter([0], [0], [0], color="gold", s=80, label="Sun")

    for name, path in FILES:
        if not path.exists():
            raise SystemExit(f"Missing {path}. Run: cd modulation && go run ./cmd/simulator")

        coords = read_body_coords(path)
        ax.plot(coords["x"], coords["y"], coords["z"], label=name)
        ax.scatter(coords["x"][-1], coords["y"][-1], coords["z"][-1], s=30)

    ax.set_xlabel("X")
    ax.set_ylabel("Y")
    ax.set_zlabel("Z")
    ax.set_title("3D orbit simulation")
    ax.legend()
    ax.grid(True)
    ax.set_box_aspect((1, 1, 0.45))

    plt.tight_layout()
    plt.savefig(OUTPUT_IMAGE, dpi=160)
    plt.show()


if __name__ == "__main__":
    main()
