package types

import (
	"backend/internal/common/constants"
	"backend/share/base"
	"strings"
	"sync"
)

type UnsDefinition struct {
	CreateTopicDto
	Lock      sync.RWMutex `json:"-"`
	fieldLock sync.RWMutex
	qosField  string
}

// GetTimestampField returns the timestamp field name
func (c *UnsDefinition) GetTimestampField() string {
	if c.TmField == "" && len(c.Fields) > 0 {
		c.fieldLock.Lock()
		// Find timestamp field (implementation depends on FieldUtils)
		if c.TmField == "" {
			for _, f := range c.Fields {
				if f.Name == constants.SysFieldCreateTime || f.Name == "timestamp" {
					c.TmField = f.Name
					break
				}
			}
		}
		c.fieldLock.Unlock()
	}

	return c.TmField
}
func (c *UnsDefinition) GetFieldDefines() *FieldDefines {
	if c.FieldDefines == nil && len(c.Fields) > 0 {
		c.fieldLock.Lock()
		if c.FieldDefines == nil {
			c.SetFields(c.Fields)
		}
		c.fieldLock.Unlock()
	}
	return c.FieldDefines
}
func (c *UnsDefinition) GetPrimaryField() []string {
	if c.FieldDefines == nil {
		c.fieldLock.Lock()
		if c.FieldDefines == nil {
			c.SetFields(c.Fields)
		}
		c.fieldLock.Unlock()
	}
	return c.PrimaryField
}

// GetQualityField returns the quality field name
func (c *UnsDefinition) GetQualityField() string {
	if len(c.Fields) > 2 && c.qosField == "" {
		// Find quality field (implementation depends on FieldUtils and dataSrcId.typeCode)
		c.fieldLock.Lock()
		if c.qosField == "" {
			found := false
			for _, f := range c.Fields {
				if f.Type == FieldTypeLong && f.IsSystemField() && !base.P2v(f.Unique) {
					c.qosField = f.Name
					found = true
					break
				}
			}
			if !found {
				c.qosField = " "
			}
		}
		c.fieldLock.Unlock()
	}
	qos := c.qosField
	if len(qos) == 1 {
		qos = strings.TrimSpace(qos)
	}
	return qos
}
