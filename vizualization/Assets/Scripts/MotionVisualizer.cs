using System;
using System.IO;
using System.Linq;
using System.Globalization;
using System.Collections.Generic;
using UnityEngine;

public class MotionVisualizer : MonoBehaviour
{
	[Header("Data")]
	public string fileName = "ml_pred.csv";

	[Tooltip("Выбери 1 или 2")]
	[Range(1, 2)]
	public int selectedSimId = 1;

	[Tooltip("Масштаб только для отображения")]
	public float visualScale = 1e-11f;

	[Header("Playback")]
	public float stepDuration = 0.02f;
	public int stepSkip = 10;
	public bool playOnStart = true;

	[Header("Visual")]
	public float trueSphereSize = 0.5f;
	public float predSphereSize = 0.35f;

	private readonly List<MotionRecord> records = new();
	private readonly Dictionary<string, GameObject> trueObjects = new();
	private readonly Dictionary<string, GameObject> predObjects = new();
	private readonly List<int> steps = new();

	private float timer = 0f;
	private int currentStepIndex = 0;
	private bool isPlaying = false;

	void Start()
	{
		LoadData();
		CreateObjectsForFirstStep();

		if (steps.Count > 0)
			ApplyStep(steps[0]);

		if (playOnStart)
			isPlaying = true;
	}

	void Update()
	{
		if (!isPlaying || steps.Count == 0)
			return;

		timer += Time.deltaTime;

		if (timer >= stepDuration)
		{
			timer = 0f;
			currentStepIndex += Mathf.Max(1, stepSkip);

			if (currentStepIndex >= steps.Count)
				currentStepIndex = 0;

			ApplyStep(steps[currentStepIndex]);
		}
	}

	void LoadData()
	{
		records.Clear();
		steps.Clear();

		string path = Path.Combine(Application.streamingAssetsPath, fileName);

		if (!File.Exists(path))
		{
			Debug.LogError("File not found: " + path);
			return;
		}

		string[] lines = File.ReadAllLines(path);

		if (lines.Length <= 1)
		{
			Debug.LogError("File is empty or has no data rows.");
			return;
		}

		for (int i = 1; i < lines.Length; i++)
		{
			if (string.IsNullOrWhiteSpace(lines[i]))
				continue;

			string[] parts = lines[i].Split(',');

			if (parts.Length < 10)
			{
				Debug.LogWarning("Skipped invalid line: " + lines[i]);
				continue;
			}

			try
			{
				int simId = int.Parse(parts[0]);
				if (simId != selectedSimId)
					continue;

				int step = int.Parse(parts[1]);
				int body = int.Parse(parts[2]);

				float xTrue = ParseFloat(parts[3]);
				float yTrue = ParseFloat(parts[4]);
				float zTrue = ParseFloat(parts[5]);

				float xPred = ParseFloat(parts[6]);
				float yPred = ParseFloat(parts[7]);
				float zPred = ParseFloat(parts[8]);

				MotionRecord record = new MotionRecord
				{
					simId = simId,
					step = step,
					body = body,
					truePosition = new Vector3(xTrue, yTrue, zTrue) * visualScale,
					predPosition = new Vector3(xPred, yPred, zPred) * visualScale
				};

				records.Add(record);
			}
			catch (Exception e)
			{
				Debug.LogWarning($"Error parsing line {i}: {e.Message}");
			}
		}

		steps.AddRange(records
			.Select(r => r.step)
			.Distinct()
			.OrderBy(s => s));

		Debug.Log($"Loaded sim_id={selectedSimId}, records={records.Count}, steps={steps.Count}");
	}

	float ParseFloat(string value)
	{
		return float.Parse(value, CultureInfo.InvariantCulture);
	}

	string BuildKey(int body)
	{
		return $"Body_{body}";
	}

	void CreateObjectsForFirstStep()
	{
		if (steps.Count == 0)
			return;

		int firstStep = steps[0];
		var firstStepRecords = records.Where(r => r.step == firstStep);

		foreach (var record in firstStepRecords)
		{
			string key = BuildKey(record.body);

			if (!trueObjects.ContainsKey(key))
			{
				GameObject trueSphere = GameObject.CreatePrimitive(PrimitiveType.Sphere);
				trueSphere.name = $"{key}_True";
				trueSphere.transform.position = record.truePosition;
				trueSphere.transform.localScale = Vector3.one * trueSphereSize;

				Renderer trueRenderer = trueSphere.GetComponent<Renderer>();
				trueRenderer.material = CreateMaterial(Color.green);

				trueObjects[key] = trueSphere;
			}

			if (!predObjects.ContainsKey(key))
			{
				GameObject predSphere = GameObject.CreatePrimitive(PrimitiveType.Sphere);
				predSphere.name = $"{key}_Pred";
				predSphere.transform.position = record.predPosition;
				predSphere.transform.localScale = Vector3.one * predSphereSize;

				Renderer predRenderer = predSphere.GetComponent<Renderer>();
				predRenderer.material = CreateMaterial(Color.red);

				predObjects[key] = predSphere;
			}
		}
	}

	void ApplyStep(int step)
	{
		var stepRecords = records.Where(r => r.step == step);

		foreach (var record in stepRecords)
		{
			string key = BuildKey(record.body);

			if (trueObjects.ContainsKey(key))
				trueObjects[key].transform.position = record.truePosition;

			if (predObjects.ContainsKey(key))
				predObjects[key].transform.position = record.predPosition;
		}
	}

	Material CreateMaterial(Color color)
	{
		Shader shader =
			Shader.Find("Universal Render Pipeline/Lit") ??
			Shader.Find("Universal Render Pipeline/Simple Lit") ??
			Shader.Find("Standard");

		Material material = new Material(shader);
		material.color = color;
		return material;
	}

	public void PausePlayback()
	{
		isPlaying = false;
	}

	public void ResumePlayback()
	{
		isPlaying = true;
	}

	public void RestartPlayback()
	{
		if (steps.Count == 0)
			return;

		currentStepIndex = 0;
		timer = 0f;
		ApplyStep(steps[0]);
	}
}