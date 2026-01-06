package types

// FieldDefines represents a collection of field definitions
//type FieldDefines struct {
//	FieldsMap     map[string]*FieldDefine // Field name -> FieldDefine
//	FieldIndexMap map[string]string       // Index -> Field name
//	UniqueKeys    map[string]bool         // Set of unique field names
//	CalcField     *FieldDefine            // Calculation field
//}

// NewFieldDefines creates a new FieldDefines from a slice
func NewFieldDefines(fields []*FieldDefine) *FieldDefines {
	fd := &FieldDefines{
		FieldsMap:     make(map[string]*FieldDefine),
		FieldIndexMap: make(map[string]string),
		UniqueKeys:    make(map[string]bool),
	}

	if len(fields) > 0 {
		for _, f := range fields {
			fd.FieldsMap[f.Name] = f

			if f.IsUnique() {
				fd.UniqueKeys[f.Name] = true
			}

			if f.Index != nil {
				fd.FieldIndexMap[*f.Index] = f.Name
			}
		}
	}

	return fd
}

// NewFieldDefinesFromMap creates a new FieldDefines from a map
func NewFieldDefinesFromMap(fieldsMap map[string]*FieldDefine) *FieldDefines {
	if len(fieldsMap) == 0 {
		return &FieldDefines{
			FieldsMap:     make(map[string]*FieldDefine),
			FieldIndexMap: make(map[string]string),
			UniqueKeys:    make(map[string]bool),
		}
	}

	fields := make([]*FieldDefine, 0, len(fieldsMap))
	for _, f := range fieldsMap {
		fields = append(fields, f)
	}

	return NewFieldDefines(fields)
}

// ToFieldDefineArray converts to array
func (fd *FieldDefines) ToFieldDefineArray() []*FieldDefine {
	if len(fd.FieldsMap) == 0 {
		return []*FieldDefine{}
	}

	fields := make([]*FieldDefine, 0, len(fd.FieldsMap))
	for _, f := range fd.FieldsMap {
		fields = append(fields, f)
	}

	return fields
}
