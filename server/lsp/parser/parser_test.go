package parser

import (
	"testing"

	"github.com/pherrymason/c3-lsp/lsp/document"
	idx "github.com/pherrymason/c3-lsp/lsp/symbols"
	"github.com/stretchr/testify/assert"
	"github.com/tliron/commonlog"
)

func createParser() Parser {
	logger := &commonlog.MockLogger{}
	return NewParser(logger)
}

func TestParses_empty_document(t *testing.T) {
	doc := document.NewDocument("empty", "")
	parser := createParser()

	symbols := parser.ParseSymbols(&doc)

	assert.Equal(t, 0, len(symbols.modules.Keys()))
}

func TestParses_TypedEnums(t *testing.T) {
	docId := "doc"
	source := `enum Colors:int { RED, BLUE, GREEN };`
	doc := document.NewDocument(docId, source)
	parser := createParser()

	t.Run("finds Colors enum identifier", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		module := symbols.Get("doc")

		assert.NotNil(t, module.Enums["Colors"])
		assert.Equal(t, "Colors", module.Enums["Colors"].GetName())
		assert.Equal(t, "int", module.Enums["Colors"].GetType())
		assert.Same(t, module.Enums["Colors"], module.Children()[0])
	})

	t.Run("reads ranges for enum", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		scope := symbols.Get("doc")
		enum := scope.Enums["Colors"]

		assert.Equal(t, idx.NewRange(0, 0, 0, 36), enum.GetDocumentRange(), "Wrong document rage")
		assert.Equal(t, idx.NewRange(0, 5, 0, 11), enum.GetIdRange(), "Wrong identifier range")
	})

	t.Run("finds defined enumerators", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		enum := symbols.Get("doc").Enums["Colors"]
		e := enum.GetEnumerator("RED")
		assert.Equal(t, "RED", e.GetName())
		assert.Equal(t, idx.NewRange(0, 18, 0, 21), e.GetIdRange())
		assert.Same(t, enum.Children()[0], e)

		e = enum.GetEnumerator("BLUE")
		assert.Equal(t, "BLUE", e.GetName())
		assert.Equal(t, idx.NewRange(0, 23, 0, 27), e.GetIdRange())
		assert.Same(t, enum.Children()[1], e)

		e = enum.GetEnumerator("GREEN")
		assert.Equal(t, "GREEN", e.GetName())
		assert.Equal(t, idx.NewRange(0, 29, 0, 34), e.GetIdRange())
		assert.Same(t, enum.Children()[2], e)
	})
}

func TestParses_UnTypedEnums(t *testing.T) {
	docId := "doc"
	source := `enum Colors { RED, BLUE, GREEN };`
	doc := document.NewDocument(docId, source)
	parser := createParser()

	t.Run("finds Colors enum identifier", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		scope := symbols.Get("doc")

		assert.NotNil(t, scope.Enums["Colors"])
		assert.Equal(t, "Colors", scope.Enums["Colors"].GetName())
		assert.Equal(t, "", scope.Enums["Colors"].GetType())
		assert.Same(t, scope.Children()[0], scope.Enums["Colors"])
	})

	t.Run("reads ranges for enum", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		enum := symbols.Get("doc").Enums["Colors"]
		assert.Equal(t, idx.NewRange(0, 0, 0, 32), enum.GetDocumentRange(), "Wrong document rage")
		assert.Equal(t, idx.NewRange(0, 5, 0, 11), enum.GetIdRange(), "Wrong identifier range")
	})

	t.Run("finds defined enumerators", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		enum := symbols.Get("doc").Enums["Colors"]
		e := enum.GetEnumerator("RED")
		assert.Equal(t, "RED", e.GetName())
		assert.Equal(t, idx.NewRange(0, 14, 0, 17), e.GetIdRange())

		e = enum.GetEnumerator("BLUE")
		assert.Equal(t, "BLUE", e.GetName())
		assert.Equal(t, idx.NewRange(0, 19, 0, 23), e.GetIdRange())

		e = enum.GetEnumerator("GREEN")
		assert.Equal(t, "GREEN", e.GetName())
		assert.Equal(t, idx.NewRange(0, 25, 0, 30), e.GetIdRange())
	})
}

