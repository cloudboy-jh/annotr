package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type CodeBlock struct {
	Type      string
	Name      string
	StartLine uint32
	EndLine   uint32
	StartByte uint32
	EndByte   uint32
	Code      string
	Context   string
}

type Parser struct {
	parser   *sitter.Parser
	language string
}

func NewParser(filename string) (*Parser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	lang, langName := getLanguage(ext)
	if lang == nil {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	return &Parser{
		parser:   parser,
		language: langName,
	}, nil
}

func (p *Parser) Language() string {
	return p.language
}

func (p *Parser) Parse(source []byte) ([]CodeBlock, error) {
	tree, err := p.parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	var blocks []CodeBlock
	p.extractBlocks(tree.RootNode(), source, &blocks)

	return blocks, nil
}

func (p *Parser) extractBlocks(node *sitter.Node, source []byte, blocks *[]CodeBlock) {
	nodeType := node.Type()

	if p.isCommentableNode(nodeType) {
		name := p.extractName(node, source)
		block := CodeBlock{
			Type:      nodeType,
			Name:      name,
			StartLine: node.StartPoint().Row,
			EndLine:   node.EndPoint().Row,
			StartByte: node.StartByte(),
			EndByte:   node.EndByte(),
			Code:      string(source[node.StartByte():node.EndByte()]),
		}
		*blocks = append(*blocks, block)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		p.extractBlocks(child, source, blocks)
	}
}

func (p *Parser) isCommentableNode(nodeType string) bool {
	commentableTypes := map[string]bool{
		"function_declaration":      true,
		"method_declaration":        true,
		"function_definition":       true,
		"method_definition":         true,
		"class_definition":          true,
		"class_declaration":         true,
		"type_declaration":          true,
		"interface_declaration":     true,
		"struct_type":               true,
		"function":                  true,
		"arrow_function":            true,
		"function_expression":       true,
		"lexical_declaration":       true,
	}
	return commentableTypes[nodeType]
}

func (p *Parser) extractName(node *sitter.Node, source []byte) string {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()
		if childType == "identifier" || childType == "name" || childType == "type_identifier" {
			return string(source[child.StartByte():child.EndByte()])
		}
		if childType == "function_declarator" || childType == "declarator" {
			return p.extractName(child, source)
		}
	}
	return ""
}

func getLanguage(ext string) (*sitter.Language, string) {
	switch ext {
	case ".go":
		return golang.GetLanguage(), "go"
	case ".py":
		return python.GetLanguage(), "python"
	case ".js":
		return javascript.GetLanguage(), "javascript"
	case ".ts":
		return typescript.GetLanguage(), "typescript"
	case ".tsx":
		return typescript.GetLanguage(), "typescript"
	default:
		return nil, ""
	}
}

func GetSupportedExtensions() []string {
	return []string{".go", ".py", ".js", ".ts", ".tsx"}
}

func IsSupportedFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, supported := range GetSupportedExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}
