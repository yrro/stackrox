package common

import (
	"fmt"
	"strings"

	"github.com/stackrox/rox/central/compliance/checks/msgfmt"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	pkgSet "github.com/stackrox/rox/pkg/set"
)

// GetProcess returns the commandline object that matches the process name
func GetProcess(ret *compliance.ComplianceReturn, processName string) (*compliance.CommandLine, bool) {
	for _, c := range ret.CommandLines {
		if strings.Contains(c.Process, processName) {
			return c, true
		}
	}
	return nil, false
}

// GetArgForFlag returns the arg that matches the passed key
func GetArgForFlag(args []*compliance.CommandLine_Args, key string) *compliance.CommandLine_Args {
	for _, a := range args {
		if a.Key == key {
			return a
		}
	}
	return nil
}

// GetValuesForFlag returns the values based on the key passes
func GetValuesForFlag(args []*compliance.CommandLine_Args, key string) []string {
	var values []string
	for _, a := range args {
		if a.Key == key {
			values = append(values, a.GetValues()...)
		}
	}
	return values
}

// GetValuesForCommandFromFlagsAndConfig returns the values for specific key from the args and unmarshalled config
func GetValuesForCommandFromFlagsAndConfig(args []*compliance.CommandLine_Args, c map[string]interface{}, key string) []string {
	values := GetValuesForFlag(args, key)

	var value interface{}
	value, ok := c[key]
	if !ok {
		value, ok = c[key+"s"]
	}
	if !ok {
		return values
	}
	switch obj := value.(type) {
	case string:
		values = append(values, obj)
	case []string:
		values = append(values, obj...)
	default:
		panic(fmt.Sprintf("Unsupported type: %T", obj))
	}
	return values
}

// CommandEvaluationFunc is a generic function that checks command lines
type CommandEvaluationFunc func(framework.ComplianceContext, []string, string, string, string)
type helperEvaluationFunc func([]string, string, string, string) (message string, passes bool)

// Info returns info with values set for the flag. Info is used when there is no strict determination of if the check is met
func Info(ctx framework.ComplianceContext, values []string, key, _, defaultVal string) {
	if len(values) == 0 {
		framework.Notef(ctx, "%q is to the default value of %q", key, defaultVal)
		return
	}
	framework.Notef(ctx, "%q is set to %q", key, msgfmt.FormatStrings(values...))
}

// Set checks whether or not a value is set in the command line
func Set(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, set)
}

// Unset checks whether or not a value is not set in the command line
func Unset(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, unset)
}

// Matches checks whether or not a value matches the target value exactly
func Matches(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, matches)
}

// OnlyContains checks whether or not a value contains only the target values (where target values are delimited by ",")
func OnlyContains(ctx framework.ComplianceContext, values []string, key, targets, defaultVal string) {
	resultWrapper(ctx, values, key, targets, defaultVal, onlyContains)
}

// NotMatches checks where or not a value matches the target value exactly
func NotMatches(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, notMatches)
}

// Contains checks where or not a value contains the target value
func Contains(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, contains)
}

// NotContains checks where or not a value contains the target value
func NotContains(ctx framework.ComplianceContext, values []string, key, target, defaultVal string) {
	resultWrapper(ctx, values, key, target, defaultVal, notContains)
}

func resultWrapper(ctx framework.ComplianceContext, values []string, key, target, defaultVal string, f helperEvaluationFunc) {
	msg, pass := f(values, key, target, defaultVal)
	if pass {
		framework.Pass(ctx, msg)
	} else {
		framework.Fail(ctx, msg)
	}
}

func set(values []string, key, _, _ string) (string, bool) {
	if len(values) > 0 {
		return fmt.Sprintf("%q is set to %s", key, msgfmt.FormatStrings(values...)), true
	}
	return fmt.Sprintf("%q is not set", key), false
}

func unset(values []string, key, _, _ string) (string, bool) {
	if len(values) == 0 {
		return fmt.Sprintf("%q is not set", key), true
	}
	return fmt.Sprintf("%q is set to %s", key, msgfmt.FormatStrings(values...)), false
}

func matches(values []string, key, target, defaultStr string) (string, bool) {
	var matchingValues []string
	var nonMatchingValues []string
	for _, v := range values {
		if strings.EqualFold(v, target) {
			matchingValues = append(matchingValues, v)
		} else {
			nonMatchingValues = append(nonMatchingValues, v)
		}
	}
	if len(matchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s", key, msgfmt.FormatStrings(matchingValues...)), true
	} else if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q is set to %q and not the target value of %q", key, msgfmt.FormatStrings(nonMatchingValues...), target), false
	} else if target == defaultStr {
		return fmt.Sprintf("%q has a default value that matches the target value of %q", key, defaultStr), true
	} else {
		return fmt.Sprintf("%q has a default value of %q that does not match the target value of %q", key, defaultStr, target), false
	}
}

