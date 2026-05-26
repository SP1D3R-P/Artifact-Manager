package a2f_manager

import "testing"

func TestLoadValidA2F(t *testing.T) {
	a2fJson := `
	{
		"metadata": {
			"name": "helloC",
			"version": {
			"Major": 1,
			"Minor": 1,
			"Patch": 1
			},
			"build_time": "2026-05-15T20:42:52.4975433+05:30",
			"source_path": "C:\\Users\\LENOVO\\Desktop\\POC3\\FtestData\\helloC",
			"checksum": "f50632a88d2dae7fb534eec52ab85ae273e430a70d2aFaa4217ac3b7c47fe41f4",
			"location": "test"
		},
		"artifact-details": {
			"name": "main.exe",
			"size": 131283
		},
		"config": {
			"build": {
			"steps": [
				{
				"cmd": "echo $OPTIMIZATION_LEVEL"
				},
				{
				"cmd": "cc main.c -O$OPTIMIZATION_LEVEL -o main.exe"
				}
			],
			"environ": {
				"OPTIMIZATION_LEVEL": "3"
			}
			},
			"exec": {
			"steps": [
				{
				"cmd": "./main.exe",
				"input": [
					"priyankar"
				]
				}
			],
			"environ": {
				"A": "D"
			}
			}
		},
		"build-details": {
			"build-time": 0,
			"ouput": "test\\BuildResult.txt"
		}
	}
	`
	loader := NewLoader()
	_, err := loader.parse([]byte(a2fJson))
	if err != nil {
		t.Errorf("The Json Should Correctly Formated but not due to %s", err)
	}
}
