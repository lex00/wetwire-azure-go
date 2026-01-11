package serialize

import (
	"testing"

	"github.com/lex00/wetwire-azure-go/intrinsics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSerializeValue_NonIntrinsic tests SerializeValue with non-intrinsic values
func TestSerializeValue_NonIntrinsic(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{"string", "hello", "hello"},
		{"int", 42, 42},
		{"float", 3.14, 3.14},
		{"bool", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SerializeValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsZeroValue_AllTypes tests the isZeroValue function with all supported types
func TestIsZeroValue_AllTypes(t *testing.T) {
	// Test struct with various field types
	type AllFieldTypes struct {
		StringField    string          `json:"stringField"`
		IntField       int             `json:"intField"`
		Int8Field      int8            `json:"int8Field"`
		Int16Field     int16           `json:"int16Field"`
		Int32Field     int32           `json:"int32Field"`
		Int64Field     int64           `json:"int64Field"`
		UintField      uint            `json:"uintField"`
		Uint8Field     uint8           `json:"uint8Field"`
		Uint16Field    uint16          `json:"uint16Field"`
		Uint32Field    uint32          `json:"uint32Field"`
		Uint64Field    uint64          `json:"uint64Field"`
		Float32Field   float32         `json:"float32Field"`
		Float64Field   float64         `json:"float64Field"`
		BoolField      bool            `json:"boolField"`
		SliceField     []string        `json:"sliceField,omitempty"`
		MapField       map[string]any  `json:"mapField,omitempty"`
		PointerField   *string         `json:"pointerField,omitempty"`
		InterfaceField any             `json:"interfaceField,omitempty"`
	}

	// Test zero values
	zero := AllFieldTypes{}
	result := ToARMResource(zero)

	// Zero values should result in zero values in output (not omitted unless omitempty)
	assert.Equal(t, "", result["stringField"])
	assert.Equal(t, 0, result["intField"])
	assert.Equal(t, int8(0), result["int8Field"])
	assert.Equal(t, int16(0), result["int16Field"])
	assert.Equal(t, int32(0), result["int32Field"])
	assert.Equal(t, int64(0), result["int64Field"])
	assert.Equal(t, uint(0), result["uintField"])
	assert.Equal(t, uint8(0), result["uint8Field"])
	assert.Equal(t, uint16(0), result["uint16Field"])
	assert.Equal(t, uint32(0), result["uint32Field"])
	assert.Equal(t, uint64(0), result["uint64Field"])
	assert.Equal(t, float32(0), result["float32Field"])
	assert.Equal(t, float64(0), result["float64Field"])
	assert.Equal(t, false, result["boolField"])

	// omitempty fields should be omitted when zero
	assert.Nil(t, result["sliceField"])
	assert.Nil(t, result["mapField"])
	assert.Nil(t, result["pointerField"])
	assert.Nil(t, result["interfaceField"])
}

// TestIsZeroValue_NonZeroValues tests isZeroValue with non-zero values
func TestIsZeroValue_NonZeroValues(t *testing.T) {
	type NonZeroStruct struct {
		StringField  string         `json:"stringField"`
		IntField     int            `json:"intField"`
		BoolField    bool           `json:"boolField"`
		Float64Field float64        `json:"float64Field"`
		SliceField   []string       `json:"sliceField,omitempty"`
		MapField     map[string]any `json:"mapField,omitempty"`
	}

	nonZero := NonZeroStruct{
		StringField:  "hello",
		IntField:     42,
		BoolField:    true,
		Float64Field: 3.14,
		SliceField:   []string{"a", "b"},
		MapField:     map[string]any{"key": "value"},
	}

	result := ToARMResource(nonZero)

	assert.Equal(t, "hello", result["stringField"])
	assert.Equal(t, 42, result["intField"])
	assert.Equal(t, true, result["boolField"])
	assert.Equal(t, 3.14, result["float64Field"])

	slice, ok := result["sliceField"].([]any)
	require.True(t, ok)
	assert.Len(t, slice, 2)

	mapField, ok := result["mapField"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "value", mapField["key"])
}

// TestStructToMap_NilPointer tests structToMap with nil pointer
func TestStructToMap_NilPointer(t *testing.T) {
	var nilPointer *struct{ Name string }
	result := ToARMResource(nilPointer)
	assert.Nil(t, result)
}

// TestStructToMap_NonStruct tests structToMap with non-struct type
func TestStructToMap_NonStruct(t *testing.T) {
	result := ToARMResource("not a struct")
	assert.Nil(t, result)
}

// TestConvertValue_NilPointer tests convertValue with nil pointer
func TestConvertValue_NilPointer(t *testing.T) {
	type WithPointer struct {
		Name *string `json:"name"`
	}

	s := WithPointer{Name: nil}
	result := ToARMResource(s)
	assert.Nil(t, result["name"])
}

// TestConvertValue_NilSlice tests convertValue with nil slice
func TestConvertValue_NilSlice(t *testing.T) {
	type WithSlice struct {
		Items []string `json:"items"`
	}

	s := WithSlice{Items: nil}
	result := ToARMResource(s)
	assert.Nil(t, result["items"])
}

// TestConvertValue_NilMap tests convertValue with nil map
func TestConvertValue_NilMap(t *testing.T) {
	type WithMap struct {
		Data map[string]string `json:"data"`
	}

	s := WithMap{Data: nil}
	result := ToARMResource(s)
	assert.Nil(t, result["data"])
}

// TestConvertValue_NilInterface tests convertValue with nil interface
func TestConvertValue_NilInterface(t *testing.T) {
	type WithInterface struct {
		Value any `json:"value"`
	}

	s := WithInterface{Value: nil}
	result := ToARMResource(s)
	assert.Nil(t, result["value"])
}

// TestConvertValue_Interface tests convertValue with non-nil interface
func TestConvertValue_Interface(t *testing.T) {
	type WithInterface struct {
		Value any `json:"value"`
	}

	s := WithInterface{Value: "hello"}
	result := ToARMResource(s)
	assert.Equal(t, "hello", result["value"])
}

// TestConvertValue_NestedStruct tests convertValue with nested struct
func TestConvertValue_NestedStruct(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner Inner `json:"inner"`
	}

	s := Outer{Inner: Inner{Name: "test"}}
	result := ToARMResource(s)

	inner, ok := result["inner"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", inner["name"])
}

// TestConvertValue_MapWithValues tests convertValue with populated map
func TestConvertValue_MapWithValues(t *testing.T) {
	type WithMap struct {
		Tags map[string]string `json:"tags"`
	}

	s := WithMap{
		Tags: map[string]string{
			"env":  "prod",
			"team": "platform",
		},
	}
	result := ToARMResource(s)

	tags, ok := result["tags"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "prod", tags["env"])
	assert.Equal(t, "platform", tags["team"])
}

// TestConvertValue_SliceWithValues tests convertValue with populated slice
func TestConvertValue_SliceWithValues(t *testing.T) {
	type WithSlice struct {
		Items []string `json:"items"`
	}

	s := WithSlice{Items: []string{"a", "b", "c"}}
	result := ToARMResource(s)

	items, ok := result["items"].([]any)
	require.True(t, ok)
	assert.Len(t, items, 3)
	assert.Equal(t, "a", items[0])
	assert.Equal(t, "b", items[1])
	assert.Equal(t, "c", items[2])
}

// TestConvertValue_SliceOfStructs tests convertValue with slice of structs
func TestConvertValue_SliceOfStructs(t *testing.T) {
	type Item struct {
		Name string `json:"name"`
	}
	type WithSlice struct {
		Items []Item `json:"items"`
	}

	s := WithSlice{
		Items: []Item{
			{Name: "item1"},
			{Name: "item2"},
		},
	}
	result := ToARMResource(s)

	items, ok := result["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 2)

	item1, ok := items[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "item1", item1["name"])
}

// TestParseJSONTag_EmptyTag tests parseJSONTag with empty tag
func TestParseJSONTag_EmptyTag(t *testing.T) {
	type WithEmptyTag struct {
		Name string `json:""`
	}

	s := WithEmptyTag{Name: "test"}
	result := ToARMResource(s)
	// Empty json tag should be skipped
	assert.Empty(t, result)
}

// TestParseJSONTag_DashTag tests parseJSONTag with "-" tag
func TestParseJSONTag_DashTag(t *testing.T) {
	type WithDashTag struct {
		Name   string `json:"-"`
		Public string `json:"public"`
	}

	s := WithDashTag{Name: "hidden", Public: "visible"}
	result := ToARMResource(s)
	assert.Nil(t, result["Name"])
	assert.Nil(t, result["-"])
	assert.Equal(t, "visible", result["public"])
}

// TestParseJSONTag_OmitEmpty tests parseJSONTag with omitempty
func TestParseJSONTag_OmitEmpty(t *testing.T) {
	type WithOmitEmpty struct {
		Required string `json:"required"`
		Optional string `json:"optional,omitempty"`
	}

	s := WithOmitEmpty{Required: "value", Optional: ""}
	result := ToARMResource(s)
	assert.Equal(t, "value", result["required"])
	assert.Nil(t, result["optional"])
}

// TestUnexportedFields tests that unexported fields are skipped
func TestUnexportedFields(t *testing.T) {
	type WithUnexported struct {
		Public  string `json:"public"`
		private string //nolint:unused // intentionally testing unexported field handling
	}

	s := WithUnexported{Public: "visible"}
	result := ToARMResource(s)
	assert.Equal(t, "visible", result["public"])
	assert.Nil(t, result["private"])
}

// TestIntrinsicInStruct tests intrinsic detection within struct
func TestIntrinsicInStruct(t *testing.T) {
	type ResourceWithIntrinsic struct {
		Name     string                    `json:"name"`
		Location intrinsics.ResourceGroupValue `json:"location"`
	}

	s := ResourceWithIntrinsic{
		Name:     "test",
		Location: intrinsics.ResourceGroupValue{Property: "location"},
	}
	result := ToARMResource(s)

	assert.Equal(t, "test", result["name"])
	assert.Equal(t, "[resourceGroup().location]", result["location"])
}

// TestAllIntrinsicTypes tests all intrinsic type serialization
func TestAllIntrinsicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    intrinsics.Intrinsic
		expected string
	}{
		{
			name:     "ResourceId",
			input:    intrinsics.ResourceId("Microsoft.Storage/storageAccounts", "myaccount"),
			expected: "[resourceId('Microsoft.Storage/storageAccounts', 'myaccount')]",
		},
		{
			name:     "Reference",
			input:    intrinsics.Ref("myresource", "2021-04-01"),
			expected: "[reference('myresource', '2021-04-01')]",
		},
		{
			name:     "ReferenceWithProperty",
			input:    intrinsics.RefProperty("myresource", "2021-04-01", "primaryEndpoints.blob"),
			expected: "[reference('myresource', '2021-04-01').primaryEndpoints.blob]",
		},
		{
			name:     "Parameter",
			input:    intrinsics.Parameters("location"),
			expected: "[parameters('location')]",
		},
		{
			name:     "Variable",
			input:    intrinsics.Variables("storageAccountName"),
			expected: "[variables('storageAccountName')]",
		},
		{
			name:     "ResourceGroup",
			input:    intrinsics.ResourceGroup(),
			expected: "[resourceGroup()]",
		},
		{
			name:     "ResourceGroupWithProperty",
			input:    intrinsics.ResourceGroupValue{Property: "location"},
			expected: "[resourceGroup().location]",
		},
		{
			name:     "Subscription",
			input:    intrinsics.Subscription{},
			expected: "[subscription()]",
		},
		{
			name:     "SubscriptionWithProperty",
			input:    intrinsics.Subscription{Property: "subscriptionId"},
			expected: "[subscription().subscriptionId]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SerializeValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestConcat tests Concat intrinsic serialization
func TestConcat(t *testing.T) {
	concat := intrinsics.Concat{Values: []any{"a", "b"}}
	result := SerializeValue(concat)
	assert.Equal(t, "[concat(...)]", result)
}

// TestUniqueString tests UniqueString intrinsic serialization
func TestUniqueString(t *testing.T) {
	unique := intrinsics.UniqueString{Values: []string{"a", "b"}}
	result := SerializeValue(unique)
	assert.Equal(t, "[uniqueString(...)]", result)
}

// TestEmptyARMTemplate tests ARM template with no resources
func TestEmptyARMTemplate(t *testing.T) {
	template := ToARMTemplate([]any{})

	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", template["$schema"])
	assert.Equal(t, "1.0.0.0", template["contentVersion"])

	resources, ok := template["resources"].([]any)
	require.True(t, ok)
	assert.Empty(t, resources)
}

// TestDeepNestedStruct tests deeply nested struct serialization
func TestDeepNestedStruct(t *testing.T) {
	type Level3 struct {
		Value string `json:"value"`
	}
	type Level2 struct {
		Level3 Level3 `json:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `json:"level2"`
	}
	type Root struct {
		Level1 Level1 `json:"level1"`
	}

	s := Root{
		Level1: Level1{
			Level2: Level2{
				Level3: Level3{
					Value: "deep",
				},
			},
		},
	}

	result := ToARMResource(s)

	level1, ok := result["level1"].(map[string]any)
	require.True(t, ok)

	level2, ok := level1["level2"].(map[string]any)
	require.True(t, ok)

	level3, ok := level2["level3"].(map[string]any)
	require.True(t, ok)

	assert.Equal(t, "deep", level3["value"])
}

// TestPointerToStruct tests pointer to struct serialization
func TestPointerToStruct(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner *Inner `json:"inner"`
	}

	inner := &Inner{Name: "test"}
	s := Outer{Inner: inner}
	result := ToARMResource(s)

	innerResult, ok := result["inner"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", innerResult["name"])
}

// TestSliceOfPointers tests slice of pointers serialization
func TestSliceOfPointers(t *testing.T) {
	type Item struct {
		Name string `json:"name"`
	}
	type Container struct {
		Items []*Item `json:"items"`
	}

	item1 := &Item{Name: "item1"}
	item2 := &Item{Name: "item2"}
	s := Container{Items: []*Item{item1, item2}}
	result := ToARMResource(s)

	items, ok := result["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 2)

	i1, ok := items[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "item1", i1["name"])
}

// TestMapWithIntrinsicValues tests map with intrinsic values
func TestMapWithIntrinsicValues(t *testing.T) {
	type Container struct {
		Config map[string]any `json:"config"`
	}

	s := Container{
		Config: map[string]any{
			"location": intrinsics.ResourceGroupValue{Property: "location"},
			"static":   "value",
		},
	}
	result := ToARMResource(s)

	config, ok := result["config"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "[resourceGroup().location]", config["location"])
	assert.Equal(t, "value", config["static"])
}

// TestSliceWithIntrinsicValues tests slice with intrinsic values
func TestSliceWithIntrinsicValues(t *testing.T) {
	type Container struct {
		Values []any `json:"values"`
	}

	s := Container{
		Values: []any{
			intrinsics.Parameters("param1"),
			"static",
			intrinsics.Variables("var1"),
		},
	}
	result := ToARMResource(s)

	values, ok := result["values"].([]any)
	require.True(t, ok)
	require.Len(t, values, 3)
	assert.Equal(t, "[parameters('param1')]", values[0])
	assert.Equal(t, "static", values[1])
	assert.Equal(t, "[variables('var1')]", values[2])
}

// TestPointerToPointer tests double pointer handling
func TestPointerToPointer(t *testing.T) {
	name := "test"
	namePtr := &name
	type Container struct {
		Name **string `json:"name"`
	}

	s := Container{Name: &namePtr}
	result := ToARMResource(s)
	// Double pointer should be dereferenced
	assert.Equal(t, "test", result["name"])
}

// TestIsZeroValue_Struct tests isZeroValue with zero struct
func TestIsZeroValue_Struct(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner Inner `json:"inner,omitempty"`
	}

	s := Outer{Inner: Inner{}}
	result := ToARMResource(s)
	// Zero struct with omitempty should be omitted
	assert.Nil(t, result["inner"])
}

// TestIsZeroValue_NonZeroStruct tests isZeroValue with non-zero struct
func TestIsZeroValue_NonZeroStruct(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner Inner `json:"inner,omitempty"`
	}

	s := Outer{Inner: Inner{Name: "test"}}
	result := ToARMResource(s)

	inner, ok := result["inner"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", inner["name"])
}

// TestSliceType tests slice type handling
func TestSliceType(t *testing.T) {
	type Container struct {
		Values []string `json:"values"`
	}

	s := Container{Values: []string{"a", "b", "c"}}
	result := ToARMResource(s)

	values, ok := result["values"].([]any)
	require.True(t, ok)
	require.Len(t, values, 3)
	assert.Equal(t, "a", values[0])
	assert.Equal(t, "b", values[1])
	assert.Equal(t, "c", values[2])
}

// TestEmptySlice tests empty slice handling
func TestEmptySlice(t *testing.T) {
	type Container struct {
		Values []string `json:"values"`
	}

	s := Container{Values: []string{}}
	result := ToARMResource(s)
	// Empty slice should be nil/omitted
	assert.Nil(t, result["values"])
}

// TestIsZeroValue_Uint tests isZeroValue with uint types
func TestIsZeroValue_Uint(t *testing.T) {
	type UintStruct struct {
		U   uint   `json:"u,omitempty"`
		U8  uint8  `json:"u8,omitempty"`
		U16 uint16 `json:"u16,omitempty"`
		U32 uint32 `json:"u32,omitempty"`
		U64 uint64 `json:"u64,omitempty"`
	}

	// Zero values with omitempty - should be omitted
	zero := UintStruct{}
	result := ToARMResource(zero)
	assert.Nil(t, result["u"])
	assert.Nil(t, result["u8"])
	assert.Nil(t, result["u16"])
	assert.Nil(t, result["u32"])
	assert.Nil(t, result["u64"])

	// Non-zero values - should be present
	nonZero := UintStruct{U: 1, U8: 2, U16: 3, U32: 4, U64: 5}
	result = ToARMResource(nonZero)
	assert.Equal(t, uint(1), result["u"])
	assert.Equal(t, uint8(2), result["u8"])
	assert.Equal(t, uint16(3), result["u16"])
	assert.Equal(t, uint32(4), result["u32"])
	assert.Equal(t, uint64(5), result["u64"])
}

// TestIsZeroValue_Float tests isZeroValue with float types
func TestIsZeroValue_Float(t *testing.T) {
	type FloatStruct struct {
		F32 float32 `json:"f32,omitempty"`
		F64 float64 `json:"f64,omitempty"`
	}

	// Zero values with omitempty - should be omitted
	zero := FloatStruct{}
	result := ToARMResource(zero)
	assert.Nil(t, result["f32"])
	assert.Nil(t, result["f64"])

	// Non-zero values - should be present
	nonZero := FloatStruct{F32: 1.5, F64: 2.5}
	result = ToARMResource(nonZero)
	assert.Equal(t, float32(1.5), result["f32"])
	assert.Equal(t, float64(2.5), result["f64"])
}

// TestIsZeroValue_SliceVariants tests isZeroValue with various slice states
func TestIsZeroValue_SliceVariants(t *testing.T) {
	type SliceStruct struct {
		Empty    []string `json:"empty,omitempty"`
		NonEmpty []string `json:"nonEmpty,omitempty"`
	}

	s := SliceStruct{
		Empty:    []string{},
		NonEmpty: []string{"a"},
	}
	result := ToARMResource(s)
	assert.Nil(t, result["empty"])
	assert.NotNil(t, result["nonEmpty"])
}

// TestIsZeroValue_Map tests isZeroValue with map types
func TestIsZeroValue_Map(t *testing.T) {
	type MapStruct struct {
		Empty    map[string]string `json:"empty,omitempty"`
		NonEmpty map[string]string `json:"nonEmpty,omitempty"`
	}

	s := MapStruct{
		Empty:    map[string]string{},
		NonEmpty: map[string]string{"k": "v"},
	}
	result := ToARMResource(s)
	assert.Nil(t, result["empty"])
	assert.NotNil(t, result["nonEmpty"])
}

// TestIsZeroValue_Pointer tests isZeroValue with pointer types
func TestIsZeroValue_Pointer(t *testing.T) {
	type PointerStruct struct {
		Nil    *string `json:"nil,omitempty"`
		NonNil *string `json:"nonNil,omitempty"`
	}

	str := "value"
	s := PointerStruct{
		Nil:    nil,
		NonNil: &str,
	}
	result := ToARMResource(s)
	assert.Nil(t, result["nil"])
	assert.Equal(t, "value", result["nonNil"])
}

// TestIsZeroValue_Interface tests isZeroValue with interface types
func TestIsZeroValue_Interface(t *testing.T) {
	type InterfaceStruct struct {
		Nil    any `json:"nil,omitempty"`
		NonNil any `json:"nonNil,omitempty"`
	}

	s := InterfaceStruct{
		Nil:    nil,
		NonNil: "value",
	}
	result := ToARMResource(s)
	assert.Nil(t, result["nil"])
	assert.Equal(t, "value", result["nonNil"])
}

// TestIsZeroValue_Bool tests isZeroValue with bool type
func TestIsZeroValue_Bool(t *testing.T) {
	type BoolStruct struct {
		False bool `json:"false,omitempty"`
		True  bool `json:"true,omitempty"`
	}

	s := BoolStruct{
		False: false,
		True:  true,
	}
	result := ToARMResource(s)
	assert.Nil(t, result["false"])
	assert.Equal(t, true, result["true"])
}

// TestParseJSONTag_NoComma tests parseJSONTag with no comma in tag
func TestParseJSONTag_NoComma(t *testing.T) {
	type SimpleTag struct {
		Name string `json:"name"`
	}

	s := SimpleTag{Name: "test"}
	result := ToARMResource(s)
	assert.Equal(t, "test", result["name"])
}
