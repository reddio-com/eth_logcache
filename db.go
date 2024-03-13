package main

import (
	"encoding/binary"
	"fmt"

	"github.com/cockroachdb/pebble"
)

var (
	DB *pebble.DB
)

func Uint64ArrayToByteSlice(array []uint64) []byte {
	byteSlice := make([]byte, 8*len(array))
	for i, v := range array {
		binary.LittleEndian.PutUint64(byteSlice[i*8:], v)
	}
	return byteSlice
}

func ByteSliceToUint64Array(byteSlice []byte) []uint64 {
	array := make([]uint64, len(byteSlice)/8)
	for i := range array {
		array[i] = binary.LittleEndian.Uint64(byteSlice[i*8:])
	}
	return array
}

func DeduplicateUint64Array(array []uint64) []uint64 {
	keys := make(map[uint64]bool)
	list := []uint64{}
	for _, entry := range array {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func SetKeyNumberArray(key string, number uint64) error {
	// Check if key exists, if exists, read array and append number to array
	// If not exists, create array with number
	keyByteSlice := []byte(key)
	existingNumberArray, err := GetKeyNumberArray(key)
	if err != nil {
		return err
	}
	if existingNumberArray != nil {
		// Append number to array
		existingNumberArray = append(existingNumberArray, number)

		// Deduplicate array
		existingNumberArray = DeduplicateUint64Array(existingNumberArray)

		numberByteSlice := Uint64ArrayToByteSlice(existingNumberArray)
		err := DB.Set(keyByteSlice, numberByteSlice, pebble.NoSync)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Create array with number
		numberArray := []uint64{number}
		numberByteSlice := Uint64ArrayToByteSlice(numberArray)
		err := DB.Set(keyByteSlice, numberByteSlice, pebble.NoSync)
		if err != nil {
			return err
		}
		return nil
	}
}

func SetKeyNumber(key string, number uint64) error {
	keyByteSlice := []byte(key)
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(number))
	err := DB.Set(keyByteSlice, value, pebble.Sync)
	if err != nil {
		return err
	}
	return nil
}

func GetKeyNumberArray(key string) ([]uint64, error) {
	keyByteSlice := []byte(key)
	value, closer, err := DB.Get(keyByteSlice)
	if err != nil {
		// Check if key is not found
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()

	numberArray := ByteSliceToUint64Array(value)
	return numberArray, nil
}

func GetKeyNumber(key string) (uint64, error) {
	keyByteSlice := []byte(key)
	value, closer, err := DB.Get(keyByteSlice)
	if err != nil {
		// Check if key is not found
		if err == pebble.ErrNotFound {
			fmt.Printf("Key %s not found\n", key)
			return 0, nil
		}
	}
	defer closer.Close()

	number := binary.LittleEndian.Uint64(value)
	return number, nil
}

func FlushDB() error {
	err := DB.Flush()
	if err != nil {
		return err
	}
	return nil
}

func LatestBlockNumberInDB() (uint64, error) {
	key := "LatestBlockNumberInDB"
	number, err := GetKeyNumber(key)
	if err != nil {
		return 0, err
	}
	return number, nil
}
