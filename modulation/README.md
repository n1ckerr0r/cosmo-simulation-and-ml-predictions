# Короткое описание проекта

в internal/ находятся вычисления шага, подсчет сил и логика работы цикла программы
в cmd/ запуск программы
в data/ выходные данные

## Демонстрационные сценарии

`go run ./cmd/simulator` сохраняет три примера с одним и более полным оборотом:

- `data/output_mercury.csv`
- `data/output_earth.csv`
- `data/output_mars.csv`

`python3 temp_python_script.py` строит простой 3D-график `data/orbits_3d.png` по этим CSV.

Тесты запускаются командой `go test ./...`.
