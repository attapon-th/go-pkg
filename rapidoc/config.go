package rapidoc

type RenderStyle string
type SchemaStyle string

const (
	RenderStyle_Read  RenderStyle = "read"
	RenderStyle_View  RenderStyle = "view"
	RenderStyle_Focus RenderStyle = "focused"

	SchemaStyle_Tree  SchemaStyle = "tree"
	SchemaStyle_Table SchemaStyle = "table"
)

type RapiDocConfig struct {
	Title       string      `json:"tiltle,omitempty"`
	SpecURL     string      `json:"spec_url,omitempty"`
	HeaderText  string      `json:"header_text,omitempty"`
	LogoURL     string      `json:"logo_url,omitempty"`
	RenderStyle RenderStyle `json:"render_style,omitempty"`
	SchemaStyle SchemaStyle `json:"schema_style,omitempty"`
}

func GetDefaultRapiDocConfig() RapiDocConfig {
	return RapiDocConfig{
		Title:       "API Documentation",
		SpecURL:     "./swagger.json",
		HeaderText:  "API Documentation",
		LogoURL:     "https://mrin9.github.io/RapiDoc/images/logo.png",
		RenderStyle: RenderStyle_Read,
		SchemaStyle: SchemaStyle_Tree,
	}
}

func HtmlTemplateRapiDoc() string {
	return `<!doctype html>
	<html>	
	<head>
		<title>{{$.Title}}</title>
		<meta charset="utf-8">
		<link href="https://fonts.googleapis.com/css2?family=Sarabun&display=swap" rel="stylesheet">
		<link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;600&family=Roboto+Mono&display=swap" rel="stylesheet">
		<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
	</head>
	
	<body>
		<rapi-doc spec-url="{{$.SpecURL}}" heading-text="{{$.HeaderText}}" regular-font="Sarabun" mono-font="'Roboto Mono'" render-style="{{$.RenderStyle}}" nav-text-color="#aaa" nav-hover-text-color="#fff" nav-accent-color="#FF0000" primary-color="#3B5998" show-header="false" show-info="true" allow-authentication="true" theme="dark" allow-try="true" allow-search="false" allow-spec-url-load="false" allow-spec-file-load="false" header-color="#3B5998" schema-style="{{$.SchemaStyle}}">
		<div slot="nav-logo" style="display: flex; align-items: center; justify-content: center;">
			<img src="{{$.LogoURL}}" style="width:50px; margin-right: 20px">
		</div>
		</rapi-doc>
	</body>
	</html>`
}
