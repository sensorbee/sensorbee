package data

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	multierror "github.com/hashicorp/go-multierror"
)

// Decoder decodes a Map into a struct.
type Decoder struct {
	config *DecoderConfig

	// This Decoder is strongly inspired by github.com/mitchellh/mapstructure.
	// The decoder may be going to support "json" compatible tags in structs.
}

// DecoderConfig is used to configure the behavior of Decoder.
type DecoderConfig struct {
	// ErrorUnused, if set to true, reports an error when a Map contains keys
	// that are not defined in a struct.
	ErrorUnused bool

	// Metadata has meta information of decode. If this is nil, meat information
	// will not be tracked.
	Metadata *DecoderMetadata

	TagName string
}

// DecoderMetadata tracks field names that are used or not used for decoding.
type DecoderMetadata struct {
	// Keys contains keys in a Map that are processed.
	Keys []string

	// Unsed contains keys in a Map that are not defined in a struct.
	Unused []string
}

// NewDecoder creates a new Decoder with the given config.
func NewDecoder(c *DecoderConfig) *Decoder {
	if c == nil {
		c = &DecoderConfig{}
	}
	if c.TagName == "" {
		c.TagName = "bql"
	}
	return &Decoder{
		config: c,
	}
}

// Decode decodes a Map into a struct. The argument must be a pointer to a
// struct.
func (d *Decoder) Decode(m Map, v interface{}) (err error) {
	defer func() {
		// Because this function heavily depends on reflect and it has many
		// chance to panic, this function catches it converts it to an error.
		// This isn't a desired handling of panics, but panics makes debugging
		// hard for SensorBee users otherwise.
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	p := reflect.ValueOf(v)
	if p.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer to a struct")
	}
	s := p.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("result must be pointer to a struct")
	}
	return d.decodeStruct("", m, s)
}

func (d *Decoder) decode(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	switch dst.Kind() {
	case reflect.Bool:
		return d.decodeBool(prefix, src, dst, weaklyTyped)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// TODO: support detecting overflows
		return d.decodeInt(prefix, src, dst, weaklyTyped)

	case reflect.Float32, reflect.Float64:
		return d.decodeFloat(prefix, src, dst, weaklyTyped)

	case reflect.String:
		return d.decodeString(prefix, src, dst, weaklyTyped)

	case reflect.Interface: // Only Value is supported
		if dst.Type() != reflect.TypeOf(func(Value) {}).In(0) {
			return fmt.Errorf("%v: interface{} other than data.Value is not supported", prefix)
		}
		return d.decodeValue(prefix, src, dst)

	case reflect.Map:
		return d.decodeMap(prefix, src, dst, weaklyTyped)

	case reflect.Slice:
		return d.decodeSlice(prefix, src, dst, weaklyTyped)

	case reflect.Struct:
		return d.decodeStruct(prefix, src, dst)

	case reflect.Ptr:
		// To decode a value to dst, dst must be addressable. However,
		// reflect.ValueOf or reflect.Zero may return non-addressable values
		// especially when value is of primitive types. The following
		// New-Indirect idiom works for such cases. First reflect.New creates
		// a pointer that points to a non-nil addressable value. Then,
		// reflect.Indirect returns an element pointed by the pointer.
		v := reflect.New(dst.Type().Elem())
		if err := d.decode(prefix, src, reflect.Indirect(v), weaklyTyped); err != nil {
			return err
		}
		dst.Set(v)
		return nil
	}
	return fmt.Errorf("%v: decoder doesn't support the type: %v", prefix, dst.Kind())
}

func (d *Decoder) decodeBool(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		b   bool
		err error
	)
	if weaklyTyped {
		b, err = ToBool(src)
	} else {
		b, err = AsBool(src)
	}
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.SetBool(b)
	return nil
}

func (d *Decoder) decodeInt(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	if _, ok := dst.Interface().(time.Duration); ok {
		return d.decodeDuration(prefix, src, dst)
	}

	var (
		i   int64
		err error
	)

	if weaklyTyped {
		i, err = ToInt(src)
	} else {
		i, err = AsInt(src)

		// special handling for floats having integer values
		if err != nil && src.Type() == TypeFloat {
			f, _ := AsFloat(src)
			i = int64(f)
			if f == float64(i) {
				err = nil
			}
		}
	}
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.SetInt(i)
	return nil
}

