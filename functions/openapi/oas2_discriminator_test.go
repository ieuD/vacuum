package openapi

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOAS2Discriminator_GetSchema(t *testing.T) {
	def := OAS2Discriminator{}
	assert.Equal(t, "oas2_discriminator", def.GetSchema().Name)
}

func TestOAS2Discriminator_RunRule(t *testing.T) {
	def := OAS2Discriminator{}
	res := def.RunRule(nil, model.RuleFunctionContext{})
	assert.Len(t, res, 0)
}

func TestOAS2Discriminator_RunRule_Success(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      tags: 
        - little
        - song
definitions:
  Song:
    discriminator: love
    type: object
    required: 
      - love`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	rule := buildOpenApiTestRuleAction(path, "oas2_discriminator", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	ctx.Index = model.NewSpecIndex(&rootNode)

	def := OAS2Discriminator{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 0)
}

func TestOAS2Discriminator_RunRule_Fail(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      tags: 
        - little
        - song
definitions:
  Song:
    discriminator: love
    type: object
    required: 
      - cuddles`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	rule := buildOpenApiTestRuleAction(path, "oas2_discriminator", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	ctx.Index = model.NewSpecIndex(&rootNode)

	def := OAS2Discriminator{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 1)
}

func TestOAS2Discriminator_RunRule_Fail_NoRequired(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      tags: 
        - little
        - song
definitions:
  Song:
    discriminator: love
    type: object`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	rule := buildOpenApiTestRuleAction(path, "oas2_discriminator", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	ctx.Index = model.NewSpecIndex(&rootNode)

	def := OAS2Discriminator{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 1)
}

func TestOAS2Discriminator_RunRule_Fail_DiscriminatorMap(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      tags: 
        - little
        - song
definitions:
  Song:
    discriminator: 
      thing: love
    type: object`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	rule := buildOpenApiTestRuleAction(path, "oas2_discriminator", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	ctx.Index = model.NewSpecIndex(&rootNode)

	def := OAS2Discriminator{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 2)
}