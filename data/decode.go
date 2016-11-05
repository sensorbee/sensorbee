package data

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/camelcase"
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
	return d.decodeStruct(m, s)
}

func (d *Decoder) decode(src Value, dst reflect.Value, weaklyTyped bool) error {
	switch dst.Kind() {
	case reflect.Bool:
		return d.decodeBool(src, dst, weaklyTyped)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// TODO: support detecting overflows
		return d.decodeInt(src, dst, weaklyTyped)

	case reflect.Float32, reflect.Float64:
		return d.decodeFloat(src, dst, weaklyTyped)

	case reflect.String:
		return d.decodeString(src, dst, weaklyTyped)

	case reflect.Interface: // Only Value is supported
		if dst.Type() != reflect.TypeOf(func(Value) {}).In(0) {
			return errors.New("interface{} other than data.Value is not supported")
		}
		return d.decodeValue(src, dst)

	case reflect.Map:
		return d.decodeMap(src, dst, weaklyTyped)

	case reflect.Slice:
		return d.decodeSlice(src, dst, weaklyTyped)

	case reflect.Struct:
		return d.decodeStruct(src, dst)

	case reflect.Ptr:
		// To decode a value to dst, dst must be addressable. However,
		// reflect.ValueOf or reflect.Zero may return non-addressable values
		// especially when value is of primitive types. The following
		// New-Indirect idiom works for such cases. First reflect.New creates
		// a pointer that points to a non-nil addressable value. Then,
		// reflect.Indirect returns an element pointed by the pointer.
		v := reflect.New(dst.Type().Elem())
		if err := d.decode(src, reflect.Indirect(v), weaklyTyped); err != nil {
			return err
		}
		dst.Set(v)
		return nil
	}
	return fmt.Errorf("decoder doesn't support the type: %v", dst.Kind())
}

func (d *Decoder) decodeBool(src Value, dst reflect.Value, weaklyTyped bool) error {
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
		return err
	}
	dst.SetBool(b)
	return nil
}

func (d *Decoder) decodeInt(src Value, dst reflect.Value, weaklyTyped bool) error {
	if _, ok := dst.Interface().(time.Duration); ok {
		return d.decodeDuration(src, dst)
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
		return err
	}
	dst.SetInt(i)
	return nil
}

func (d *Decoder) decodeDuration(src Value, dst reflect.Value) error {
	// As described in decodeTimestamp, decodeDuration also assumes that the
	// value is always weaklytyped.
	dur, err := ToDuration(src)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(dur))
	return nil
}

func (d *Decoder) decodeFloat(src Value, dst reflect.Value, weaklyTyped bool) error {
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
		return err
	}
	dst.SetFloat(f)
	return nil
}

func (d *Decoder) decodeString(src Value, dst reflect.Value, weaklyTyped bool) error {
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
		return err
	}
	dst.SetString(s)
	return nil
}

func (d *Decoder) decodeValue(src Value, dst reflect.Value) error {
	dst.Set(reflect.ValueOf(src))
	return nil
}

func (d *Decoder) decodeMap(src Value, dst reflect.Value, weaklyTyped bool) error {
	if src.Type() != TypeMap {
		return fmt.Errorf("cannot decode to a map: %v", src.Type())
	}
	m, _ := AsMap(src)

	t := dst.Type()
	if k := t.Key().Kind(); k != reflect.String {
		return fmt.Errorf("key must be string: %v", k)
	}
	valueType := t.Elem()

	res := reflect.MakeMap(t)
	for k, e := range m {
		v := reflect.Indirect(reflect.New(valueType))
		if err := d.decode(e, v, weaklyTyped); err != nil {
			// TODO: this should probably be multierror, too.
			return err
		}
		res.SetMapIndex(reflect.ValueOf(k), v)
	}
	dst.Set(res)
	return nil
}

func (d *Decoder) decodeSlice(src Value, dst reflect.Value, weaklyTyped bool) error {
	if dst.Type().Elem().Kind() == reflect.Uint8 {
		return d.decodeBlob(src, dst, weaklyTyped)
	}

	if src.Type() != TypeArray {
		return fmt.Errorf("cannot decode to an array: %v", src.Type())
	}
	a, _ := AsArray(src)

	res := reflect.MakeSlice(dst.Type(), len(a), len(a))
	for i, e := range a {
		v := res.Index(i)
		if err := d.decode(e, v, weaklyTyped); err != nil {
			// TODO: this should probably be multierror, too.
			return err
		}
	}
	dst.Set(res)
	return nil
}

func (d *Decoder) decodeBlob(src Value, dst reflect.Value, weaklyTyped bool) error {
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
		return err
	}
	dst.Set(reflect.ValueOf(b))
	return nil
}

func (d *Decoder) decodeStruct(src Value, dst reflect.Value) error {
	if dst.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
		return d.decodeTimestamp(src, dst)
	}

	m, err := AsMap(src)
	if err != nil {
		return errors.New("struct can only be decoded from a map")
	}

	// Aggregates all error informations to help users debug BQL.
	// TODO: replace this with github.com/hashicorp/go-multierror
	var errs []error

	t := dst.Type()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
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
				errs = append(errs, fmt.Errorf("%v field has an undefined option: %v", f.Name, opt))
			}
		}

		name := strings.TrimSpace(opts[0])
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		src, ok := m[name]
		if !ok {
			if required {
				errs = append(errs, fmt.Errorf("%v is required but missing", name))
			}
			continue
		}

		if err := d.decode(src, dst.Field(i), weaklyTyped); err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		// TODO: flatten errors
		return errs[0]
	}
	return nil
}

func (d *Decoder) decodeTimestamp(src Value, dst reflect.Value) error {
	// src is provided as a part of BQL and there's no way in BQL to construct
	// Timestamp directly. So, conversions to time.Time and Timestamp is always
	// considered weaklytyped.

	t, err := ToTimestamp(src)
	if err != nil {
		return err
	}
	switch dst.Interface().(type) {
	case time.Time:
		dst.Set(reflect.ValueOf(t))
	case Timestamp:
		dst.Set(reflect.ValueOf(Timestamp(t)))
	default:
		return errors.New("only time.Time and data.Timestamp can be used for decoding a timestamp")
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
