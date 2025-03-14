package customplanmodifier

type UnknownAttributeReplacment[ResourceInfo any, TPFType any] struct {
	Name string
	Call func(TPFType, AttributeChanges, ResourceInfo, PlanModifyDiffer) *TPFType
}
