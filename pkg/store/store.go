package store

import (
	"os"
	"syscall"
	"unsafe"
)

type Record struct {
	Timestamp int64
	Point     float64
}

type DataStore struct {
	mmap         []byte
	maxEntrySize int64
	nextEntry    int64
	sizePerEntry int64
	mapFull      bool
}

func NewDataStore(entry int64, filePath string) (*DataStore, error) {
	mapFile, err := os.OpenFile(filePath, syscall.O_CREAT|syscall.O_RDWR, 0644)
	defer mapFile.Close()

	if err != nil {
		return nil, err
	}

	fdInfo, err := os.Stat(mapFile.Name())
	if err != nil {
		return nil, err
	}

	sizePerEntry := int64(unsafe.Sizeof(int64(0)) + unsafe.Sizeof(float64(0)))
	reqSize := entry * sizePerEntry

	// cleanup all entry if entry size had changed.
	if fdInfo.Size() != reqSize {
		zeroes := make([]byte, fdInfo.Size())
		mapFile.Write(zeroes)

		err := syscall.Ftruncate(int(mapFile.Fd()), int64(reqSize))
		if err != nil {
			return nil, err
		}
	}

	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, int(reqSize), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	var nextEntry int64
	nextEntry = 0
	for i := int64(0); i < reqSize; i += sizePerEntry {
		tmp := (*int64)(unsafe.Pointer(&mmap[i]))
		if *tmp == int64(0) {
			break
		}
		nextEntry++
	}
	if nextEntry >= entry {
		nextEntry = 0
	}

	return &DataStore{
		mmap:         mmap,
		maxEntrySize: entry,
		nextEntry:    nextEntry,
		sizePerEntry: sizePerEntry,
	}, nil
}

func (d *DataStore) WriteRecord(timestamp int64, point float64) {
	timestampEntryPtr := (*int64)(unsafe.Pointer(&d.mmap[d.sizePerEntry*d.nextEntry]))
	*timestampEntryPtr = timestamp
	pointEntryPtr := (*float64)(unsafe.Pointer(&d.mmap[d.sizePerEntry*d.nextEntry+int64(unsafe.Sizeof(int64(0)))]))
	*pointEntryPtr = point
	d.nextEntry += 1

	if d.nextEntry >= d.maxEntrySize {
		d.mapFull = true
		d.nextEntry = 0
	}
}

func (d *DataStore) ReadAll() []Record {
	var records []Record

	var readEntryLimit int64
	if d.mapFull {
		readEntryLimit = d.sizePerEntry * d.maxEntrySize
	} else {
		readEntryLimit = d.sizePerEntry * d.nextEntry
	}

	for i := int64(0); i < readEntryLimit; i += d.sizePerEntry {
		timestampEntryPtr := (*int64)(unsafe.Pointer(&d.mmap[i]))
		pointEntryPtr := (*float64)(unsafe.Pointer(&d.mmap[i+int64(unsafe.Sizeof(int64(0)))]))

		records = append(records, Record{
			Timestamp: *timestampEntryPtr,
			Point:     *pointEntryPtr,
		})
	}

	return records
}

func (d *DataStore) Close() error {
	return syscall.Munmap(d.mmap)
}
