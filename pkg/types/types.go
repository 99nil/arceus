/*
Copyright © 2021 zc2638 <zc2638@qq.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package types

import "strings"

const (
	Group          = "arceus"
	CustomGroup    = "custom." + Group
	Version        = "v1"
	Kind           = "CustomResourceDefine"
	TemplateKind   = "Template"
	QuickStartKind = "QuickStart"
)

type JSONOperation struct {
	Op    string      `json:"op"` // add|remove|replace
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type TypeMeta struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ArceusResourceDefinition struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	// spec describes how the user wants the resources to appear
	Spec ArceusResourceDefinitionSpec `json:"spec"`
}

type ArceusResourceDefinitionSpec struct {
	// group is the API group of the defined custom resource.
	// The custom resources are served under `/apis/<group>/...`.
	// Must match the name of the CustomResourceDefinition (in the form `<names.plural>.<group>`).
	Group string `json:"group" protobuf:"bytes,1,opt,name=group"`
	// names specify the resource and kind names for the custom resource.
	Names ArceusResourceDefinitionNames `json:"names" protobuf:"bytes,3,opt,name=names"`
	// versions is the list of all API versions of the defined custom resource.
	// Version names are used to compute the order in which served versions are listed in API discovery.
	// If the version string is "kube-like", it will sort above non "kube-like" version strings, which are ordered
	// lexicographically. "Kube-like" versions start with a "v", then are followed by a number (the major version),
	// then optionally the string "alpha" or "beta" and another number (the minor version). These are sorted first
	// by GA > beta > alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing
	// major version, then minor version. An example sorted list of versions:
	// v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
	Versions []ArceusResourceDefinitionVersion `json:"versions"`
}

type ArceusResourceDefinitionNames struct {
	// plural is the plural name of the resource to serve.
	// The custom resources are served under `/apis/<group>/<version>/.../<plural>`.
	// Must match the name of the CustomResourceDefinition (in the form `<names.plural>.<group>`).
	// Must be all lowercase.
	Plural string `json:"plural,omitempty"`
	// singular is the singular name of the resource. It must be all lowercase. Defaults to lowercased `kind`.
	// +optional
	Singular string `json:"singular,omitempty"`
	// shortNames are short names for the resource, exposed in API discovery documents,
	// and used by clients to support invocations like `kubectl get <shortname>`.
	// It must be all lowercase.
	// +optional
	ShortNames []string `json:"shortNames,omitempty"`
	// kind is the serialized kind of the resource. It is normally CamelCase and singular.
	// Custom resource instances will use this value as the `kind` attribute in API calls.
	Kind string `json:"kind"`
	// listKind is the serialized kind of the list for this resource. Defaults to "`kind`List".
	// +optional
	ListKind string `json:"listKind,omitempty" protobuf:"bytes,5,opt,name=listKind"`
	// categories is a list of grouped resources this custom resource belongs to (e.g. 'all').
	// This is published in API discovery documents, and used by clients to support invocations like
	// `kubectl get all`.
	// +optional
	Categories []string `json:"categories,omitempty"`
}

type ArceusResourceDefinitionVersion struct {
	// name is the version name, e.g. “v1”, “v2beta1”, etc.
	// The custom resources are served under this version at `/apis/<group>/<version>/...` if `served` is true.
	Name string `json:"name"`
	// schema describes the schema used for validation, pruning, and defaulting of this version of the custom resource.
	Schema *ArceusResourceValidation `json:"schema,omitempty"`
}

type ArceusResourceValidation struct {
	// openAPIV3Schema is the OpenAPI v3 schema to use for validation and pruning.
	OpenAPIV3Schema *JSONSchemaProps `json:"openAPIV3Schema,omitempty"`
}

// Locale is a language code definition. (e.g. en、zh、de、ja).
type Locale string

// JSONSchemaProps is a JSON-Schema following Specification Draft 4 (http://json-schema.org/).
type JSONSchemaProps struct {
	ID string `json:"id,omitempty"`

	// Schema represents a schema url.
	Schema       string            `json:"$schema,omitempty"`
	Ref          *string           `json:"$ref,omitempty"`
	Description  string            `json:"description"`
	Descriptions map[Locale]string `json:"descriptions,omitempty"`
	Type         string            `json:"type,omitempty"`

	// format is an OpenAPI v3 format string. Unknown formats are ignored. The following formats are validated:
	//
	// - bsonobjectid: a bson object ID, i.e. a 24 characters hex string
	// - uri: an URI as parsed by Golang net/url.ParseRequestURI
	// - email: an email address as parsed by Golang net/mail.ParseAddress
	// - hostname: a valid representation for an Internet host name, as defined by RFC 1034, section 3.1 [RFC1034].
	// - ipv4: an IPv4 IP as parsed by Golang net.ParseIP
	// - ipv6: an IPv6 IP as parsed by Golang net.ParseIP
	// - cidr: a CIDR as parsed by Golang net.ParseCIDR
	// - mac: a MAC address as parsed by Golang net.ParseMAC
	// - uuid: an UUID that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid3: an UUID3 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?3[0-9a-f]{3}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid4: an UUID4 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?4[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - uuid5: an UUID5 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?5[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - isbn: an ISBN10 or ISBN13 number string like "0321751043" or "978-0321751041"
	// - isbn10: an ISBN10 number string like "0321751043"
	// - isbn13: an ISBN13 number string like "978-0321751041"
	// - creditcard: a credit card number defined by the regex ^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$ with any non digit characters mixed in
	// - ssn: a U.S. social security number following the regex ^\\d{3}[- ]?\\d{2}[- ]?\\d{4}$
	// - hexcolor: an hexadecimal color code like "#FFFFFF: following the regex ^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$
	// - rgbcolor: an RGB color code like rgb like "rgb(255,255,2559"
	// - byte: base64 encoded binary data
	// - password: any kind of string
	// - date: a date string like "2006-01-02" as defined by full-date in RFC3339
	// - duration: a duration string like "22 ns" as parsed by Golang time.ParseDuration or compatible with Scala duration format
	// - datetime: a date time string like "2014-12-15T19:30:20.000Z" as defined by date-time in RFC3339.
	Format string `json:"format,omitempty"`

	Title string `json:"title,omitempty"`
	// default is a default value for undefined object fields.
	Default              *string                    `json:"default,omitempty"`
	Maximum              *float64                   `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                       `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64                   `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                       `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int64                     `json:"maxLength,omitempty"`
	MinLength            *int64                     `json:"minLength,omitempty"`
	Pattern              string                     `json:"pattern,omitempty"`
	MaxItems             *int64                     `json:"maxItems,omitempty"`
	MinItems             *int64                     `json:"minItems,omitempty"`
	UniqueItems          bool                       `json:"uniqueItems,omitempty"`
	MultipleOf           *float64                   `json:"multipleOf,omitempty"`
	Enum                 []string                   `json:"enum,omitempty"`
	MaxProperties        *int64                     `json:"maxProperties,omitempty"`
	MinProperties        *int64                     `json:"minProperties,omitempty"`
	Required             []string                   `json:"required,omitempty"`
	Items                *JSONSchemaProps           `json:"items,omitempty"`
	AllOf                []JSONSchemaProps          `json:"allOf,omitempty"`
	OneOf                []JSONSchemaProps          `json:"oneOf,omitempty"`
	AnyOf                []JSONSchemaProps          `json:"anyOf,omitempty"`
	Not                  *JSONSchemaProps           `json:"not,omitempty"`
	Properties           map[string]JSONSchemaProps `json:"properties,omitempty"`
	AdditionalProperties *JSONSchemaPropsOrBool     `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]JSONSchemaProps `json:"patternProperties,omitempty"`
	Dependencies         JSONSchemaDependencies     `json:"dependencies,omitempty"`
	AdditionalItems      *JSONSchemaPropsOrBool     `json:"additionalItems,omitempty"`
	Definitions          JSONSchemaDefinitions      `json:"definitions,omitempty"`
	ExternalDocs         *ExternalDocumentation     `json:"externalDocs,omitempty"`
	Example              *string                    `json:"example,omitempty"`
	Nullable             bool                       `json:"nullable,omitempty"`
}

// JSONSchemaPropsOrBool represents JSONSchemaProps or a boolean value.
// Defaults to true for the boolean property.
type JSONSchemaPropsOrBool struct {
	Allows bool             `protobuf:"varint,1,opt,name=allows"`
	Schema *JSONSchemaProps `protobuf:"bytes,2,opt,name=schema"`
}

// JSONSchemaDependencies represent a dependencies property.
type JSONSchemaDependencies map[string]JSONSchemaPropsOrStringArray

// JSONSchemaPropsOrStringArray represents a JSONSchemaProps or a string array.
type JSONSchemaPropsOrStringArray struct {
	Schema   *JSONSchemaProps `protobuf:"bytes,1,opt,name=schema"`
	Property []string         `protobuf:"bytes,2,rep,name=property"`
}

// JSONSchemaDefinitions contains the models explicitly defined in this spec.
type JSONSchemaDefinitions map[string]JSONSchemaProps

// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	Description string `json:"description,omitempty" protobuf:"bytes,1,opt,name=description"`
	URL         string `json:"url,omitempty" protobuf:"bytes,2,opt,name=url"`
}

type QuickStartRule struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	// spec describes how the user wants the resources to appear
	Spec QuickStartRuleSpec `json:"spec"`
}

type QuickStartRuleSpec struct {
	Group     string               `json:"group"`
	Version   string               `json:"version"`
	Input     JSONSchemaProps      `json:"input"`
	Templates []RuleTemplateDefine `json:"templates"`
	Relate    []RuleRelate         `json:"relate"`
	Defines   []RuleDefine         `json:"defines"`
	Settings  []RuleSetting        `json:"settings"`
}

type RuleDefine struct {
	Path  string        `json:"path"`
	Value string        `json:"value"`
	Src   []interface{} `json:"src"`
}

type RuleTemplateDefine struct {
	Name     string                     `json:"name"`
	Template RuleTemplateResourceDefine `json:"template"`
}

type RuleTemplateResourceDefine struct {
	Name    string `json:"name"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

type RuleRelate struct {
	From RuleRelateFrom `json:"from"`
	To   RuleRelateTo   `json:"to"`
}

type RuleRelateFrom struct {
	TypeMeta
	Field string `json:"field"`
}

type RuleRelateTo struct {
	TypeMeta
	Fields []string `json:"fields"`
}

type RuleSetting struct {
	Path    string          `json:"path"`
	Targets []SettingTarget `json:"targets"`
}

type SettingTarget struct {
	Name   string               `json:"name"`   // template name
	Sub    string               `json:"sub"`    // template sub resource name
	Fields []SettingTargetField `json:"fields"` // field path
}

type SettingTargetField struct {
	Path string `json:"path"`
	Op   string `json:"op"`
}

type QuickStart struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	// spec describes how the user wants the resources to appear
	Spec QuickStartSpec `json:"spec"`
}

type QuickStartSpec struct {
	Rule []QuickStartSpecRule `json:"rule"`
	Data string               `json:"data"`
}

type QuickStartSpecRule struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type KValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KValuePairs []KValuePair

func (ps KValuePairs) Filter() KValuePairs {
	m := make(map[string]string)
	for _, p := range ps {
		m[p.Key] = p.Value
	}
	pairs := make(KValuePairs, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, KValuePair{
			Key:   k,
			Value: v,
		})
	}
	return pairs
}

func ParseKValuePairs(values []string) KValuePairs {
	if len(values) == 0 {
		return nil
	}
	pairs := make(KValuePairs, 0, len(values))
	for _, v := range values {
		vs := strings.SplitN(v, "=", 2)
		if len(vs) != 2 {
			continue
		}
		pairs = append(pairs, KValuePair{
			Key:   vs[0],
			Value: vs[1],
		})
	}
	return pairs
}