func TestParse_fault(t *testing.T) {
	docId := "doc"
	source := `fault IOResult
	{
	  IO_ERROR,
	  PARSE_ERROR
	};`

	doc := document.NewDocument(docId, source)
	parser := createParser()

	t.Run("finds Fault identifier", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		scope := symbols.Get("doc")
		assert.NotNil(t, scope.Faults["IOResult"])
		assert.Equal(t, "IOResult", scope.Faults["IOResult"].GetName())
		assert.Equal(t, "", scope.Faults["IOResult"].GetType())
		assert.Same(t, scope.Children()[0], scope.Faults["IOResult"])
	})

	t.Run("reads ranges for fault", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		found := symbols.Get("doc").Faults["IOResult"]
		assert.Equal(t, idx.NewRange(0, 0, 4, 2), found.GetDocumentRange(), "Wrong document rage")
		assert.Equal(t, idx.NewRange(0, 6, 0, 14), found.GetIdRange(), "Wrong identifier range")
	})

	t.Run("finds defined fault constants", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		fault := symbols.Get("doc").Faults["IOResult"]
		e := fault.GetConstant("IO_ERROR")
		assert.Equal(t, "IO_ERROR", e.GetName())
		assert.Equal(t, idx.NewRange(2, 3, 2, 11), e.GetIdRange())
		assert.Same(t, fault.Children()[0], e)

		e = fault.GetConstant("PARSE_ERROR")
		assert.Equal(t, "PARSE_ERROR", e.GetName())
		assert.Equal(t, idx.NewRange(3, 3, 3, 14), e.GetIdRange())
		assert.Same(t, fault.Children()[1], e)
	})
}

func TestParse_interface(t *testing.T) {
	module := "x"
	docId := "doc"
	source := `interface MyName
	{
		fn String method();
	};`

	doc := document.NewDocument(docId, source)
	parser := createParser()

	t.Run("finds interface", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		expected := idx.NewInterfaceBuilder("MyName", module, docId).
			Build()

		module := symbols.Get("doc")
		interfac := module.Interfaces["MyName"]
		assert.NotNil(t, interfac)
		assert.Same(t, module.Children()[0], interfac)

		assert.Equal(t, expected.GetName(), interfac.GetName())
	})

	t.Run("reads ranges for interface", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		found := symbols.Get("doc").Interfaces["MyName"]
		assert.Equal(t, idx.NewRange(0, 0, 3, 2), found.GetDocumentRange(), "Wrong document rage")
		assert.Equal(t, idx.NewRange(0, 10, 0, 16), found.GetIdRange(), "Wrong identifier range")
	})

	t.Run("finds defined methods in interface", func(t *testing.T) {
		symbols := parser.ParseSymbols(&doc)

		module := symbols.Get("doc")
		_interface := module.Interfaces["MyName"]
		m := _interface.GetMethod("method")
		assert.Equal(t, "method", m.GetName())
		assert.Equal(t, "String", m.GetReturnType())
		assert.Equal(t, idx.NewRange(2, 12, 2, 18), m.GetIdRange())
		assert.Equal(t, module.Children()[0], _interface)
	})
}

func TestExtractSymbols_finds_definition(t *testing.T) {
	source := `module mod;
	def Kilo = int;
	def KiloPtr = Kilo*;
	def MyFunction = fn void (Allocator*, JSONRPCRequest*, JSONRPCResponse*);
	def MyMap = HashMap(<String, Feature>);
	`
	// TODO: Missing def different definition examples. See parser.nodeToDef
	mod := "mod"
	doc := document.NewDocument("x", source)
	parser := createParser()

	symbols := parser.ParseSymbols(&doc)

	expectedDefKilo := idx.NewDefBuilder("Kilo", mod, "x").
		WithResolvesTo("int").
		WithIdentifierRange(1, 5, 1, 9).
		WithDocumentRange(1, 1, 1, 16).
		Build()

	expectedDefKiloPtr := idx.NewDefBuilder("KiloPtr", mod, "x").
		WithResolvesTo("Kilo*").
		WithIdentifierRange(2, 5, 2, 12).
		WithDocumentRange(2, 1, 2, 21).
		Build()

	expectedDefFunction := idx.NewDefBuilder("MyFunction", mod, "x").
		WithResolvesTo("fn void (Allocator*, JSONRPCRequest*, JSONRPCResponse*)").
		WithIdentifierRange(3, 5, 3, 15).
		WithDocumentRange(3, 1, 3, 74).
		Build()

	expectedDefTypeWithGenerics := idx.NewDefBuilder("MyMap", mod, "x").
		WithResolvesTo("HashMap(<String, Feature>)").
		WithIdentifierRange(4, 5, 4, 10).
		WithDocumentRange(4, 1, 4, 40).
		Build()

	module := symbols.Get(mod)
	assert.Equal(t, expectedDefKilo, module.Defs["Kilo"])
	assert.Same(t, module.Children()[0], module.Defs["Kilo"])

	assert.Equal(t, expectedDefKiloPtr, module.Defs["KiloPtr"])
	assert.Same(t, module.Children()[1], module.Defs["KiloPtr"])

	assert.Equal(t, expectedDefFunction, module.Defs["MyFunction"])
	assert.Same(t, module.Children()[2], module.Defs["MyFunction"])

	assert.Equal(t, expectedDefTypeWithGenerics, module.Defs["MyMap"])
	assert.Same(t, module.Children()[3], module.Defs["MyMap"])
}

