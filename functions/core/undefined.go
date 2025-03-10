// Copyright 2020-2021 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT

package core

import (
	"fmt"
	"github.com/daveshanley/vacuum/model"
	vacuumUtils "github.com/daveshanley/vacuum/utils"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// Undefined is a rule that will check if a field has not been defined.
type Undefined struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the Undefined rule.
func (u Undefined) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{
		Name: "undefined",
	}
}

// RunRule will execute the Undefined rule, based on supplied context and a supplied []*yaml.Node slice.
func (u Undefined) RunRule(nodes []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if len(nodes) <= 0 {
		return nil
	}

	var results []model.RuleFunctionResult

	message := context.Rule.Message

	pathValue := "unknown"
	if path, ok := context.Given.(string); ok {
		pathValue = path
	}

	ruleMessage := context.Rule.Description

	for _, node := range nodes {

		fieldNode, _ := utils.FindKeyNode(context.RuleAction.Field, node.Content)
		if fieldNode != nil {
			var val = ""
			if context.RuleAction.Field != "" {
				val = fmt.Sprintf("'%s' ", context.RuleAction.Field)
			}
			results = append(results, model.RuleFunctionResult{
				Message: vacuumUtils.SuppliedOrDefault(message, fmt.Sprintf("%s: `%s` must be undefined]",
					ruleMessage, val)),
				StartNode: fieldNode,
				EndNode:   fieldNode,
				Path:      pathValue,
				Rule:      context.Rule,
			})
		}
	}
	return results
}
