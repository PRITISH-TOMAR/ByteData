package tests

import (
	"fmt"
	"os"
	"testing"

	"byted/internal/wal"
)

func TestWAL(t *testing.T) {
	// temp WAL file for testing
	testPath := "test_wal.log"
	defer os.Remove(testPath) // clean up after test

	// create WAL
	w, err := wal.New(testPath)
	if err != nil {
		t.Fatal("Failed to create WAL:", err)
	}
	defer w.Close()

	// append some records
	lsn1, err := w.AppendPut([]byte("k01"), []byte("val-01"))
	if err != nil {
		t.Fatal("AppendPut failed:", err)
	}
	t.Logf("Appended PUT record LSN=%d", lsn1)

	lsn2, err := w.AppendDelete([]byte("k02"))
	if err != nil {
		t.Fatal("AppendDelete failed:", err)
	}
	t.Logf("Appended DELETE record LSN=%d", lsn2)

	// reopen the WAL file in read mode and dump bytes (for debug)
	f, err := os.Open(testPath)
	if err != nil {
		t.Fatal("Failed to reopen WAL file:", err)
	}
	defer f.Close()

	stat, _ := f.Stat()
	buf := make([]byte, stat.Size())
	_, err = f.Read(buf)
	if err != nil {
		t.Fatal("Failed to read WAL file:", err)
	}

	// print raw WAL bytes
	fmt.Println("Raw WAL file bytes:")
	for i, b := range buf {
		fmt.Printf("%02X ", b)
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
