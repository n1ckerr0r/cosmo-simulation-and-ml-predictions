from pathlib import Path
import warnings

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score

warnings.filterwarnings("ignore")

try:
    from catboost import CatBoostRegressor
    CATBOOST_AVAILABLE = True
except ImportError:
    CATBOOST_AVAILABLE = False


RESULTS_DIR = Path("../results")
OUTPUT_PRED_FILE = RESULTS_DIR / "ml_pred.csv"

INPUT_FILES = [
    RESULTS_DIR / "program_model1.csv",
    RESULTS_DIR / "program_model2.csv",
    RESULTS_DIR / "program_model3.csv",
]


def load_simulation_file(file_path: Path, sim_id: str) -> pd.DataFrame:
    df = pd.read_csv(file_path)

    required_columns = {"step", "body", "x", "y", "z", "energy"}
    missing = required_columns - set(df.columns)
    if missing:
        raise ValueError(f"В файле {file_path.name} нет колонок: {missing}")

    df = df.copy()
    df["sim_id"] = sim_id
    return df


def build_features(df: pd.DataFrame) -> pd.DataFrame:
    """
    Для каждого тела в каждой симуляции строим признаки:
    - текущие координаты
    - координаты на 1 и 2 шага назад
    - разности (условно скорости)
    Цель:
    - координаты на следующем шаге
    """
    df = df.sort_values(["sim_id", "body", "step"]).copy()

    group_cols = ["sim_id", "body"]

    for col in ["x", "y", "z"]:
        df[f"{col}_lag1"] = df.groupby(group_cols)[col].shift(1)
        df[f"{col}_lag2"] = df.groupby(group_cols)[col].shift(2)

    df["vx_approx"] = df["x"] - df["x_lag1"]
    df["vy_approx"] = df["y"] - df["y_lag1"]
    df["vz_approx"] = df["z"] - df["z_lag1"]

    df["target_x"] = df.groupby(group_cols)["x"].shift(-1)
    df["target_y"] = df.groupby(group_cols)["y"].shift(-1)
    df["target_z"] = df.groupby(group_cols)["z"].shift(-1)

    df = df.dropna().reset_index(drop=True)

    return df


def train_test_split_by_time(df: pd.DataFrame, test_size: float = 0.2):
    """
    Делим по времени
    """
    max_step = df["step"].max()
    split_step = int(max_step * (1 - test_size))

    train_df = df[df["step"] <= split_step].copy()
    test_df = df[df["step"] > split_step].copy()

    return train_df, test_df, split_step


def evaluate_regression(y_true: np.ndarray, y_pred: np.ndarray, model_name: str):
    mae = mean_absolute_error(y_true, y_pred)
    rmse = np.sqrt(mean_squared_error(y_true, y_pred))
    r2 = r2_score(y_true, y_pred)

    print(f"\nМетрики для {model_name}:")
    print(f"MAE  = {mae:.6e}")
    print(f"RMSE = {rmse:.6e}")
    print(f"R2   = {r2:.6f}")

    return {"model": model_name, "mae": mae, "rmse": rmse, "r2": r2}


def main():
    all_dfs = []
    for file_path in INPUT_FILES:
        if not file_path.exists():
            print(f"[WARNING] Файл не найден: {file_path}")
            continue
        if file_path.stat().st_size == 0:
            print(f"[WARNING] Файл пустой: {file_path}")
            continue

        sim_id = file_path.stem
        df = load_simulation_file(file_path, sim_id)
        all_dfs.append(df)

    if not all_dfs:
        raise RuntimeError("Нет входных данных для обучения.")

    raw_df = pd.concat(all_dfs, ignore_index=True)
    print("Общая форма сырых данных:", raw_df.shape)

    ml_df = build_features(raw_df)
    print("Форма после построения признаков:", ml_df.shape)

    feature_cols = [
        "sim_id",
        "step",
        "body",
        "energy",
        "x", "y", "z",
        "x_lag1", "y_lag1", "z_lag1",
        "x_lag2", "y_lag2", "z_lag2",
        "vx_approx", "vy_approx", "vz_approx",
    ]

    target_cols = ["target_x", "target_y", "target_z"]

    data = ml_df[feature_cols + target_cols].copy()

    data["sim_id"] = data["sim_id"].astype("category").cat.codes

    train_df, test_df, split_step = train_test_split_by_time(data, test_size=0.2)

    print(f"Граница train/test по step: {split_step}")
    print("Train shape:", train_df.shape)
    print("Test shape:", test_df.shape)

    X_train = train_df[feature_cols]
    X_test = test_df[feature_cols]

    y_train = train_df[target_cols]
    y_test = test_df[target_cols]

    linear_model = LinearRegression()
    linear_model.fit(X_train, y_train)
    linear_pred = linear_model.predict(X_test)

    linear_metrics = evaluate_regression(y_test.values, linear_pred, "LinearRegression")

    best_model_name = "LinearRegression"
    best_pred = linear_pred
    best_metrics = linear_metrics

    if CATBOOST_AVAILABLE:
        cat_models = {}
        cat_pred_parts = []

        for target in target_cols:
            y_train_target = y_train[target]
            unique_count = y_train_target.nunique()

            if unique_count == 1:
                const_value = y_train_target.iloc[0]

                preds = np.full(shape=len(X_test), fill_value=const_value)
                cat_models[target] = None
            else:
                model = CatBoostRegressor(
                    iterations=300,
                    depth=6,
                    learning_rate=0.05,
                    loss_function="RMSE",
                    eval_metric="RMSE",
                    verbose=0,
                    random_seed=42
                )
                model.fit(X_train, y_train_target)
                preds = model.predict(X_test)

                cat_models[target] = model

            cat_pred_parts.append(preds)

        cat_pred = np.column_stack(cat_pred_parts)
        cat_metrics = evaluate_regression(y_test.values, cat_pred, "CatBoostRegressor")

        if cat_metrics["rmse"] < best_metrics["rmse"]:
            best_model_name = "CatBoostRegressor"
            best_pred = cat_pred
            best_metrics = cat_metrics
    else:
        print("\n[INFO] CatBoost не установлен. Пока обучена только LinearRegression.")
        print("Установка: pip install catboost")

    pred_df = test_df[["sim_id", "step", "body"]].copy()
    pred_df["x_true"] = y_test["target_x"].values
    pred_df["y_true"] = y_test["target_y"].values
    pred_df["z_true"] = y_test["target_z"].values

    pred_df["x_pred"] = best_pred[:, 0]
    pred_df["y_pred"] = best_pred[:, 1]
    pred_df["z_pred"] = best_pred[:, 2]
    pred_df["model_name"] = best_model_name

    pred_df.to_csv(OUTPUT_PRED_FILE, index=False)

    print(f"\nЛучший результат: {best_model_name}")
    print(f"Файл с предсказаниями сохранен: {OUTPUT_PRED_FILE}")

    print("\nИтоговые метрики лучшей модели:")
    print(f"MAE  = {best_metrics['mae']:.6e}")
    print(f"RMSE = {best_metrics['rmse']:.6e}")
    print(f"R2   = {best_metrics['r2']:.6f}")


if __name__ == "__main__":
    main()