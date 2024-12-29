// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package styles

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Button() templ.CSSClass {
	templ_7745c5c3_CSSBuilder := templruntime.GetBuilder()
	templ_7745c5c3_CSSBuilder.WriteString(`background-color:var(--color-primary);`)
	templ_7745c5c3_CSSBuilder.WriteString(`border:none;`)
	templ_7745c5c3_CSSBuilder.WriteString(`color:white;`)
	templ_7745c5c3_CSSBuilder.WriteString(`padding:12px 24px;`)
	templ_7745c5c3_CSSBuilder.WriteString(`text-align:center;`)
	templ_7745c5c3_CSSBuilder.WriteString(`text-decoration:none;`)
	templ_7745c5c3_CSSBuilder.WriteString(`display:inline-block;`)
	templ_7745c5c3_CSSBuilder.WriteString(`font-size:16px;`)
	templ_7745c5c3_CSSBuilder.WriteString(`margin:4px 2px;`)
	templ_7745c5c3_CSSBuilder.WriteString(`cursor:pointer;`)
	templ_7745c5c3_CSSBuilder.WriteString(`border-radius:8px;`)
	templ_7745c5c3_CSSBuilder.WriteString(`transition:background-color 0.3s ease, transform 0.2s ease;`)
	templ_7745c5c3_CSSID := templ.CSSID(`Button`, templ_7745c5c3_CSSBuilder.String())
	return templ.ComponentCSSClass{
		ID:    templ_7745c5c3_CSSID,
		Class: templ.SafeCSS(`.` + templ_7745c5c3_CSSID + `{` + templ_7745c5c3_CSSBuilder.String() + `}`),
	}
}

var _ = templruntime.GeneratedTemplate