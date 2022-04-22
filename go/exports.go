package polysdk

import (
	"polysdk/internal/polysign"
)

// polysign exports
const (
	// NOTE: if raw api has path args, pass them by this body object field
	XPolyBodyHideArgs = polysign.XPolyBodyHideArgs

	// NOTE: if origin body is not object, define customer body by this field
	XPolyCustomerBodyRoot = polysign.XPolyCustomerBodyRoot
)
