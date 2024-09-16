package shared

import (
	"fmt"
)

func ConstructPBFilePath(collectionName, recordID, fileID string) string {
	return fmt.Sprintf("/api/files/%s/%s/%s", collectionName, recordID, fileID)
}
