// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package owasp

import (
	"fmt"
	"github.com/daveshanley/vacuum/model"
	vacuumUtils "github.com/daveshanley/vacuum/utils"
	"github.com/pb33f/doctor/model/high/base"
	"gopkg.in/yaml.v3"
)

type RatelimitRetry429 struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the DefineError rule.
func (r RatelimitRetry429) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{Name: "ratelimit_retry_429"}
}

// RunRule will execute the DefineError rule, based on supplied context and a supplied []*yaml.Node slice.
func (r RatelimitRetry429) RunRule(_ []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	var results []model.RuleFunctionResult

	if context.DrDocument == nil {
		return results
	}

	for pathPairs := context.DrDocument.V3Document.Paths.PathItems.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		for opPairs := pathPairs.Value().GetOperations().First(); opPairs != nil; opPairs = opPairs.Next() {
			opValue := opPairs.Value()
			opType := opPairs.Key()

			responses := opValue.Responses.Codes

			for respPairs := responses.First(); respPairs != nil; respPairs = respPairs.Next() {
				resp := respPairs.Value()
				respCode := respPairs.Key()
				if respCode == "429" {

					var node *yaml.Node
					if resp.Headers != nil {
						foundHeader := resp.Headers.GetOrZero("Retry-After")
						if foundHeader == nil {
							lowCodes := opValue.Responses.Value.GoLow().Codes
							for lowCodePairs := lowCodes.First(); lowCodePairs != nil; lowCodePairs = lowCodePairs.Next() {
								lowCodeKey := lowCodePairs.Key()
								codeCodeVal := lowCodeKey.KeyNode.Value
								if codeCodeVal == "429" {
									node = lowCodeKey.KeyNode
								}
							}
							result := model.RuleFunctionResult{
								Message: vacuumUtils.SuppliedOrDefault(context.Rule.Message,
									"missing 'Retry-After' header for 429 error response"),
								StartNode: node,
								EndNode:   node,
								Path:      fmt.Sprintf("$.paths.%s.%s.responses.429", pathPairs.Key(), opType),
								Rule:      context.Rule,
							}
							resp.AddRuleFunctionResult(base.ConvertRuleResult(&result))
							results = append(results, result)
						}
					}
				}
			}
		}
	}
	return results
}
