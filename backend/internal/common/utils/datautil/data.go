package datautil

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/idutil"
	"backend/internal/types"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetBlobTagFromFieldDefine is a helper for GetBlobTag using a FieldDefine object.
// Corresponds to DataUtils.getBlobTag(String fileName, FieldDefine blobField)
func GetBlobTagFromFieldDefine(fileName string, blobField *types.FieldDefine) string {
	return GetBlobTag(fileName, blobField.Name, types.FieldType(blobField.Type))
}

// GetBlobTag generates a unique tag for a BLOB.
// Corresponds to DataUtils.getBlobTag(String fileName, String fieldName, FieldType fieldType)
func GetBlobTag(fileName, fieldName string, fieldType types.FieldType) string {
	return fmt.Sprintf("%s---%s-%s-%d", fieldType.String(), fileName, fieldName, idutil.NextID())
}

// SaveBlobDataFromFieldDefine is a helper for SaveBlobData using a FieldDefine object.
// Corresponds to DataUtils.saveBlobData(String fileAlias, FieldDefine fieldDefine, Object blobData)
func SaveBlobDataFromFieldDefine(fileAlias string, fieldDefine *types.FieldDefine, blobData any) (string, error) {
	return SaveBlobData(fileAlias, fieldDefine.Name, types.FieldType(fieldDefine.Type), blobData)
}

// SaveBlobData saves BLOB data (either base64 string or []byte) to a file.
// The file is named with a unique blob tag, which is returned.
// Corresponds to DataUtils.saveBlobData(String fileAlias, String fieldName, FieldType fieldType, Object blobData)
func SaveBlobData(fileAlias, fieldName string, fieldType types.FieldType, blobData any) (string, error) {
	if blobData == nil {
		return "", nil
	}

	var bytes []byte
	var err error

	switch data := blobData.(type) {
	case string:
		// The decoder is lenient, so we'll rely on it to error out for invalid base64.
		bytes, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("blobData is not a valid base64 string: %w", err)
		}
	case []byte:
		bytes = data
	default:
		return "", fmt.Errorf("unsupported blobData type: %T", blobData)
	}

	if len(bytes) == 0 {
		return "", nil
	}

	blobTag := GetBlobTag(fileAlias, fieldName, fieldType)
	filePath := filepath.Join(constants.BlobPath, blobTag)

	// Ensure the directory exists
	if err := os.MkdirAll(constants.BlobPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create blob directory: %w", err)
	}

	logx.Infof("saving blob: %s, size: %d", blobTag, len(bytes))
	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		logx.Errorf("failed to save blob data: %s", err)
		return "", fmt.Errorf("failed to write blob file: %s", err)
	}

	return blobTag, nil
}

// GetBlobData reads blob data from a file identified by the blob tag and returns it as a base64 string.
// Corresponds to DataUtils.getBlobData
func GetBlobData(blobTag string) (string, error) {
	if blobTag == "" {
		return "", nil
	}

	filePath := filepath.Join(constants.BlobPath, blobTag)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("blob file not found: %s", filePath)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read blob file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

// HandleBlobMsg processes a JSON message, replacing blob tags with their base64 encoded data.
// Corresponds to DataUtils.handleBolb (typo in Java)
func HandleBlobMsg(msg string, def *types.CreateTopicDto) (string, error) {
	if msg == "" || def == nil || !def.HasBlobField {
		return msg, nil
	}

	blobFields := def.FilterAllBlobField()
	if len(blobFields) == 0 {
		return msg, nil
	}

	var oldMsg map[string]any
	if err := json.Unmarshal([]byte(msg), &oldMsg); err != nil {
		return msg, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	data, dataOk := oldMsg["data"].(map[string]any)
	payload, payloadOk := oldMsg["payload"].(string)

	if !dataOk {
		return msg, nil // No 'data' field to process
	}

	modified := false
	for _, blobField := range blobFields {
		valueObj, ok := data[blobField.Name]
		if !ok {
			continue
		}

		valueStr, isString := valueObj.(string)
		if !isString {
			continue
		}

		if strings.HasPrefix(valueStr, types.FieldTypeBlob+"---") {
			blobValue, err := GetBlobData(valueStr)
			if err != nil {
				// Log the error but continue, matching Java's behavior of not stopping.
				logx.Errorf("failed to get blob data for tag: %s, error: %s", valueStr, err)
				continue
			}
			data[blobField.Name] = blobValue
			if payloadOk {
				payload = strings.Replace(payload, valueStr, blobValue, 1)
			}
			modified = true
		} else if strings.HasPrefix(valueStr, types.FieldTypeLBlob+"---") {
			data[blobField.Name] = ""
			if payloadOk {
				payload = strings.Replace(payload, valueStr, "", 1)
			}
			modified = true
		}
	}

	if !modified {
		return msg, nil
	}

	oldMsg["data"] = data
	if payloadOk {
		oldMsg["payload"] = payload
	}

	newMsgBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return msg, fmt.Errorf("failed to marshal modified message: %w", err)
	}

	return string(newMsgBytes), nil
}

// GetDefaultValue returns the default value for a given field type.
// Corresponds to DataUtils.getDefaultValue
func GetDefaultValue(fieldType types.FieldType) any {
	switch fieldType {
	case types.FieldTypeString, types.FieldTypeBlob:
		return ""
	case types.FieldTypeLBlob:
		return "-"
	case types.FieldTypeInteger, types.FieldTypeLong, types.FieldTypeFloat, types.FieldTypeDouble:
		return 0
	case types.FieldTypeBoolean:
		return false
	case types.FieldTypeDatetime:
		return 0
	default:
		return nil
	}
}

// TransformEmptyValue creates a JSON-like map with default values for a given UNS definition.
// Corresponds to DataUtils.transEmptyValue
func TransformEmptyValue(uns *types.CreateTopicDto, strTypeQos bool) map[string]any {
	data := make(map[string]any)
	if uns == nil || uns.Fields == nil {
		return data
	}

	for _, field := range uns.Fields {
		if field.Name != constants.QosField &&
			field.Name != constants.SysFieldCreateTime &&
			field.Name != constants.SysSaveTime {
			data[field.Name] = GetDefaultValue(types.FieldType(field.Type))
		}
	}

	// Bad quality status
	badQuality := int64(-9223372036854775808) // This is the signed int64 representation of 0x8000000000000000
	if strTypeQos {
		data[constants.QosField] = fmt.Sprintf("%x", uint64(badQuality))
	} else {
		data[constants.QosField] = badQuality
	}
	data[constants.SysFieldCreateTime] = SystemTimeMillis()

	return data
}

// SystemTimeMillis returns the current system time in milliseconds.
// Corresponds to System.currentTimeMillis()
func SystemTimeMillis() int64 {
	return time.Now().UnixMilli()
}
