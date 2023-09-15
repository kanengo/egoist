package metadata

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dapr/kit/ptr"
	"github.com/kanengo/egoist/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

const (
	RawPayloadKey = "rawPayload"

	ContentType = "contentType"
)

func DecodeMetadata(input any, result any) error {
	// avoids a common mistake of passing the metadata struct, instead of the properties map
	// if input is of type struct, cast it to metadata.Base and access the Properties instead
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Properties")
		if f.IsValid() && f.Kind() == reflect.Map {
			input = f.Interface().(map[string]string)
		}
	}

	inputMap, err := cast.ToStringMapStringE(input)
	if err != nil {
		return fmt.Errorf("input object cannot be cast to map[string]string: %w", err)
	}

	// Handle aliases
	err = resolveAliases(inputMap, result)
	if err != nil {
		return fmt.Errorf("failed to resolve aliases: %w", err)
	}

	// Finally, decode the metadata using mapstructure
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeDurationArrayHookFunc(),
			toTimeDurationHookFunc(),
			toTruthyBoolHookFunc(),
			toStringArrayHookFunc(),
			toByteSizeHookFunc(),
		),
		Metadata:         nil,
		Result:           result,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(inputMap)
	return err
}

func resolveAliases(md map[string]string, result any) error {
	// Get the list of all keys in the map
	keys := make(map[string]string, len(md))
	for k := range md {
		lk := strings.ToLower(k)

		// Check if there are duplicate keys after lowercasing
		_, ok := keys[lk]
		if ok {
			return fmt.Errorf("key %s is duplicate in the metadata", lk)
		}

		keys[lk] = k
	}

	// Error if result is not pointer to struct, or pointer to pointer to struct
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("not a pointer: %s", t.Kind().String())
	}
	t = t.Elem()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", t.Kind().String())
	}

	// Iterate through all the properties of result to see if anyone has the "mapstructurealiases" property
	for i := 0; i < t.NumField(); i++ {
		currentField := t.Field(i)

		// Ignored fields that are not exported or that don't have a "mapstructure" tag
		mapstructureTag := currentField.Tag.Get("mapstructure")
		if !currentField.IsExported() || mapstructureTag == "" {
			continue
		}

		// If the current property has a value in the metadata, then we don't need to handle aliases
		_, ok := keys[strings.ToLower(mapstructureTag)]
		if ok {
			continue
		}

		// Check if there's a "mapstructurealiases" tag
		aliasesTag := strings.ToLower(currentField.Tag.Get("mapstructurealiases"))
		if aliasesTag == "" {
			continue
		}

		// Look for the first alias that has a value
		var mdKey string
		for _, alias := range strings.Split(aliasesTag, ",") {
			mdKey, ok = keys[alias]
			if !ok {
				continue
			}

			// We found an alias
			md[mapstructureTag] = md[mdKey]
			break
		}
	}

	return nil
}

func toTruthyBoolHookFunc() mapstructure.DecodeHookFunc {
	stringType := reflect.TypeOf("")
	boolType := reflect.TypeOf(true)
	boolPtrType := reflect.TypeOf(ptr.Of(true))

	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f == stringType && t == boolType {
			return utils.IsTruthy(data.(string)), nil
		}
		if f == stringType && t == boolPtrType {
			return ptr.Of(utils.IsTruthy(data.(string))), nil
		}
		return data, nil
	}
}

func toStringArrayHookFunc() mapstructure.DecodeHookFunc {
	stringType := reflect.TypeOf("")
	stringSliceType := reflect.TypeOf([]string{})
	stringSlicePtrType := reflect.TypeOf(ptr.Of([]string{}))

	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f == stringType && t == stringSliceType {
			return strings.Split(data.(string), ","), nil
		}
		if f == stringType && t == stringSlicePtrType {
			return ptr.Of(strings.Split(data.(string), ",")), nil
		}
		return data, nil
	}
}

func toTimeDurationArrayHookFunc() mapstructure.DecodeHookFunc {
	convert := func(input string) ([]time.Duration, error) {
		parts := strings.Split(input, ",")
		res := make([]time.Duration, 0, len(parts))
		for _, v := range parts {
			input := strings.TrimSpace(v)
			if input == "" {
				continue
			}
			val, err := time.ParseDuration(input)
			if err != nil {
				// If we can't parse the duration, try parsing it as int64 seconds
				seconds, errParse := strconv.ParseInt(input, 10, 0)
				if errParse != nil {
					return nil, errors.Join(err, errParse)
				}
				val = time.Duration(seconds * int64(time.Second))
			}
			res = append(res, val)
		}
		return res, nil
	}

	stringType := reflect.TypeOf("")
	durationSliceType := reflect.TypeOf([]time.Duration{})
	durationSlicePtrType := reflect.TypeOf(ptr.Of([]time.Duration{}))

	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f == stringType && t == durationSliceType {
			inputArrayString := data.(string)
			return convert(inputArrayString)
		}
		if f == stringType && t == durationSlicePtrType {
			inputArrayString := data.(string)
			res, err := convert(inputArrayString)
			if err != nil {
				return nil, err
			}
			return ptr.Of(res), nil
		}
		return data, nil
	}
}

