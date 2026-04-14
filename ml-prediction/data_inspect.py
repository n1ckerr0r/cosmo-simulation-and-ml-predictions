from pathlib import Path
import pandas as pd

RESULTS_DIR = Path("../results")
csv_files = sorted(RESULTS_DIR.glob("program_model*.csv"))

if not csv_files:
    raise FileNotFoundError("Не найдены файлы program_model*.csv в папке results")

print("Найдены файлы:")
for file in csv_files:
    print("-", file.name)

for file in csv_files:
    print("\n" + "=" * 80)
    print(f"Файл: {file.name}")
    print("=" * 80)

    if not file.exists():
        print("Файл не существует")
        continue

    size = file.stat().st_size
    print(f"Размер файла: {size} байт")

    with open(file, "r", encoding="utf-8", errors="replace") as f:
        first_lines = []
        for _ in range(5):
            line = f.readline()
            if not line:
                break
            first_lines.append(repr(line))

    print("Первые строки файла:")
    if first_lines:
        for line in first_lines:
            print(line)
    else:
        print("Файл пустой")

    try:
        df = pd.read_csv(file)
        print("Размерность таблицы:", df.shape)
        print("Колонки:", df.columns.tolist())
        print(df.head())
    except Exception as e:
        print("Ошибка чтения pandas:", type(e).__name__, str(e))