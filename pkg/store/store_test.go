package store

import (
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	d, err := NewDataStore(3, "./test")
	defer d.Close()
	defer os.Remove("./test")

	if err != nil {
		t.Errorf("failed to create store: %v", err)
	}

	testdata := []struct {
		timestamp int64
		point     float64
	}{
		{
			timestamp: 1622914377,
			point:     32.4,
		},
		{
			timestamp: 1622914377,
			point:     32.5,
		},
		{
			timestamp: 1622914377,
			point:     32.6,
		},
	}

	for _, test := range testdata {
		d.WriteRecord(test.timestamp, test.point)
	}

	records := d.ReadAll()
	i := 0

	for _, record := range records {
		if record.Timestamp != testdata[i].timestamp {
			t.Errorf("not equal timestamp: %v, expected %v", record.Timestamp, testdata[i].timestamp)
		}
		if record.Point != testdata[i].point {
			t.Errorf("not equal data point: %v, expected %v", record.Point, testdata[i].point)
		}
		i++
	}
}

func TestStoreOver(t *testing.T) {
	d, err := NewDataStore(1, "./test")
	defer d.Close()
	defer os.Remove("./test")

	if err != nil {
		t.Errorf("failed to create store: %v", err)
	}

	testdata := []struct {
		timestamp int64
		point     float64
	}{
		{
			timestamp: 1622914377,
			point:     32.4,
		},
		{
			timestamp: 1622914377,
			point:     32.5,
		},
	}

	for _, test := range testdata {
		d.WriteRecord(test.timestamp, test.point)
	}

	records := d.ReadAll()
	if len(records) != 1 {
		t.Errorf("invalid length: %v, expected %v", len(records), 1)
	}
	record := records[0]

	if record.Timestamp != testdata[1].timestamp {
		t.Errorf("not equal timestamp: %v, expected %v", record.Timestamp, testdata[1].timestamp)
	}
	if record.Point != testdata[1].point {
		t.Errorf("not equal data point: %v, expected %v", record.Point, testdata[1].point)
	}
}

func TestRead(t *testing.T) {
	d, err := NewDataStore(2, "./test")

	if err != nil {
		t.Errorf("failed to create store: %v", err)
	}

	testdata := []struct {
		timestamp int64
		point     float64
	}{
		{
			timestamp: 1622914377,
			point:     32.4,
		},
		{
			timestamp: 1622914377,
			point:     32.5,
		},
	}

	d.WriteRecord(testdata[0].timestamp, testdata[0].point)

	d.Close()

	d2, err := NewDataStore(2, "./test")
	if err != nil {
		t.Errorf("failed to create store: %v", err)
	}

	d2.WriteRecord(testdata[1].timestamp, testdata[1].point)

	records := d.ReadAll()
	i := 0

	for _, record := range records {
		if record.Timestamp != testdata[i].timestamp {
			t.Errorf("not equal timestamp: %v, expected %v", record.Timestamp, testdata[i].timestamp)
		}
		if record.Point != testdata[i].point {
			t.Errorf("not equal data point: %v, expected %v", record.Point, testdata[i].point)
		}
		i++
	}

	d2.Close()
	os.Remove("./test")
}
