package artifact

import (
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	data := []byte("test")
	checksum := calculateChecksum(data)
	if checksum == "" {
		t.Error("checksum should not be empty")
	}

	if len(checksum) == 0 {
		t.Error("checksum length should be > 0")
	}
}

func TestCalculateChecksumConsistency(t *testing.T) {
	data := []byte("same data")

	checksum1 := calculateChecksum(data)
	checksum2 := calculateChecksum(data)

	if checksum1 != checksum2 {
		t.Errorf("checksums should be consistent: '%s' != '%s'", checksum1, checksum2)
	}
}

func TestCalculateChecksumDifferentData(t *testing.T) {
	data1 := []byte("data1")
	data2 := []byte("data2")

	checksum1 := calculateChecksum(data1)
	checksum2 := calculateChecksum(data2)

	if checksum1 == checksum2 {
		t.Errorf("different data should produce different checksums: '%s' == '%s'", checksum1, checksum2)
	}
}

func TestCalculateChecksumEmpty(t *testing.T) {
	data := []byte{}
	checksum := calculateChecksum(data)

	if checksum == "" {
		t.Error("checksum should not be empty even for empty data")
	}
}