func (d *Decoder) decodeDuration(prefix string, src Value, dst reflect.Value) error {
	// As described in decodeTimestamp, decodeDuration also assumes that the
	// value is always weaklytyped.
	dur, err := ToDuration(src)
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.Set(reflect.ValueOf(dur))
	return nil
}

func (d *Decoder) decodeFloat(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		f   float64
		err error
	)
	if weaklyTyped {
		f, err = ToFloat(src)
	} else {
		f, err = AsFloat(src)
		if err != nil && src.Type() == TypeInt {
			i, _ := AsInt(src)
			f = float64(i)
			err = nil
		}
	}
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.SetFloat(f)
	return nil
}

func (d *Decoder) decodeString(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		s   string
		err error
	)
	if weaklyTyped {
		s, err = ToString(src)
	} else {
		s, err = AsString(src)
	}
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.SetString(s)
	return nil
}

func (d *Decoder) decodeValue(prefix string, src Value, dst reflect.Value) error {
	dst.Set(reflect.ValueOf(src))
	return nil
}

func (d *Decoder) decodeMap(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	if src.Type() != TypeMap {
		return fmt.Errorf("%v: cannot decode to a map: %v", prefix, src.Type())
	}
	m, _ := AsMap(src)

	t := dst.Type()
	if k := t.Key().Kind(); k != reflect.String {
		return fmt.Errorf("%v: key must be string: %v", prefix, k)
	}
	valueType := t.Elem()

	var errs *multierror.Error
	res := reflect.MakeMap(t)
	for k, e := range m {
		v := reflect.Indirect(reflect.New(valueType))
		if err := d.decode(fmt.Sprintf(`%v["%v"]`, prefix, k), e, v, weaklyTyped); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		res.SetMapIndex(reflect.ValueOf(k), v)
	}
	dst.Set(res)
	if errs == nil {
		return nil // DO NOT return errs even if it's nil
	}
	return errs
}

func (d *Decoder) decodeSlice(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	if dst.Type().Elem().Kind() == reflect.Uint8 {
		return d.decodeBlob(prefix, src, dst, weaklyTyped)
	}

	if src.Type() != TypeArray {
		return fmt.Errorf("%v: cannot decode to an array: %v", prefix, src.Type())
	}
	a, _ := AsArray(src)

	var errs *multierror.Error
	res := reflect.MakeSlice(dst.Type(), len(a), len(a))
	for i, e := range a {
		v := res.Index(i)
		if err := d.decode(fmt.Sprintf("%v[%v]", prefix, i), e, v, weaklyTyped); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
	}
	dst.Set(res)
	if errs == nil {
		return nil // DO NOT return errs even if it's nil
	}
	return errs
}

func (d *Decoder) decodeBlob(prefix string, src Value, dst reflect.Value, weaklyTyped bool) error {
	var (
		b   Blob
		err error
	)
	if weaklyTyped {
		b, err = ToBlob(src)
	} else {
		b, err = AsBlob(src)
	}
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	dst.Set(reflect.ValueOf(b))
	return nil
}

func (d *Decoder) decodeStruct(prefix string, src Value, dst reflect.Value) error {
	if dst.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
		if prefix == "" {
			// time.Time or Timestamp is passed directly to Decode. They have
			// to be a field of a struct.
			return errors.New("timestamp cannot be decoded directly")
		}
		return d.decodeTimestamp(prefix, src, dst)
	}

	m, err := AsMap(src)
	if err != nil {
		return fmt.Errorf("%v: struct can only be decoded from a map", prefix)
	}

	errPrefix := prefix + "."
	if prefix == "" {
		errPrefix = ""
	}

	var unused map[string]struct{}
	if d.config.ErrorUnused {
		unused = make(map[string]struct{}, len(m))
		for k := range m {
			unused[k] = struct{}{}
		}
	}

	// Accumulate all error informations to help users debug BQL.
	errs := d.iterateField(errPrefix, m, unused, dst)
	if d.config.ErrorUnused && len(unused) > 0 {
		keys := make([]string, len(unused))
		i := 0
		for k := range unused {
			keys[i] = errPrefix + k
			i++
		}
		if d.config.Metadata != nil {
			d.config.Metadata.Unused = append(d.config.Metadata.Unused, keys...)
		}

		errs = multierror.Append(errs, fmt.Errorf("%v: unused keys: %v", prefix, strings.Join(keys, ", ")))
	}
	if errs == nil {
		// To avoid nil != nil problem due to type mismatch, this function has
		// to return nil explicitly. Don't do "return errs" even if errs == nil.
		return nil
	}

	// TODO: set formatter
	return errs
}

