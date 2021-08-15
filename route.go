package router

type (
	route struct {
		method      string
		handlerFunc func(context *XmusContext)
		middleware  []*middleware
	}
)
