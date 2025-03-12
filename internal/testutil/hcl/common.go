package hcl

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"

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
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.10.1")),
	}
	execPath, err := installer.Install(context.Background())
	if err != nil {
		panic(err)
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

func PrettyHCL(t *testing.T, content string) string {
	t.Helper()
	builder := strings.Builder{}
	fmt := getTF().Format(t.Context(), io.NopCloser(strings.NewReader(content)), &builder)
	require.NoError(t, fmt)
	formatted := builder.String()
	return formatted
}

func CanonicalHCL(t *testing.T, def string) string {
	t.Helper()
	return string(hclwrite.Format(GetDefParser(t, def).Bytes()))
}

func GetDefParser(t *testing.T, def string) *hclwrite.File {
	t.Helper()
	parser, diags := hclwrite.ParseConfig([]byte(def), "", hcl.Pos{Line: 1, Column: 1})
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
