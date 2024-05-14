package indexables

import (
	"fmt"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Struct struct {
	members    []StructMember
	isUnion    bool
	implements []string
	BaseIndexable
}

func NewStruct(name string, interfaces []string, members []StructMember, module string, docId string, idRange Range, docRange Range) Struct {
	return Struct{
		members:    members,
		isUnion:    false,
		implements: interfaces,
		BaseIndexable: NewBaseIndexable(
			name,
			module,
			docId,
			idRange,
			docRange,
			protocol.CompletionItemKindStruct,
		),
	}
}

func NewUnion(name string, members []StructMember, module string, docId string, idRange Range, docRange Range) Struct {
	return Struct{
		members: members,
		isUnion: true,
		BaseIndexable: NewBaseIndexable(
			name,
			module,
			docId,
			idRange,
			docRange,
			protocol.CompletionItemKindStruct,
		),
	}
}

func (s Struct) GetMembers() []StructMember {
	return s.members
}

func (s Struct) GetInterfaces() []string {
	return s.implements
}

func (s Struct) IsUnion() bool {
	return s.isUnion
}

func (s Struct) GetHoverInfo() string {
	return fmt.Sprintf("%s", s.name)
}

type StructMember struct {
	name     string
	baseType string
	BaseIndexable
}

func (m StructMember) GetName() string {
	return m.name
}

func (m StructMember) GetType() string {
	return m.baseType
}

func (m StructMember) GetIdRange() Range {
	return m.idRange
}

func (m StructMember) GetDocumentRange() Range {
	return m.docRange
}

func (m StructMember) GetDocumentURI() string {
	return m.documentURI
}

func (s StructMember) GetHoverInfo() string {
	return fmt.Sprintf("%s %s", s.baseType, s.name)
}
func (s StructMember) GetKind() protocol.CompletionItemKind {
	return s.Kind
}
func (s StructMember) GetModuleString() string {
	return s.moduleString
}

func (s StructMember) GetModule() ModulePath {
	return s.module
}

func (s StructMember) IsSubModuleOf(module ModulePath) bool {
	return s.module.IsSubModuleOf(module)
}

func NewStructMember(name string, baseType string, posRange Range, module string, docId string) StructMember {
	return StructMember{
		name:     name,
		baseType: baseType,
		BaseIndexable: BaseIndexable{
			idRange:      posRange,
			documentURI:  docId,
			moduleString: module,
			module:       NewModulePathFromString(module),
			Kind:         protocol.CompletionItemKindField,
		},
	}
}