func (d *Decoder) iterateField(prefix string, m Map, unused map[string]struct{}, dst reflect.Value) *multierror.Error {
	// FIXME: iterateField decodes the same value multiple times if a struct and
	// its embedded struct have the same field names.

	// Accumulate all error informations to help users debug BQL.
	var errs *multierror.Error

	if dst.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
		errs = multierror.Append(errs, fmt.Errorf("%v: time.Time and data.Timestamp cannot be embedded", prefix))
		return errs
	}

	t := dst.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		if f.Anonymous { // process embedded field
			if f.Type.Kind() == reflect.Struct {
				if err := d.iterateField(prefix, m, unused, dst.Field(i)); err != nil {
					errs = multierror.Append(errs, err)
				}
				continue

			} else if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
				v := reflect.New(f.Type.Elem())
				if err := d.iterateField(prefix, m, unused, reflect.Indirect(v)); err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
				dst.Field(i).Set(v)
				continue
			}
			errs = multierror.Append(errs, fmt.Errorf("%v: unsupported embedded field: %v", prefix, f.Name))
		}

		opts := strings.Split(f.Tag.Get(d.config.TagName), ",")

		// parse options
		var (
			required    bool
			weaklyTyped bool
		)
		for ti, opt := range opts {
			if ti == 0 { // skip name
				continue
			}
			switch opt {
			case "required":
				required = true
			case "weaklytyped":
				weaklyTyped = true
			default:
				errs = multierror.Append(errs, fmt.Errorf("%v%v: an undefined option: %v", prefix, f.Name, opt))
			}
		}

		name := strings.TrimSpace(opts[0])
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		if d.config.ErrorUnused {
			delete(unused, name)
		}
		src, ok := m[name]
		if !ok {
			if required {
				errs = multierror.Append(errs, fmt.Errorf("%v%v: required but missing", prefix, name))
			}
			continue
		}
		if d.config.Metadata != nil {
			d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
		}

		if err := d.decode(prefix+name, src, dst.Field(i), weaklyTyped); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (d *Decoder) decodeTimestamp(prefix string, src Value, dst reflect.Value) error {
	// src is provided as a part of BQL and there's no way in BQL to construct
	// Timestamp directly. So, conversions to time.Time and Timestamp is always
	// considered weaklytyped.

	t, err := ToTimestamp(src)
	if err != nil {
		return fmt.Errorf("%v: %v", prefix, err)
	}
	switch dst.Interface().(type) {
	case time.Time:
		dst.Set(reflect.ValueOf(t))
	case Timestamp:
		dst.Set(reflect.ValueOf(Timestamp(t)))
	default:
		return fmt.Errorf("%v: only time.Time and data.Timestamp can be used for decoding a timestamp", prefix)
	}
	return nil
}

func toSnakeCase(name string) string {
	words := camelcase.Split(name)
	buf := bytes.NewBuffer(nil)
	for i, w := range words {
		if i > 0 {
			buf.WriteString("_")
		}
		buf.WriteString(strings.ToLower(w))
	}
	return string(buf.Bytes())
}

// Decode decodes a Map into a struct. The argument must be a pointer to a
// struct.
func Decode(m Map, v interface{}) error {
	return NewDecoder(&DecoderConfig{
		ErrorUnused: true,
	}).Decode(m, v)
}
