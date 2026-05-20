package hcl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/require"
)

var tf *tfexec.Terraform
var tfMutex sync.Mutex

func getTF() *tfexec.Terraform {
	tfMutex.Lock()
	defer tfMutex.Unlock()
	if tf != nil {
		return tf
	}
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		panic("terraform not found in PATH: " + err.Error())
	}
	tempDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tf, err = tfexec.NewTerraform(tempDir, execPath)
	if err != nil {
		panic(err)
	}
	return tf
}

func GetAttrVal(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	ret := make(map[string]cty.Value)
	AddAttributes(t, body, ret)
	for _, block := range body.Blocks {
		ret[block.Type] = GetAttrVal(t, block.Body)
	}
	return cty.ObjectVal(ret)
}

func AddAttributes(t *testing.T, body *hclsyntax.Body, ret map[string]cty.Value) {
	t.Helper()
	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		require.False(t, diags.HasErrors(), "failed to parse attribute %s: %s", name, diags.Error())
		ret[name] = val
	}
}

func PrettyHCL(tb testing.TB, content string) string {
	tb.Helper()
	builder := strings.Builder{}
	err := getTF().Format(tb.Context(), io.NopCloser(strings.NewReader(content)), &builder)
	require.NoError(tb, err)
	formatted := builder.String()
	return formatted
}

func CanonicalHCL(t *testing.T, def string) string {
	t.Helper()
	return string(hclwrite.Format(GetDefParser(t, def).Bytes()))
}

func GetDefParser(t *testing.T, def string) *hclwrite.File {
	t.Helper()
	parser, diags := hclwrite.ParseConfig([]byte(def), "", hcl.InitialPos)
	require.False(t, diags.HasErrors(), "failed to parse def: %s", diags.Error())
	return parser
}

func GetBlockBody(t *testing.T, block *hclwrite.Block) *hclsyntax.Body {
	t.Helper()
	parser, diags := hclparse.NewParser().ParseHCL(block.Body().BuildTokens(nil).Bytes(), "")
	require.False(t, diags.HasErrors(), "failed to parse block: %s", diags.Error())

	body, ok := parser.Body.(*hclsyntax.Body)
	require.True(t, ok, "unexpected *hclsyntax.Body type: %T", parser.Body)
	return body
}

// StringSliceToHCL converts a Go string slice to an HCL literal.
// Returns "null" for nil, "[]" for empty, or a quoted list like `["a", "b"]`.
func StringSliceToHCL(slice []string) string {
	if slice == nil {
		return "null"
	}
	if len(slice) == 0 {
		return "[]"
	}
	return fmt.Sprintf("[%s]", `"`+strings.Join(slice, `", "`)+`"`)
}
