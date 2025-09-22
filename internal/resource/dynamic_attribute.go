package resource

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"

	"github.com/hashicorp/terraform-plugin-codegen-framework/internal/convert"
	"github.com/hashicorp/terraform-plugin-codegen-framework/internal/model"
	generatorschema "github.com/hashicorp/terraform-plugin-codegen-framework/internal/schema"
)

type GeneratorDynamicAttribute struct {
	AssociatedExternalType   *generatorschema.AssocExtType
	ComputedOptionalRequired convert.ComputedOptionalRequired
	CustomType               convert.CustomTypePrimitive
	Default                  convert.DefaultDynamic
	DeprecationMessage       convert.DeprecationMessage
	Description              convert.Description
	PlanModifiers            convert.PlanModifiers
	Sensitive                convert.Sensitive
	Validators               convert.Validators
}

func NewGeneratorDynamicAttribute(name string, a *resource.DynamicAttribute) (GeneratorDynamicAttribute, error) {
	if a == nil {
		return GeneratorDynamicAttribute{}, fmt.Errorf("*resource.DynamicAttribute is nil")
	}

	c := convert.NewComputedOptionalRequired(a.ComputedOptionalRequired)

	ctp := convert.NewCustomTypePrimitive(a.CustomType, a.AssociatedExternalType, name)

	ds := convert.NewDefaultDynamic(a.Default)

	dm := convert.NewDeprecationMessage(a.DeprecationMessage)

	d := convert.NewDescription(a.Description)

	pm := convert.NewPlanModifiers(convert.PlanModifierTypeString, a.PlanModifiers.CustomPlanModifiers())

	s := convert.NewSensitive(a.Sensitive)

	v := convert.NewValidators(convert.ValidatorTypeString, a.Validators.CustomValidators())

	return GeneratorDynamicAttribute{
		AssociatedExternalType:   generatorschema.NewAssocExtType(a.AssociatedExternalType),
		ComputedOptionalRequired: c,
		CustomType:               ctp,
		Default:                  ds,
		DeprecationMessage:       dm,
		Description:              d,
		PlanModifiers:            pm,
		Sensitive:                s,
		Validators:               v,
	}, nil
}

func (g GeneratorDynamicAttribute) GeneratorSchemaType() generatorschema.Type {
	return generatorschema.GeneratorDynamicAttribute
}

func (g GeneratorDynamicAttribute) Imports() *generatorschema.Imports {
	imports := generatorschema.NewImports()

	imports.Append(g.CustomType.Imports())
	imports.Append(g.Default.Imports())
	imports.Append(g.PlanModifiers.Imports())
	imports.Append(g.Validators.Imports())

	if g.AssociatedExternalType != nil {
		imports.Append(generatorschema.AssociatedExternalTypeImports())
	}

	imports.Append(g.AssociatedExternalType.Imports())

	return imports
}

func (g GeneratorDynamicAttribute) Equal(ga generatorschema.GeneratorAttribute) bool {
	h, ok := ga.(GeneratorDynamicAttribute)

	if !ok {
		return false
	}

	if !g.AssociatedExternalType.Equal(h.AssociatedExternalType) {
		return false
	}

	if !g.ComputedOptionalRequired.Equal(h.ComputedOptionalRequired) {
		return false
	}

	if !g.CustomType.Equal(h.CustomType) {
		return false
	}

	if !g.Default.Equal(h.Default) {
		return false
	}

	if !g.DeprecationMessage.Equal(h.DeprecationMessage) {
		return false
	}

	if !g.Description.Equal(h.Description) {
		return false
	}

	if !g.PlanModifiers.Equal(h.PlanModifiers) {
		return false
	}

	if !g.Sensitive.Equal(h.Sensitive) {
		return false
	}

	return g.Validators.Equal(h.Validators)
}

func (g GeneratorDynamicAttribute) Schema(name generatorschema.FrameworkIdentifier) (string, error) {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("%q: schema.DynamicAttribute{\n", name))
	b.Write(g.CustomType.Schema())
	b.Write(g.ComputedOptionalRequired.Schema())
	b.Write(g.Sensitive.Schema())
	b.Write(g.Description.Schema())
	b.Write(g.DeprecationMessage.Schema())
	b.Write(g.PlanModifiers.Schema())
	b.Write(g.Validators.Schema())
	b.Write(g.Default.Schema())
	b.WriteString("},")

	return b.String(), nil
}

func (g GeneratorDynamicAttribute) ModelField(name generatorschema.FrameworkIdentifier) (model.Field, error) {
	field := model.Field{
		Name:      name.ToPascalCase(),
		TfsdkName: name.ToString(),
		ValueType: model.DynamicValueType,
	}

	customValueType := g.CustomType.ValueType()

	if customValueType != "" {
		field.ValueType = customValueType
	}

	return field, nil
}

func (g GeneratorDynamicAttribute) CustomTypeAndValue(name string) ([]byte, error) {
	if g.AssociatedExternalType == nil {
		return nil, nil
	}

	var buf bytes.Buffer

	dynamicType := generatorschema.NewCustomDynamicType(name)

	b, err := dynamicType.Render()

	if err != nil {
		return nil, err
	}

	buf.Write(b)

	// TODO: this doesn't exist yet
	dynamicValue := generatorschema.NewCustomDynamicValue(name)

	b, err = dynamicValue.Render()

	if err != nil {
		return nil, err
	}

	buf.Write(b)

	return buf.Bytes(), nil
}

func (g GeneratorDynamicAttribute) ToFromFunctions(name string) ([]byte, error) {
	if g.AssociatedExternalType == nil {
		return nil, nil
	}

	toFrom := generatorschema.NewToFromDynamic(name, g.AssociatedExternalType)

	b, err := toFrom.Render()

	if err != nil {
		return nil, err
	}

	return b, nil
}

// AttrType returns a string representation of a basetypes.DynamicTypable type.
func (g GeneratorDynamicAttribute) AttrType(name generatorschema.FrameworkIdentifier) (string, error) {
	if g.AssociatedExternalType != nil {
		return fmt.Sprintf("%sType{}", name.ToPascalCase()), nil
	}

	return "basetypes.DynamicType{}", nil
}

func (g GeneratorDynamicAttribute) AttrValue(name generatorschema.FrameworkIdentifier) string {
	if g.AssociatedExternalType != nil {
		return fmt.Sprintf("%sValue", name.ToPascalCase())
	}

	return "basetypes.DynamicValue"
}

func (g GeneratorDynamicAttribute) To() (generatorschema.ToFromConversion, error) {
	if g.AssociatedExternalType != nil {
		return generatorschema.ToFromConversion{
			AssocExtType: g.AssociatedExternalType,
		}, nil
	}

	return generatorschema.ToFromConversion{
		Default: "ValueStringPointer",
	}, nil
}

func (g GeneratorDynamicAttribute) From() (generatorschema.ToFromConversion, error) {
	if g.AssociatedExternalType != nil {
		return generatorschema.ToFromConversion{
			AssocExtType: g.AssociatedExternalType,
		}, nil
	}

	return generatorschema.ToFromConversion{
		Default: "StringPointerValue",
	}, nil
}
