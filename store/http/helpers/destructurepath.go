package helpers

import (
	"strings"
)

func DestructurePath(url string) (string, string, string) {
	var storage, collection, documentId string
	paths := strings.Split(url, "/")

	for pathCount := 0; pathCount < len(paths); pathCount++ {
		if strings.ToLower(paths[pathCount]) == "storages" {
			if len(paths) > pathCount+1 {
				storage = paths[pathCount+1]
				if len(paths) > pathCount+3 {
					collection = paths[pathCount+3]
					if len(paths) > pathCount+4 {
						documentId = paths[pathCount+4]
					}
				}
			}
		}
	}
	return storage, collection, documentId
}
