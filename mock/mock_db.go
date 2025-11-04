package mock

import (
	"encoding/json"
	"fmt"
	"os"
)

var incrementalIds = make(map[string]int)

const defaultDBPath = "database/"
const storedIds = defaultDBPath + "ids.json"

type MockDB[T any] struct {
	FilePath string
}

func NewDb[T any](filePath string) *MockDB[T] {
	db := &MockDB[T]{
		FilePath: defaultDBPath + filePath,
	}

	if err := db.createFileIfNotExist(); err != nil {
		fmt.Println("warning: failed to ensure file exists:", err)
	}
	return db
}

func (db *MockDB[T]) InsertData(data T) error {
	datas, err := db.GetData()
	if err != nil {
		return err
	}

	id := db.getIncrementId()

	dataMap := make(map[string]interface{})
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(rawData, &dataMap); err != nil {
		return err
	}
	dataMap["id"] = id

	var dataWithID T
	rawWithID, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(rawWithID, &dataWithID); err != nil {
		return err
	}

	datas = append(datas, dataWithID)

	raw, err := json.MarshalIndent(datas, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.FilePath, raw, 0644); err != nil {
		return err
	}

	return nil
}

func (db *MockDB[T]) GetData() ([]T, error) {
	var datas []T

	file, err := os.ReadFile(db.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return []T{}, nil
	}

	if err := json.Unmarshal(file, &datas); err != nil {
		return nil, err
	}

	return datas, nil
}

func (db *MockDB[T]) UpdateData(datas []T) error {
	raw, err := json.MarshalIndent(datas, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.FilePath, raw, 0644); err != nil {
		return err
	}

	return nil
}

func (db *MockDB[T]) createFileIfNotExist() error {
	if _, err := os.Stat(db.FilePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.WriteFile(db.FilePath, []byte("[]"), 0644); err != nil {
		return err
	}

	fmt.Printf("Created new %s with empty array\n", db.FilePath)
	return nil
}

func (db *MockDB[T]) getIncrementId() int {

	id, exists := incrementalIds[db.FilePath]
	if !exists {
		id = 0
	}
	incremented := id + 1
	incrementalIds[db.FilePath] = incremented
	raw, err := json.MarshalIndent(incrementalIds, "", "  ")
	if err != nil {
		return 0
	}

	if err := os.WriteFile(storedIds, raw, 0644); err != nil {
		return 0
	}

	return incremented
}

func InitDbIncremental() {
	// read file
	file, err := os.ReadFile(storedIds)
	if err != nil {
		if os.IsNotExist(err) {
			incrementalIds = make(map[string]int)

			raw, err := json.MarshalIndent(incrementalIds, "", "  ")
			if err != nil {
				fmt.Println("warning: failed to marshal ids.json:", err)
				return
			}

			if err := os.WriteFile(storedIds, raw, 0644); err != nil {
				fmt.Println("warning: failed to write ids.json:", err)
				return
			}
			return
		}
		fmt.Println("warning: failed to read ids.json:", err)
		return
	}

	if len(file) == 0 {
		incrementalIds = make(map[string]int)
		return
	}

	if err := json.Unmarshal(file, &incrementalIds); err != nil {
		fmt.Println("warning: failed to unmarshal ids.json:", err)
		return
	}

	// write back to file
	raw, err := json.MarshalIndent(incrementalIds, "", "  ")
	if err != nil {
		fmt.Println("warning: failed to marshal ids.json:", err)
		return
	}

	if err := os.WriteFile(storedIds, raw, 0644); err != nil {
		fmt.Println("warning: failed to write ids.json:", err)
		return
	}
}
