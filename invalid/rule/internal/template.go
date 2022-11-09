package internal

import (
	_ "embed"
	"strings"
	"text/template"
)

var (
	//go:embed template/default/single.tmpl
	defalutSingleStr string
	//go:embed template/default/ptr_single.tmpl
	defalutPtrSingleStr  string
	defaultSingleTmpl    *template.Template
	defaultPtrSingleTmpl *template.Template

	//go:embed template/duck/single.tmpl
	duckSingleStr string
	//go:embed template/duck/slice.tmpl
	duckSliceStr string
	//go:embed template/duck/ptr_slice.tmpl
	duckPtrSliceStr  string
	duckSingleTmpl   *template.Template
	duckSliceTmpl    *template.Template
	duckPtrSliceTmpl *template.Template

	//go:embed template/enum/single.tmpl
	enumSingleStr string
	//go:embed template/enum/ptr_single.tmpl
	enumPtrSingleStr string
	//go:embed template/enum/slice.tmpl
	enumSliceStr string
	//go:embed template/enum/ptr_slice.tmpl
	enumPtrSliceStr   string
	enumSingleTmpl    *template.Template
	enumPtrSingleTmpl *template.Template
	enumSliceTmpl     *template.Template
	enumPtrSliceTmpl  *template.Template

	//go:embed template/equal/single.tmpl
	equalSingleStr string
	//go:embed template/equal/ptr_single.tmpl
	equalPtrSingleStr string
	//go:embed template/equal/slice.tmpl
	equalSliceStr string
	//go:embed template/equal/ptr_slice.tmpl
	equalPtrSliceStr   string
	equalSingleTmpl    *template.Template
	equalPtrSingleTmpl *template.Template
	equalSliceTmpl     *template.Template
	equalPtrSliceTmpl  *template.Template

	//go:embed template/length/single.tmpl
	lengthSingleStr string
	//go:embed template/length/ptr_single.tmpl
	lengthPtrSingleStr string
	//go:embed template/length/slice.tmpl
	lengthSliceStr string
	//go:embed template/length/ptr_slice.tmpl
	lengthPtrSliceStr   string
	lengthSingleTmpl    *template.Template
	lengthPtrSingleTmpl *template.Template
	lengthSliceTmpl     *template.Template
	lengthPtrSliceTmpl  *template.Template

	//go:embed template/not/single.tmpl
	notSingleStr string
	//go:embed template/not/ptr_single.tmpl
	notPtrSingleStr string
	//go:embed template/not/slice.tmpl
	notSliceStr string
	//go:embed template/not/ptr_slice.tmpl
	notPtrSliceStr   string
	notSingleTmpl    *template.Template
	notPtrSingleTmpl *template.Template
	notSliceTmpl     *template.Template
	notPtrSliceTmpl  *template.Template

	//go:embed template/range/single.tmpl
	rangeSingleStr string
	//go:embed template/range/ptr_single.tmpl
	rangePtrSingleStr string
	//go:embed template/range/slice.tmpl
	rangeSliceStr string
	//go:embed template/range/ptr_slice.tmpl
	rangePtrSliceStr   string
	rangeSingleTmpl    *template.Template
	rangePtrSingleTmpl *template.Template
	rangeSliceTmpl     *template.Template
	rangePtrSliceTmpl  *template.Template

	//go:embed template/range_time/single.tmpl
	rangeTimeSingleStr string
	//go:embed template/range_time/ptr_single.tmpl
	rangeTimePtrSingleStr string
	//go:embed template/range_time/slice.tmpl
	rangeTimeSliceStr string
	//go:embed template/range_time/ptr_slice.tmpl
	rangeTimePtrSliceStr   string
	rangeTimeSingleTmpl    *template.Template
	rangeTimePtrSingleTmpl *template.Template
	rangeTimeSliceTmpl     *template.Template
	rangeTimePtrSliceTmpl  *template.Template

	//go:embed template/regexp/single.tmpl
	regexpSingleStr string
	//go:embed template/regexp/ptr_single.tmpl
	regexpPtrSingleStr string
	//go:embed template/regexp/slice.tmpl
	regexpSliceStr string
	//go:embed template/regexp/ptr_slice.tmpl
	regexpPtrSliceStr   string
	regexpSingleTmpl    *template.Template
	regexpPtrSingleTmpl *template.Template
	regexpSliceTmpl     *template.Template
	regexpPtrSliceTmpl  *template.Template

	//go:embed template/unique/single.tmpl
	uniqueSingleStr string
	//go:embed template/unique/ptr_single.tmpl
	uniquePtrSingleStr  string
	uniqueSingleTmpl    *template.Template
	uniquePtrSingleTmpl *template.Template
)

