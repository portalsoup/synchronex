package hashcode

import (
	"encoding/json"
	"hash/fnv"
	"log"
)

func HashCode(data interface{}) uint32 {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	hasher := fnv.New32a()
	hasher.Write(bytes)
	return hasher.Sum32()
}
