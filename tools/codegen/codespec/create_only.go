package codespec

func setCreateOnlyInAttrs(attrs []Attribute) {
	for i := range attrs {
		setCreateOnlyValue(&attrs[i])

		switch {
		case attrs[i].ListNested != nil:
			setCreateOnlyInAttrs(attrs[i].ListNested.NestedObject.Attributes)
		case attrs[i].SingleNested != nil:
			setCreateOnlyInAttrs(attrs[i].SingleNested.NestedObject.Attributes)
		case attrs[i].SetNested != nil:
			setCreateOnlyInAttrs(attrs[i].SetNested.NestedObject.Attributes)
		case attrs[i].MapNested != nil:
			setCreateOnlyInAttrs(attrs[i].MapNested.NestedObject.Attributes)
		}
	}
}

func setCreateOnlyValue(attr *Attribute) {
	if attr.ComputedOptionalRequired == Computed {
		return
	}

	// captures case of path param attributes (no present in request body) and properties which are only present in post request
	if attr.ReqBodyUsage == OmitAlways || attr.ReqBodyUsage == OmitInUpdateBody {
		attr.CreateOnly = true
	}
}