func init() {
	var err error

	prefix :=
		`{{- if len_gt_0 .tags }}
        OUT{{ .index }}:
		for _, paramTag := range tags {
            for _, condTag := range {{ .tags }} {
                if paramTag == condTag {
        {{- end }}`
	suffix :=
		`{{- if len_gt_0 .tags }}
					break OUT{{ .index }}
				}
			}
		}
		{{- end }}`

	// default
	{
		defaultSingleTmpl, err = template.
			New("invalid_defalut").
			Funcs(funcMap).
			Parse(prefix + defalutSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		defaultPtrSingleTmpl, err = template.
			New("invalid_defalut_ptr").
			Funcs(funcMap).
			Parse(prefix + defalutPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// duck
	{
		duckSingleTmpl, err = template.
			New("invalid_duck").
			Funcs(funcMap).
			Parse(prefix + duckSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		duckSliceTmpl, err = template.
			New("invalid_duck_slice").
			Funcs(funcMap).
			Parse(prefix + duckSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		duckPtrSliceTmpl, err = template.
			New("invalid_duck_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + duckPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// enum
	{
		enumSingleTmpl, err = template.
			New("invalid_enum").
			Funcs(funcMap).
			Parse(prefix + enumSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		enumPtrSingleTmpl, err = template.
			New("invalid_enum_ptr").
			Funcs(funcMap).
			Parse(prefix + enumPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		enumSliceTmpl, err = template.
			New("invalid_enum_slice").
			Funcs(funcMap).
			Parse(prefix + enumSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		enumPtrSliceTmpl, err = template.
			New("invalid_enum_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + enumPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// equal
	{
		equalSingleTmpl, err = template.
			New("invalid_equal").
			Funcs(funcMap).
			Parse(prefix + equalSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		equalPtrSingleTmpl, err = template.
			New("invalid_equal_ptr").
			Funcs(funcMap).
			Parse(prefix + equalPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		equalSliceTmpl, err = template.
			New("invalid_equal_slice").
			Funcs(funcMap).
			Parse(prefix + equalSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		equalPtrSliceTmpl, err = template.
			New("invalid_equal_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + equalPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// length / length_utf8
	{
		lengthSingleTmpl, err = template.
			New("invalid_length").
			Funcs(funcMap).
			Parse(prefix + lengthSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		lengthPtrSingleTmpl, err = template.
			New("invalid_length_ptr").
			Funcs(funcMap).
			Parse(prefix + lengthPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		lengthSliceTmpl, err = template.
			New("invalid_length_slice").
			Funcs(funcMap).
			Parse(prefix + lengthSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		lengthPtrSliceTmpl, err = template.
			New("invalid_length_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + lengthPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// not
	{
		notSingleTmpl, err = template.
			New("invalid_not").
			Funcs(funcMap).
			Parse(prefix + notSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		notPtrSingleTmpl, err = template.
			New("invalid_not_ptr").
			Funcs(funcMap).
			Parse(prefix + notPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		notSliceTmpl, err = template.
			New("invalid_not_slice").
			Funcs(funcMap).
			Parse(prefix + notSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		notPtrSliceTmpl, err = template.
			New("invalid_not_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + notPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// range
	{
		rangeSingleTmpl, err = template.
			New("invalid_range").
			Funcs(funcMap).
			Parse(prefix + rangeSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		rangePtrSingleTmpl, err = template.
			New("invalid_range_ptr").
			Funcs(funcMap).
			Parse(prefix + rangePtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		rangeSliceTmpl, err = template.
			New("invalid_range_slice").
			Funcs(funcMap).
			Parse(prefix + rangeSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		rangePtrSliceTmpl, err = template.
			New("invalid_range_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + rangePtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// range time
	{
		rangeTimeSingleTmpl, err = template.
			New("invalid_range_time").
			Funcs(funcMap).
			Parse(prefix + rangeTimeSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		rangeTimePtrSingleTmpl, err = template.
			New("invalid_range_time_ptr").
			Funcs(funcMap).
			Parse(prefix + rangeTimePtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		rangeTimeSliceTmpl, err = template.
			New("invalid_range_time_slice").
			Funcs(funcMap).
			Parse(prefix + rangeTimeSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		rangeTimePtrSliceTmpl, err = template.
			New("invalid_range_time_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + rangeTimePtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// regexp
	{
		regexpSingleTmpl, err = template.
			New("invalid_regexp").
			Funcs(funcMap).
			Parse(prefix + regexpSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		regexpPtrSingleTmpl, err = template.
			New("invalid_regexp_ptr").
			Funcs(funcMap).
			Parse(prefix + regexpPtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		regexpSliceTmpl, err = template.
			New("invalid_regexp_slice").
			Funcs(funcMap).
			Parse(prefix + regexpSliceStr + suffix)
		if err != nil {
			panic(err)
		}

		regexpPtrSliceTmpl, err = template.
			New("invalid_regexp_slice_ptr").
			Funcs(funcMap).
			Parse(prefix + regexpPtrSliceStr + suffix)
		if err != nil {
			panic(err)
		}
	}

	// unique
	{
		uniqueSingleTmpl, err = template.
			New("invalid_unique").
			Funcs(funcMap).
			Parse(prefix + uniqueSingleStr + suffix)
		if err != nil {
			panic(err)
		}

		uniquePtrSingleTmpl, err = template.
			New("invalid_unique_ptr").
			Funcs(funcMap).
			Parse(prefix + uniquePtrSingleStr + suffix)
		if err != nil {
			panic(err)
		}
	}
}

func GetTmpl(iden, typ string) *template.Template {
	if strings.HasPrefix(typ, "*[]") {
		typ = "ptrslice"
	} else if strings.HasPrefix(typ, "[]") {
		typ = "slice"
	} else if strings.HasPrefix(typ, "*") {
		typ = "ptrsingle"
	} else {
		typ = "single"
	}

	val := map[string]map[string]*template.Template{
		Default: {
			"single":    defaultSingleTmpl,
			"ptrsingle": defaultPtrSingleTmpl,
			"slice":     defaultSingleTmpl,
			"ptrslice":  defaultPtrSingleTmpl,
		},
		Duck: {
			"single":    duckSingleTmpl,
			"ptrsingle": duckSingleTmpl,
			"slice":     duckSliceTmpl,
			"ptrslice":  duckPtrSliceTmpl,
		},
		Enum: {
			"single":    enumSingleTmpl,
			"slice":     enumSliceTmpl,
			"ptrsingle": enumPtrSingleTmpl,
			"ptrslice":  enumPtrSliceTmpl,
		},
		Equal: {
			"single":    equalSingleTmpl,
			"slice":     equalSliceTmpl,
			"ptrsingle": equalPtrSingleTmpl,
			"ptrslice":  equalPtrSliceTmpl,
		},
		Length: {
			"single":    lengthSingleTmpl,
			"slice":     lengthSliceTmpl,
			"ptrsingle": lengthPtrSingleTmpl,
			"ptrslice":  lengthPtrSliceTmpl,
		},
		LengthUtf8: {
			"single":    lengthSingleTmpl,
			"slice":     lengthSliceTmpl,
			"ptrsingle": lengthPtrSingleTmpl,
			"ptrslice":  lengthPtrSliceTmpl,
		},
		Not: {
			"single":    notSingleTmpl,
			"slice":     notSliceTmpl,
			"ptrsingle": notPtrSingleTmpl,
			"ptrslice":  notPtrSliceTmpl,
		},
		Range: {
			"single":    rangeSingleTmpl,
			"slice":     rangeSliceTmpl,
			"ptrsingle": rangePtrSingleTmpl,
			"ptrslice":  rangePtrSliceTmpl,
		},
		RangeTime: {
			"single":    rangeTimeSingleTmpl,
			"slice":     rangeTimeSliceTmpl,
			"ptrsingle": rangeTimePtrSingleTmpl,
			"ptrslice":  rangeTimePtrSliceTmpl,
		},
		Regexp: {
			"single":    regexpSingleTmpl,
			"slice":     regexpSliceTmpl,
			"ptrsingle": regexpPtrSingleTmpl,
			"ptrslice":  regexpPtrSliceTmpl,
		},
		Unique: {
			"single":    uniqueSingleTmpl,
			"ptrsingle": uniquePtrSingleTmpl,
			"slice":     uniqueSingleTmpl,
			"ptrslice":  uniquePtrSingleTmpl,
		},
	}

	return val[iden][typ]
}