func IsRawPayload(props map[string]string) (bool, error) {
	if val, ok := props[RawPayloadKey]; ok && val != "" {
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("%s value must be a valid boolean: actual is '%s'", RawPayloadKey, val)
		}

		return boolVal, nil
	}

	return false, nil
}

func TryGetContentType(props map[string]string) (string, bool) {
	if val, ok := props[ContentType]; ok && val != "" {
		return val, true
	}

	return "", false
}

type ComponentType string

const (
	BindingType            ComponentType = "bindings"
	StateStoreType         ComponentType = "state"
	SecretStoreType        ComponentType = "secretstores"
	PubSubType             ComponentType = "pubsub"
	LockStoreType          ComponentType = "lock"
	ConfigurationStoreType ComponentType = "configuration"
	MiddlewareType         ComponentType = "middleware"
	CryptoType             ComponentType = "crypto"
	NameResolutionType     ComponentType = "nameresolution"
	WorkflowType           ComponentType = "workflows"
)

// IsValid returns true if the component type is valid.
func (t ComponentType) IsValid() bool {
	switch t {
	case BindingType, StateStoreType,
		SecretStoreType, PubSubType,
		LockStoreType, ConfigurationStoreType,
		MiddlewareType, CryptoType,
		NameResolutionType, WorkflowType:
		return true
	default:
		return false
	}
}

// BuiltInMetadataProperties returns the built-in metadata properties for the given component type.
// These are normally parsed by the runtime.
func (t ComponentType) BuiltInMetadataProperties() []string {
	switch t {
	case StateStoreType:
		return []string{
			"actorStateStore",
			"keyPrefix",
		}
	case LockStoreType:
		return []string{
			"keyPrefix",
		}
	default:
		return nil
	}
}

type MetadataField struct {
	// Field type
	Type string
	// True if the field should be ignored by the metadata analyzer
	Ignored bool
	// True if the field is deprecated
	Deprecated bool
	// Aliases used for old, deprecated names
	Aliases []string
}

type MetadataMap map[string]MetadataField

// GetMetadataInfoFromStructType converts a struct to a map of field name (or struct tag) to field type.
// This is used to generate metadata documentation for components.
func GetMetadataInfoFromStructType(t reflect.Type, metadataMap *MetadataMap, componentType ComponentType) error {
	// Return if not struct or pointer to struct.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", t.Kind().String())
	}

	if *metadataMap == nil {
		*metadataMap = MetadataMap{}
	}

	for i := 0; i < t.NumField(); i++ {
		currentField := t.Field(i)
		// fields that are not exported cannot be set via the mapstructure metadata decoding mechanism
		if !currentField.IsExported() {
			continue
		}
		mapStructureTag := currentField.Tag.Get("mapstructure")
		// we are not exporting this field using the mapstructure tag mechanism
		if mapStructureTag == "-" {
			continue
		}

		// If there's a "mdonly" tag, that metadata option is only included for certain component types
		if mdOnlyTag := currentField.Tag.Get("mdonly"); mdOnlyTag != "" {
			include := false
			onlyTags := strings.Split(mdOnlyTag, ",")
			for _, tag := range onlyTags {
				if tag == string(componentType) {
					include = true
					break
				}
			}
			if !include {
				continue
			}
		}

		mdField := MetadataField{
			Type: currentField.Type.String(),
		}

		// If there's a mdignore tag and that's truthy, the field should be ignored by the metadata analyzer
		mdField.Ignored = utils.IsTruthy(currentField.Tag.Get("mdignore"))

		// If there's a "mddeprecated" tag, the field may be deprecated
		mdField.Deprecated = utils.IsTruthy(currentField.Tag.Get("mddeprecated"))

		// If there's a "mdaliases" tag, the field contains aliases
		// The value is a comma-separated string
		if mdAliasesTag := currentField.Tag.Get("mdaliases"); mdAliasesTag != "" {
			mdField.Aliases = strings.Split(mdAliasesTag, ",")
		}

		// Handle mapstructure tags and get the field name
		mapStructureTags := strings.Split(mapStructureTag, ",")
		numTags := len(mapStructureTags)
		if numTags > 1 && mapStructureTags[numTags-1] == "squash" && currentField.Anonymous {
			// traverse embedded struct
			_ = GetMetadataInfoFromStructType(currentField.Type, metadataMap, componentType)
			continue
		}
		var fieldName string
		if numTags > 0 && mapStructureTags[0] != "" {
			fieldName = mapStructureTags[0]
		} else {
			fieldName = currentField.Name
		}

		// Add the field
		(*metadataMap)[fieldName] = mdField
	}

	return nil
}
