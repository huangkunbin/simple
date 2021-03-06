package simpleapi

var appTemplate = `
// Code generated by simpleapi. DO NOT EDIT.
package {{Package}}

import (
	"simple/pkg/simplenet"
	"simple/pkg/simpleapi"
	"google.golang.org/protobuf/proto"
)

{{range .Imports}}
import {{.Name}} "{{.Path}}"
{{end}}

{{range .Services}}

func (s *{{.Name}}) ServiceID() byte {
	return {{.ID}}
}

func (s *{{.Name}}) NewRequest(id byte) (simpleapi.Message) {
	switch id {
	{{range .Requests}}
	case {{.ID}}:
		return &{{.Name}}{}
	{{end}}
	}
	return nil
}

func (s *{{.Name}}) NewResponse(id byte) (simpleapi.Message) {
	switch id {
	{{range .Responses}}
	case {{.ID}}:
		return &{{.Name}}{}
	{{end}}
	}
	return nil
}

func (s *{{.Name}}) HandleRequest(session simplenet.ISession, req simpleapi.Message) {
	switch req.MessageID() {
	{{range .Handlers}}
	case {{.ID}}:
		{{.InvokeCode}}
	{{end}}
	default:
		panic("Unhandled Message Type")
	}
}
{{end}}

{{range .Messages}}
func (m *{{.Name}}) ServiceID() byte {
	return {{.Service.ID}}
}

func (m *{{.Name}}) MessageID() byte {
	return {{.ID}}
}

func (m *{{.Name}}) Identity() string {
	return "{{.Service.Name}}.{{.Name}}"
}

func (m *{{.Name}}) Marshal() ([]byte, error) {
	return proto.Marshal(m)
}

func (m *{{.Name}}) Unmarshal(p []byte) error {
	return proto.Unmarshal(p, m)
}

{{end}}
`
