package openapi_functions

import (
	"github.com/daveshanley/vaccum/model"
	"github.com/daveshanley/vaccum/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOperd_GetSchema(t *testing.T) {
	def := TagDefined{}
	assert.Equal(t, "tag_defined", def.GetSchema().Name)
}

func TestTagDefined_RunRule(t *testing.T) {
	def := TagDefined{}
	res := def.RunRule(nil, model.RuleFunctionContext{})
	assert.Len(t, res, 0)
}

func TestTagDefined_RunRule_Success(t *testing.T) {

	yml := `tags:
  - name: "princess"
  - name: "prince"
  - name: "hope"
  - name: "naughty_dog"
paths:
  /melody:
    post:
      tags:
       - "princess"
       - "hope"
  /maddox:
    get:
      tags:
       - "prince"
       - "hope"
  /ember:
    get:
      tags:
       - "naughty_dog"`

	path := "$"

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "tag_defined", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)

	def := TagDefined{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
	//assert.Equal(t, "the 'get' operation at path '/ember' contains a duplicate operationId 'littleSong'", res[0].Message)
}

func TestTagDefined_RunRule_Fail(t *testing.T) {

	yml := `tags:
  - name: "princess"
  - name: "prince"
  - name: "hope"
  - name: "naughty_dog"
paths:
  /melody:
    post:
      tags:
       - "princess"
       - "hope"
  /maddox:
    get:
      tags:
       - "prince"
       - "hope"
  /ember:
    get:
      tags:
       - "such_a_naughty_dog"`

	path := "$"

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "tag_defined", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)

	def := TagDefined{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 1)
	assert.Equal(t, "the 'get' operation at path '/ember' contains a tag 'such_a_naughty_dog', "+
		"that is not defined in the global document tags", res[0].Message)
}