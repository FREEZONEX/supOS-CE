package types

func (r *UpdateTemplateFieldsAndDescReq) ResolveDescription() *string {
	if r == nil {
		return nil
	}
	if r.Description != nil {
		return r.Description
	}
	return r.ModelDescription
}
