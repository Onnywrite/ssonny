package postgres_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/Onnywrite/ssonny/internal/domain/models"
	"github.com/google/uuid"
)

type Post struct {
	Id      uint64
	Title   *string `qt:"VARCHAR(64)" qc:"UNIQUE" qd:"NULL"`
	Content string
	Author  User `db:"user_fk"`
}

type User struct {
	Id           uuid.UUID
	Nickname     *string `qt:"VARCHAR(32)" qc:"UNIQUE"`
	Email        string  `qt:"VARCHAR(345)" qc:"UNIQUE"`
	IsVerified   bool    `db:"verified"`
	Gender       *string `qt:"VARCHAR(16)"`
	PasswordHash *string `qt:"CHAR(60)"`

	Birthday  *time.Time `qt:"DATE"`
	CreatedAt time.Time  `qd:"NOW()"`
	UpdatedAt time.Time  `qd:"NOW()"`
	DeletedAt *time.Time
}

func TestMe(t *testing.T) {
	fmt.Println(CreateTableFor[models.App]())
	t.FailNow()
}

func CreateTableFor[T any]() string {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		panic("not struct")
	}
	b := strings.Builder{}
	b.WriteString("CREATE TABLE")
	b.WriteString(" " + toSnake(t.Name()) + "s (\n")
	for i := range t.NumField() {
		b.WriteByte('\t')
		field := t.Field(i)

		// name
		fieldName := getName(field)
		b.WriteString(" " + fieldName)

		nullable := false
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			nullable = true
			fieldType = fieldType.Elem()
		}

		// type
		if strings.HasSuffix(fieldName, "id") {
			if nullable {
				panic("nullable primary key")
			}
			isUuid := false
			switch fieldType.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
				b.WriteString(" SMALLSERIAL")
			case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
				b.WriteString(" SERIAL")
			case reflect.Int64, reflect.Uint64:
				b.WriteString(" BIGSERIAL")
			default:
				if fieldType.ConvertibleTo(reflect.TypeFor[uuid.UUID]()) {
					b.WriteString(" UUID")
					isUuid = true
				} else {
					panic("invalid primary key type")
				}
			}
			b.WriteString(" PRIMARY KEY")
			if constraints, ok := field.Tag.Lookup("qc"); ok {
				b.WriteString(" " + constraints)
			}
			if isUuid {
				b.WriteString(" DEFAULT gen_random_uuid()")
			}
		} else {
			if qtTag, ok := field.Tag.Lookup("qt"); ok {
				b.WriteString(" " + qtTag)
			} else if strings.Contains(strings.ToLower(fieldType.Name()), "time") {
				b.WriteString(" TIMESTAMP")
			} else if fieldType.Kind() == reflect.Struct {
				if !strings.HasSuffix(fieldName, "_fk") {
					b.WriteString("_fk")
				}
				b.WriteString(" " + referenceStruct(fieldType))
			} else {
				b.WriteString(" " + switchType(fieldType))
			}

			if !nullable {
				b.WriteString(" NOT NULL")
			}
			if constraints, ok := field.Tag.Lookup("qc"); ok {
				b.WriteString(" " + constraints)
			}
			if defaultValue, ok := field.Tag.Lookup("qd"); ok {
				b.WriteString(" DEFAULT " + defaultValue)
			}
		}

		if i != t.NumField()-1 {
			b.WriteString(",\n")
		}
	}
	b.WriteString("\n);")
	return b.String()
}

func referenceStruct(t reflect.Type) string {
	if t.Kind() != reflect.Struct {
		panic("not struct")
	}

	for i := range t.NumField() {
		field := t.Field(i)

		fieldName := getName(field)

		if strings.HasSuffix(fieldName, "id") {
			return switchType(field.Type)
		}
	}
	panic("no id found in struct " + t.Name())
}

func getName(f reflect.StructField) (fieldName string) {
	if dbTag, ok := f.Tag.Lookup("db"); ok {
		fieldName = toSnake(dbTag)
	} else {
		fieldName = toSnake(f.Name)
	}
	return
}

func switchType(t reflect.Type) string {
	if strings.Contains(strings.ToLower(t.Name()), "uuid") {
		return "UUID"
	}
	switch t.Kind() {
	case reflect.String:
		return "TEXT"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.Int8, reflect.Int16:
		return "SMALLINT"
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint8, reflect.Uint16:
		return "SMALLINT"
	case reflect.Uint, reflect.Uint32:
		return "INT"
	case reflect.Uint64:
		return "BIGINT"
	case reflect.Array:
		return switchType(t.Elem()) + "[" + strconv.FormatInt(int64(t.Len()), 10) + "]"
	case reflect.Slice:
		return switchType(t.Elem()) + "[]"
	}
	panic("invalid type " + t.Name())
}

func toSnake(s string) string {
	b := strings.Builder{}
	b.Grow(len(s))
	bytes := []byte(s)
	for i, e := range bytes {
		if unicode.IsUpper(rune(e)) && i > 0 {
			if i != 0 {
				b.WriteRune('_')
			}
			b.WriteRune(unicode.ToLower(rune(e)))
		} else if unicode.IsUpper(rune(e)) && i == 0 {
			b.WriteRune(unicode.ToLower(rune(e)))
		} else {
			b.WriteByte(e)
		}
	}
	return b.String()
}
