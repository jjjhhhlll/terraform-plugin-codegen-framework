package convert

import (
	"fmt"

	generatorschema "github.com/hashicorp/terraform-plugin-codegen-framework/internal/schema"
	specschema "github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

type DefaultDynamic struct {
	dynamicDefault *specschema.DynamicDefault
}

func NewDefaultDynamic(b *specschema.DynamicDefault) DefaultDynamic {
	return DefaultDynamic{
		dynamicDefault: b,
	}
}

func (d DefaultDynamic) Equal(other DefaultDynamic) bool {
	return d.dynamicDefault.Equal(other.dynamicDefault)
}

func (d DefaultDynamic) Imports() *generatorschema.Imports {
	imports := generatorschema.NewImports()

	if d.dynamicDefault == nil {
		return imports
	}

	if d.dynamicDefault.Custom != nil {
		for _, i := range d.dynamicDefault.Custom.Imports {
			if len(i.Path) > 0 {
				imports.Add(i)
			}
		}
	}

	return imports
}

func (d DefaultDynamic) Schema() []byte {
	if d.dynamicDefault == nil {
		return nil
	}

	if d.dynamicDefault.Custom != nil && d.dynamicDefault.Custom.SchemaDefinition != "" {
		return []byte(fmt.Sprintf("Default: %s,\n", d.dynamicDefault.Custom.SchemaDefinition))
	}

	return nil
}
