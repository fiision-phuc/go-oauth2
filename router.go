package oauth2

type Router interface {
	/** Should router allow method to pass through. */
	ShouldAllow(method string) bool
	/** Handle http request. */
	ServeRequest(context *Context) bool

	/** Create a group of related functions. */
	Group(urlGroup string, function func(r Router))

	/** Handle COPY method. */
	Copy(urlPath string, handler interface{})
	/** Handle DELETE method. */
	Delete(urlPath string, handler interface{})
	/** Handle GET method. */
	Get(urlPath string, handler interface{})
	/** Handle HEAD method. */
	Head(urlPath string, handler interface{})
	/** Handle LINK method. */
	Link(urlPath string, handler interface{})
	/** Handle OPTIONS method. */
	Options(urlPath string, handler interface{})
	/** Handle PATCH method. */
	Patch(urlPath string, handler interface{})
	/** Handle POST method. */
	Post(urlPath string, handler interface{})
	/** Handle PURGE method. */
	Purge(urlPath string, handler interface{})
	/** Handle PUT method. */
	Put(urlPath string, handler interface{})
	/** Handle UNLINK method. */
	Unlink(urlPath string, handler interface{})
}
