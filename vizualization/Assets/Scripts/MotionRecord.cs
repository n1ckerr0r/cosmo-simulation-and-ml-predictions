using UnityEngine;

[System.Serializable]
public class MotionRecord
{
	public int simId;
	public int step;
	public int body;
	public Vector3 truePosition;
	public Vector3 predPosition;
}