package core

import (
	"encoding/xml"
)

/*
	type server struct {
    XMLName    xml.Name `xml:"server"`
    Port string `xml:"port"`
    Package  string `xml:"package"`
}
	GoS Xml Strucs
	Links are used to import methods
*/
/*
	Struct for Go method handling
*/
type gosArch struct {
	interfacedMethods []string
	methodlimits      []string
	keeplocal         []string
	webmethods        []string
	structs           []string
	objects           []string
}

type gos struct {
	//
	XMLName xml.Name `xml:"gos"`
	// Port your webapplication will liston on
	Port string `xml:"port"`
	// Valid hostname to be used with secure cookies.
	Domain string `xml:"domain"`
	// OpenFaaS gateway to deploy application to.
	Gate string `xml:"gateway,attr"`
	//
	Debug string `xml:"trace"`
	// name of written Gofile
	Output string `xml:"output"`
	// URI to 500 code page.
	ErrorPage string `xml:"error"`
	// URI to 400 code page.
	NPage string `xml:"not_found"`
	// GoS deployment type. : faas, webapp, package
	Type string `xml:"deploy"`
	// Block of code to be ran on application start.
	// Only applies to files with deploy type : webapp
	Main          string            `xml:"main"`
	Variables     []GlobalVariables `xml:"var"`
	WriteOut      bool
	Export        string   `xml:"export"`
	Key           string   `xml:"key"`
	Session       string   `xml:"session"`
	Template_path string   `xml:"templatePath"`
	Web_root      string   `xml:"webroot"`
	Package       string   `xml:"package"`
	Web           string   `xml:"web"`
	Tmpl          string   `xml:"tmpl"`
	RootImports   []Import `xml:"import"`
	Init_Func     string   `xml:"init"`
	Header        Header   `xml:"header"`
	// Web server functions. As well as any <func>
	// tag.
	Methods     Methods   `xml:"methods"`
	PostCommand []string  `xml:"sh"`
	Timers      Timers    `xml:"timers"`
	Templates   Templates `xml:"templates"`
	// Web service endpoints.
	Endpoints        Endpoints `xml:"endpoints"`
	FolderRoot, Name string
	// Set Prod to true to build your application in production mode.
	Prod bool
}

type Pgos struct {
	gos
}

type GlobalVariables struct {
	XMLName xml.Name `xml:"var"`
	Name    string   `xml:",innerxml"`
	Type    string   `xml:"type,attr"`
}

type Import struct {
	XMLName xml.Name `xml:"import"`
	Src     string   `xml:"src,attr"`
}

type Header struct {
	XMLName xml.Name `xml:"header"`
	Structs []Struct `xml:"struct"`
	Objects []Object `xml:"object"`
}

type Methods struct {
	XMLName xml.Name `xml:"methods"`
	Methods []Method `xml:"method"`
}

type Timers struct {
	XMLName xml.Name `xml:"timers"`
	Timers  []Timer  `xml:"timer"`
}

type Templates struct {
	XMLName   xml.Name   `xml:"templates"`
	Templates []Template `xml:"template"`
}

// Webservice endpoints.
type Endpoints struct {
	XMLName xml.Name `xml:"endpoints"`
	// Array of Web service URIs.
	Endpoints []Endpoint `xml:"end"`
}

/*
	Nested values within GoS root file
*/

type VGos struct {
	XMLName xml.Name `xml:"gos"`
	Objects []Object `xml:"object"`
	Structs []Struct `xml:"struct"`
	Methods []Method `xml:"method"`
}

type Struct struct {
	XMLName    xml.Name `xml:"struct"`
	Name       string   `xml:"name,attr"`
	Attributes string   `xml:",innerxml"`
}

type Object struct {
	XMLName xml.Name `xml:"object"`
	Name    string   `xml:"name,attr"`
	Templ   string   `xml:"struct,attr"`
	Methods string   `xml:",innerxml"`
}

type Method struct {
	XMLName    xml.Name    `xml:"method"`
	Method     string      `xml:",innerxml"`
	Comment    xml.Comment `xml:",comment"`
	Name       string      `xml:"name,attr"`
	Variables  string      `xml:"var,attr"`
	Limit      string      `xml:"limit,attr"`
	Object     string      `xml:"object,attr"`
	Autoface   string      `xml:"autoface,attr"`
	Keeplocal  string      `xml:"keep-local,attr"`
	Testi      string      `xml:"testi,attr"`
	Testo      string      `xml:"testo,attr"`
	Man        string      `xml:"m,attr"`
	Returntype string      `xml:"return,attr"`
}

type Timer struct {
	XMLName  xml.Name `xml:"timer"`
	Method   string   `xml:",innerxml"`
	Interval string   `xml:"interval,attr"`
	Name     string   `xml:"name,attr"`
	Unit     string   `xml:"unit,attr"`
}

type Template struct {
	XMLName xml.Name `xml:"template"`
	// User-entered template identifier.
	Name string `xml:"name,attr"`
	// Template file relative to tmpl folder,
	// without .tmpl suffix as well.
	TemplateFile string `xml:"tmpl,attr"`
	//
	Bundle string `xml:"bundle,attr"`
	// Interface to use with template.
	Struct    string `xml:"struct,attr"`
	ForcePath bool
	Comment   xml.Comment `xml:",comment"`
}

type Endpoint struct {
	XMLName xml.Name `xml:"end"`
	// User-entered endpoint URI.
	Path string `xml:"path,attr"`
	// Code block to run on URI load.
	Method string `xml:",innerxml"`
	// User-entered endpoint request verb type.
	// If a URI matches a request but the verb does not
	// the endpoint will not load.
	Type  string `xml:"type,attr"`
	Testi string `xml:"testi,attr"`
	Testo string `xml:"testo,attr"`
	Id    string `xml:"id,attr"`
}
