// This file is automatically generated by qtc from "error.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line template/error.qtpl:1
package template

//line template/error.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line template/error.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line template/error.qtpl:2
type InternalErrorPage struct {
	Backtrace interface{}
}

//line template/error.qtpl:7
func (ep *InternalErrorPage) StreamBody(qw422016 *qt422016.Writer) {
	//line template/error.qtpl:7
	qw422016.N().S(`
<pre>
    internal error occurred!

    `)
	//line template/error.qtpl:11
	qw422016.E().V(ep.Backtrace)
	//line template/error.qtpl:11
	qw422016.N().S(`
</pre>
`)
//line template/error.qtpl:13
}

//line template/error.qtpl:13
func (ep *InternalErrorPage) WriteBody(qq422016 qtio422016.Writer) {
	//line template/error.qtpl:13
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line template/error.qtpl:13
	ep.StreamBody(qw422016)
	//line template/error.qtpl:13
	qt422016.ReleaseWriter(qw422016)
//line template/error.qtpl:13
}

//line template/error.qtpl:13
func (ep *InternalErrorPage) Body() string {
	//line template/error.qtpl:13
	qb422016 := qt422016.AcquireByteBuffer()
	//line template/error.qtpl:13
	ep.WriteBody(qb422016)
	//line template/error.qtpl:13
	qs422016 := string(qb422016.B)
	//line template/error.qtpl:13
	qt422016.ReleaseByteBuffer(qb422016)
	//line template/error.qtpl:13
	return qs422016
//line template/error.qtpl:13
}

//line template/error.qtpl:16
type NotFoundPage struct{}

//line template/error.qtpl:19
func (ep *NotFoundPage) StreamBody(qw422016 *qt422016.Writer) {
	//line template/error.qtpl:19
	qw422016.N().S(`
<pre>
    404 error!
</pre>
`)
//line template/error.qtpl:23
}

//line template/error.qtpl:23
func (ep *NotFoundPage) WriteBody(qq422016 qtio422016.Writer) {
	//line template/error.qtpl:23
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line template/error.qtpl:23
	ep.StreamBody(qw422016)
	//line template/error.qtpl:23
	qt422016.ReleaseWriter(qw422016)
//line template/error.qtpl:23
}

//line template/error.qtpl:23
func (ep *NotFoundPage) Body() string {
	//line template/error.qtpl:23
	qb422016 := qt422016.AcquireByteBuffer()
	//line template/error.qtpl:23
	ep.WriteBody(qb422016)
	//line template/error.qtpl:23
	qs422016 := string(qb422016.B)
	//line template/error.qtpl:23
	qt422016.ReleaseByteBuffer(qb422016)
	//line template/error.qtpl:23
	return qs422016
//line template/error.qtpl:23
}

//line template/error.qtpl:26
type GeneralError struct {
	Message string
}

//line template/error.qtpl:31
func (ge *GeneralError) StreamBody(qw422016 *qt422016.Writer) {
	//line template/error.qtpl:31
	qw422016.N().S(`
<pre>
    `)
	//line template/error.qtpl:33
	qw422016.E().S(ge.Message)
	//line template/error.qtpl:33
	qw422016.N().S(`
</pre>
`)
//line template/error.qtpl:35
}

//line template/error.qtpl:35
func (ge *GeneralError) WriteBody(qq422016 qtio422016.Writer) {
	//line template/error.qtpl:35
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line template/error.qtpl:35
	ge.StreamBody(qw422016)
	//line template/error.qtpl:35
	qt422016.ReleaseWriter(qw422016)
//line template/error.qtpl:35
}

//line template/error.qtpl:35
func (ge *GeneralError) Body() string {
	//line template/error.qtpl:35
	qb422016 := qt422016.AcquireByteBuffer()
	//line template/error.qtpl:35
	ge.WriteBody(qb422016)
	//line template/error.qtpl:35
	qs422016 := string(qb422016.B)
	//line template/error.qtpl:35
	qt422016.ReleaseByteBuffer(qb422016)
	//line template/error.qtpl:35
	return qs422016
//line template/error.qtpl:35
}
