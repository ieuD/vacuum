package core

import (
	"fmt"
	"github.com/daveshanley/vaccum/model"
	"github.com/daveshanley/vaccum/utils"
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
)

const (
	flat   string = "flat"
	camel  string = "camel"
	pascal string = "pascal"
	kebab  string = "kebab"
	cobol  string = "cobol"
	snake  string = "snake"
	macro  string = "macro"
)

type Casing struct {
	flat                  string
	camel                 string
	pascal                string
	kebab                 string
	cobol                 string
	snake                 string
	macro                 string
	separatorPattern      string
	ignoreDigits          bool
	separatorChar         string
	separatorAllowLeading bool
	compiled              bool
}

var casingTypes = []string{flat, camel, pascal, kebab, cobol, snake, macro}

func (c Casing) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{
		Required: []string{"type"},
		Properties: []model.RuleFunctionProperty{
			{
				Name: "type",
				Description: fmt.Sprintf("'casing' requires a 'type' to be supplied, which can be one of:"+
					" '%s'", casingTypes),
			},
		},
		ErrorMessage: "'alphabetical' function has invalid options supplied. Example valid options are 'type' = 'camel'" +
			" or 'disallowDigits' = true",
	}
}

func (c Casing) RunRule(nodes []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if len(nodes) != 1 { // there can only be a single node passed in to this function.
		return nil
	}

	var casingType string

	// check supplied type
	props := utils.ConvertInterfaceIntoStringMap(context.Options)
	if props["type"] == "" {
		return nil
	} else {
		casingType = props["type"]
	}

	// pull out props
	if props["disallowDigits"] != "" {
		c.ignoreDigits, _ = strconv.ParseBool(props["disallowDigits"])
	}

	if props["separator.char"] != "" {
		c.separatorChar = props["separator.char"]
	}

	if props["separator.allowLeading"] != "" {
		c.separatorAllowLeading, _ = strconv.ParseBool(props["separator.allowLeading"])
	}

	// if a separator is defined, and can be used as a leading char, and the node value is that
	// char (rune, what ever), then we're done.
	if len(nodes[0].Value) == 1 &&
		c.separatorChar != "" &&
		c.separatorAllowLeading == true &&
		c.separatorChar == nodes[0].Value {
		return nil
	}

	var results []model.RuleFunctionResult
	var pattern string

	if c.compiled == false {
		c.compileExpressions()
	}

	switch casingType {
	case camel:
		pattern = c.camel
		break
	case pascal:
		pattern = c.pascal
		break
	case kebab:
		pattern = c.kebab
		break
	case cobol:
		pattern = c.cobol
		break
	case snake:
		pattern = c.snake
		break
	case macro:
		pattern = c.macro
		break
	case flat:
		pattern = c.flat
		break
	}

	if c.separatorChar == "" {
		rx := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
		if !rx.MatchString(nodes[0].Value) {
			results = append(results, model.RuleFunctionResult{
				Message: fmt.Sprintf("'%s' is not %s case!", nodes[0].Value, casingType),
			})
		}
	} else {

		c.separatorPattern = fmt.Sprintf("[%s]", regexp.QuoteMeta(c.separatorChar))
		leadingSepPattern := ""
		var leadingPattern = ""
		leadingSepPattern = c.separatorPattern
		if c.separatorAllowLeading {
			leadingPattern = fmt.Sprintf("^(?:%[1]s)?%[3]s(?:%[2]s%[3]s)*$", leadingSepPattern, c.separatorPattern, pattern)
		} else {
			leadingPattern = fmt.Sprintf("^(?:%[1]s)+(?:%[2]s%[1]s)*$", pattern, c.separatorPattern)
		}

		rx := regexp.MustCompile(leadingPattern)
		if !rx.MatchString(nodes[0].Value) {
			results = append(results, model.RuleFunctionResult{
				Message: fmt.Sprintf("'%s' is not %s case!", nodes[0].Value, casingType),
			})
		}

	}

	return results
}

func (c *Casing) compileExpressions() {

	digits := "0-9"
	if c.ignoreDigits {
		digits = ""
	}

	c.flat = fmt.Sprintf("[a-z][a-z%[1]s]*", digits)
	c.camel = fmt.Sprintf("[a-z][a-z%[1]s]*(?:[A-Z%[1]s](?:[a-z%[1]s]+|$))*", digits)
	c.pascal = fmt.Sprintf("[A-Z][a-z%[1]s]*(?:[A-Z%[1]s](?:[a-z%[1]s]+|$))*", digits)
	c.kebab = fmt.Sprintf("[a-z%[1]s-]+", digits)
	c.cobol = fmt.Sprintf("[A-Z%[1]s-]+", digits)
	c.snake = fmt.Sprintf("[a-z%[1]s_]+", digits)
	c.macro = fmt.Sprintf("[A-Z%[1]s_]+", digits)
	c.compiled = true
}