func TestExtractSymbols_find_macro(t *testing.T) {
	/*
		sourceCode := `
		macro void log(LogLevel $level, String format, args...) {
			if (log_level != OFF && $level <= log_level) {
				io::fprintf(&log_file, "[%s] ", $level)!!;
				io::fprintfn(&log_file, format, ...args)!!;
			}
		}`*/
	source := `
	macro m(x) {
    	return x + 2;
	}`

	doc := document.NewDocument("docId", source)
	parser := createParser()
	symbols := parser.ParseSymbols(&doc)

	module := symbols.Get("docid")
	fn := module.GetChildrenFunctionByName("m")
	assert.True(t, fn.IsSome())
	assert.Equal(t, "m", fn.Get().GetName())
	assert.Equal(t, "x", fn.Get().Variables["x"].GetName())
	assert.Equal(t, "", fn.Get().Variables["x"].GetType().String())
	assert.Same(t, module.NestedScopes()[0], fn.Get())
}

func TestExtractSymbols_find_module(t *testing.T) {
	t.Run("finds anonymous module", func(t *testing.T) {
		source := `int value = 1;`

		doc := document.NewDocument("file name.c3", source)
		parser := createParser()
		symbols := parser.ParseSymbols(&doc)
		fn := symbols.Get("file_name")

		assert.Equal(t, "file_name", fn.GetModuleString(), "Function module is wrong")
	})

	t.Run("finds single module in single file", func(t *testing.T) {
		source := `
	module foo;
	int value = 1;
	`

		doc := document.NewDocument("docId", source)
		parser := createParser()
		symbols := parser.ParseSymbols(&doc)

		module := symbols.Get("foo")
		assert.Equal(t, "foo", module.GetModuleString(), "module name is wrong")
	})

	t.Run("finds different modules defined in single file", func(t *testing.T) {
		source := `
	module foo;
	int value = 1;

	module foo2;
	int value = 2;`

		doc := document.NewDocument("docid", source)
		parser := createParser()
		symbols := parser.ParseSymbols(&doc)

		module := symbols.Get("foo")
		assert.Equal(t, "foo", module.GetModuleString(), "module name is wrong")
		assert.Equal(t, "foo", module.GetName(), "module name is wrong")
		assert.Equal(t, idx.NewRange(1, 1, 2, 15), module.GetDocumentRange(), "Wrong range for foo module")

		module = symbols.Get("foo2")
		assert.Equal(t, "foo2", module.GetModuleString(), "module name is wrong")
		assert.Equal(t, "foo2", module.GetName(), "module name is wrong")
		assert.Equal(t, idx.NewRange(4, 1, 5, 15), module.GetDocumentRange(), "Wrong range for foo2 module")
	})
}

func TestExtractSymbols_find_imports(t *testing.T) {
	source := `
	module foo;
	import some, other, foo::bar::final;
	import another;
	import another2;
	int value = 1;
	`

	doc := document.NewDocument("docid", source)
	parser := createParser()
	symbols := parser.ParseSymbols(&doc)

	assert.Equal(t, []string{"some", "other", "foo::bar::final", "another", "another2"}, symbols.Get("foo").Imports)
}