func notMatches(values []string, key, target, defaultStr string) (string, bool) {
	var matchingValues []string
	var nonMatchingValues []string
	for _, v := range values {
		if strings.EqualFold(v, target) {
			matchingValues = append(matchingValues, v)
		} else {
			nonMatchingValues = append(nonMatchingValues, v)
		}
	}
	if len(matchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s which matches %q", key, msgfmt.FormatStrings(matchingValues...), target), false
	} else if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s which does not match %q", key, msgfmt.FormatStrings(nonMatchingValues...), target), true
	} else if target == defaultStr {
		return fmt.Sprintf("%q has a default value that matches the target value of %q", key, defaultStr), false
	} else {
		return fmt.Sprintf("%q has a default value of %q that does not match the target value of %q", key, defaultStr, target), true
	}
}

func onlyContains(values []string, key, targets, defaults string) (string, bool) {
	var matchingValues []string
	var nonMatchingValues []string

	targetSet := pkgSet.NewStringSet(strings.Split(targets, ",")...)
	for _, v := range values {
		if targetSet.Contains(v) {
			matchingValues = append(matchingValues, v)
		} else {
			nonMatchingValues = append(nonMatchingValues, v)
		}
	}

	if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s which contains values other than target values in %q", key, msgfmt.FormatStrings(nonMatchingValues...), targets), false
	} else if len(matchingValues) > 0 {
		numMatches := "some"
		if len(matchingValues) == targetSet.Cardinality() {
			numMatches = "all"
		}
		return fmt.Sprintf("%q is set to %s which contains %s target values in %q", key, msgfmt.FormatStrings(matchingValues...), numMatches, targets), true
	}

	defaultSet := pkgSet.NewStringSet(strings.Split(defaults, ";")...)
	for _, t := range targetSet.AsSlice() {
		if defaultSet.Contains(t) {
			matchingValues = append(matchingValues, t)
		} else {
			nonMatchingValues = append(nonMatchingValues, t)
		}
	}

	if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q has a default values %q which contains values other than target values in %q", key, msgfmt.FormatStrings(nonMatchingValues...), targets), false
	}
	numMatches := "some"
	if len(matchingValues) == targetSet.Cardinality() {
		numMatches = "all"
	}
	return fmt.Sprintf("%q has a default values %q which contains %s target values in %q", key, msgfmt.FormatStrings(matchingValues...), numMatches, targets), true
}

func contains(values []string, key, target, defaultStr string) (string, bool) {
	var matchingValues []string
	var nonMatchingValues []string
	for _, v := range values {
		if strings.Contains(v, target) {
			matchingValues = append(matchingValues, v)
		} else {
			nonMatchingValues = append(nonMatchingValues, v)
		}
	}
	if len(matchingValues) > 0 {
		return fmt.Sprintf("%q contains %s", key, msgfmt.FormatStrings(matchingValues...)), true
	} else if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q is set to %q and does not contain the target value of %q", key, msgfmt.FormatStrings(nonMatchingValues...), target), false
	} else if strings.Contains(defaultStr, target) {
		return fmt.Sprintf("%q has a default value that contains the target value of %q", key, defaultStr), true
	} else {
		return fmt.Sprintf("%q has a default value of %q that does not contain the target value of %q", key, defaultStr, target), false
	}
}

func notContains(values []string, key, target, defaultStr string) (string, bool) {
	var matchingValues []string
	var nonMatchingValues []string
	for _, v := range values {
		if strings.Contains(v, target) {
			matchingValues = append(matchingValues, v)
		} else {
			nonMatchingValues = append(nonMatchingValues, v)
		}
	}

	if len(matchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s which contains %q", key, msgfmt.FormatStrings(matchingValues...), target), false
	} else if len(nonMatchingValues) > 0 {
		return fmt.Sprintf("%q is set to %s which does not contain %q", key, msgfmt.FormatStrings(nonMatchingValues...), target), true
	} else if !strings.Contains(defaultStr, target) {
		return fmt.Sprintf("%q does not contain %q", key, target), true
	} else {
		return fmt.Sprintf("%q has a default value of %q that contains %q", key, defaultStr, target), false
	}
}
