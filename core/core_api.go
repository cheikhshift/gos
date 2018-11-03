package core

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/fatih/color"
	"github.com/tj/go-spin"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"
	//"go/types"
)

const (
	StdLen  = 16
	UUIDLen = 20
)

var primitives = []string{"Bool",
	"Int",
	"Int8",
	"Int16",
	"Int32",
	"Int64",
	"Uint",
	"Uint8",
	"Uint16",
	"Uint32",
	"Uint64",
	"Uintptr",
	"Float32",
	"Float64",
	"Complex64",
	"Complex128",
	"String",
	"UnsafePointer",
	"UntypedBool",
	"UntypedInt",
	"UntypedRune",
	"UntypedFloat",
	"UntypedComplex",
	"UntypedString",
	"UntypedNil",
	"Byte",
	"Rune"}

var GOHOME = os.ExpandEnv("$GOPATH") + "/src/"
var available_methods []string
var int_methods []string
var api_methods []string
var int_mappings []string
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var StdNums = []byte("OYZ0123456789")
var DAMP = "&&"
var AMP = "&"

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//Local DEBUG tools of Gos

func IsInImports(imports []*ast.ImportSpec, unfoundvar string) bool {

	for _, v := range imports {
		ab := strings.Split(strings.Replace(v.Path.Value, "\"", "", 2), "/")
		if ab[len(ab)-1] == unfoundvar {
			return true
		}

	}

	return false
}

func isBuiltin(pkgj string) bool {

	for _, v := range primitives {
		if v == pkgj || pkgj == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func IsInSlice(qry string, slic []string) bool {
	for _, j := range slic {
		if j == qry {
			return true
		}
	}
	return false
}

func DoSpin(chn chan int) {
	s := spin.New()
	s.Set(spin.Spin4)
	var lck bool
	go func() {
		<-chn
		lck = true
	}()
	for !lck {
		fmt.Printf("\r  \033[36mcomputing\033[m %s ", s.Next())
		time.Sleep(160 * time.Millisecond)
	}
	fmt.Printf("\n")
	return
}

func (d *gos) DeleteEnd(id string) {

	temp := []Endpoint{}
	for _, v := range d.Endpoints.Endpoints {
		if v.Id != id {
			temp = append(temp, v)
		}
	}
	d.Endpoints.Endpoints = temp

}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CheckFile(source string) (file, line, reason string) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.

	o, err := parser.ParseFile(fset, source, nil, parser.SpuriousErrors)
	if err != nil {
		eb := strings.Split(err.Error(), ":")
		if len(eb) > 3 {
			file = eb[0]
			line = eb[1]
			reason = eb[3]
		}
		return

	}

	if len(o.Unresolved) > 0 {
		file = "UR"
		line = "UR"
		nset := []string{}

		for _, v := range o.Unresolved {
			if !IsInImports(o.Imports, v.Name) && !IsInSlice(v.Name, nset) && !isBuiltin(v.Name) {
				nset = append(nset, v.Name)
			}
		}

		reason = strings.Join(nset, ",")

	}

	return
}

func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, er := os.Stat(source)
		if er != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

//updates
// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		err = os.MkdirAll(dest, fi.Mode())
		if err != nil {
			return err
		}
	}

	// create dest dir

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = os.MkdirAll(dfp, fi.Mode())
			if err != nil {
				log.Println(err)
			}
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return
}

// A struct for returning custom error messages
type CustomError struct {
	What string
}

// Returns the error message defined in What as a string
func (e *CustomError) Error() string {
	return e.What
}

func NewLen(length int) string {
	return NewLenChars(length, StdChars)
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

//Get config of current project.
func Config() (*gos, error) {
	return LoadGos("./gos.gxml")
}

func (template *gos) AddToMainFunc(str string) error {
	body, err := ioutil.ReadFile(template.Output)
	if err != nil {
		return err
	}

	strnew := strings.Replace(string(body), "//+++extendgxmlmain+++", fmt.Sprintf(`
		%s
		//+++extendgxmlmain+++`, str), 1)
	var strbytes = []byte(strnew)
	_ = ioutil.WriteFile(template.Output, strbytes, 0700)

	return nil
}

func NewID(length int) string {
	return NewLenChars(length, StdNums)
}

// NewLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

func Process(template *gos, r string, web string, tmpl string) (local_string string) {
	// r = GOHOME + GoS Project
	arch := gosArch{}

	var pk []string
	if strings.Contains(os.Args[1], "--") {
		pwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		pwd = strings.Replace(pwd, "\\", "/", -1)
		pk = strings.Split(strings.Trim(pwd, "/"), "/")

	} else {
		pk = strings.Split(strings.Trim(os.Args[2], "/"), "/")
	}


	if template.Type == "webapp" || template.Type == "faas" || template.Type == "package" {

		packageLibrary := ""
		var packageImports = `
		// File generated by Gopher Sauce
		// DO NOT EDIT!!
		package exported
				`

		if template.Type == "webapp" {
			var pprofadd string
			

			if !template.Prod {
				pprofadd = `_ "net/http/pprof"
			`
			}
			local_string = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package main 
		
import (
		gosweb "github.com/cheikhshift/gos/web"
	   %s`, pprofadd)



		} else {
			local_string = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package %s 
import (
		gosweb "github.com/cheikhshift/gos/web"
	 	`, pk[len(pk)-1])

	 		
		}





		for _, pkg := range template.Packages {
			packageImports += fmt.Sprintf(`import "%s"
			`, pkg.Path)

			pathComp := strings.Split(pkg.Path, "/")
			pkgNm := pathComp[len(pathComp) - 1]

			packageLibrary += fmt.Sprintf(`
			func Load%s() %s.PKG { 
				return %s.PKG{}
			}`, pkg.Name, pkgNm, pkgNm)
		}
		// process packages
		packageFinal := fmt.Sprintf(`%s

			%s`, packageImports, packageLibrary)

		os.MkdirAll("./api/exported", 0700)
		ioutil.WriteFile("./api/exported/packages.go", []byte(packageFinal), 0700)

		wd, _ := os.Getwd()
		if !strings.Contains(runtime.GOOS, "indows") {
			wd = strings.Replace(wd, filepath.Join( os.ExpandEnv("$GOPATH"), "src" ) + "/","", 1)
		} else {
			wd = strings.Replace(wd, filepath.Join( os.ExpandEnv("$GOPATH"), "src" ) + "\\","", 1)
		}

		os.MkdirAll("./api/templates", 0700)
		os.MkdirAll("./api/handlers", 0700)
		os.MkdirAll("./api/globals", 0700)

		os.MkdirAll("./api/methods", 0700)
		os.MkdirAll("./api/sessions", 0700)
		os.MkdirAll("./types", 0700)


		var appname = strings.TrimSuffix(strings.Split(strings.Join(pk, "/"), "/src/")[1], "/")


		// if template.Type == "webapp" {
		log.Println("Checking templates")

		if _, err := os.Stat(fmt.Sprintf("%s%s", TrimSuffix(os.ExpandEnv("$GOPATH"), "/"), "/src/github.com/gotpl/gtfmt/")); os.IsNotExist(err) {
			RunCmd("go get github.com/gotpl/gtfmt")
		}

		if _, err := os.Stat(fmt.Sprintf("%s%s", TrimSuffix(os.ExpandEnv("$GOPATH"), "/"), "/src/golang.org/x/tools/cmd/goimports/")); os.IsNotExist(err) {
			RunCmd("go get golang.org/x/tools/cmd/goimports")
		}

		for _, templ := range template.Templates.Templates {
			log.Println("Checking : ", templ.Name)
			var goatResponse string
			if !strings.Contains(runtime.GOOS, "indows") {
				goatResponse, _ = RunCmdSmart(fmt.Sprintf("gtfmt %s/%s.tmpl", tmpl, templ.TemplateFile))
			} else {
				goatResponse, _ = RunCmdSmart(fmt.Sprintf("gtfmt %s\\%s.tmpl", tmpl, strings.Replace(templ.TemplateFile, "/", "\\", -1)))
			}
			if goatResponse != "" {
				color.Yellow("Warning!!!!!!")
				log.Println(goatResponse)
			}
		}

		var TraceOpt, TraceOpen, TraceParam, TraceinFunc, TraCFt, TraceGet, TraceTemplate, TraceError string

		Netimports := []string{"net/http", "time", "github.com/gorilla/sessions", "github.com/cheikhshift/db", "bytes", "encoding/json", "fmt", "html", "html/template", "github.com/fatih/color", "strings", "log", "github.com/elazarl/go-bindata-assetfs"}
		if strings.Contains(template.Type, "webapp") {
			Netimports = append(Netimports, "os")
			Netimports = append(Netimports, "os/signal")
			Netimports = append(Netimports, "context")

			Netimports = append(Netimports, fmt.Sprintf(`sessionStore "%s/api/sessions"`, appname))
			Netimports = append(Netimports, fmt.Sprintf(`%s/api/handlers`, appname))
			Netimports = append(Netimports, fmt.Sprintf(`%s/types`, appname))
			Netimports = append(Netimports, fmt.Sprintf(`%s/api/assets`, appname))



		}
		/*
			Methods before so that we can create to correct delegate method for each object
		*/
		TraceParam = `)`
		TraceinFunc = `)`
		TraCFt = `)`

		if !template.Prod && template.Type != "faas" && template.Type != "package" {
			TraCFt = `, opentracing.Span)`
			TraceParam = `, span)`
			TraceinFunc = `, span opentracing.Span)`
			TraceOpt = `//wheredefault`
			TraceOpen = ` span := opentracing.StartSpan(fmt.Sprintf("%s %s",r.Method,r.URL.Path) )
  				defer span.Finish()
  			  carrier := opentracing.HTTPHeadersCarrier(r.Header)
			if err := span.Tracer().Inject(span.Context(),opentracing.HTTPHeaders,carrier);  err != nil {
		        log.Fatalf("Could not inject span context into header: %v", err)
		    }

`
			TraceTemplate = ` 
  				var sp opentracing.Span
			    opName := fmt.Sprintf("Building template %s%s", p.R.URL.Path, ".tmpl")
			  
			  if true {
			   carrier := opentracing.HTTPHeadersCarrier(p.R.Header)
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier); if err != nil {
			        sp = opentracing.StartSpan(opName)
			    } else {
			        sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			    }
			}
			  defer sp.Finish()
		`
			TraceGet = ` 
  				var sp opentracing.Span
			    opName := fmt.Sprintf(fmt.Sprintf("Web:/%s", r.URL.Path) )
			  
			  if true {
			   carrier := opentracing.HTTPHeadersCarrier(r.Header)
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier); if err != nil {
			        sp = opentracing.StartSpan(opName)
			    } else {
			        sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			    }
			}
			  defer sp.Finish()
		`
			TraceError = ` span.SetTag("error", true)
            span.LogEvent(fmt.Sprintf("%s request at %s, reason : %s ", r.Method, r.URL.Path, err) )`
			Netimports = append(Netimports, "github.com/opentracing/opentracing-go")

			Netimports = append(Netimports, `net`, `net/url`, `sourcegraph.com/sourcegraph/appdash`, `appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"`, `sourcegraph.com/sourcegraph/appdash/traceapp`)
		}

		for _, imp := range template.Methods.Methods {
			if !contains(available_methods, imp.Name) {
				available_methods = append(available_methods, imp.Name)
			}
		}
		apiraw := ``
		for _, imp := range template.Endpoints.Endpoints {
			imp.Method = strings.Replace(imp.Method, `&lt;`, `<`, -1)
			est := ``
			if !template.Prod && template.Type != "faas" {
				est = fmt.Sprintf(`	
					lastLine := ""
						var sp opentracing.Span
					    opName := fmt.Sprintf(" [%s]%%s %%s", r.Method,r.URL.Path)
					  
					  if true {
					   carrier := opentracing.HTTPHeadersCarrier(r.Header)
					wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier); if err != nil {
					        sp = opentracing.StartSpan(opName)
					    } else {
					        sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
					    }
						}
					 	defer sp.Finish()
					defer func() {
					       if n := recover(); n != nil {
					          log.Println("Web request (%s) failed at line :",gosweb.GetLine("%s", lastLine),"Of file:%s :", strings.TrimSpace(lastLine))
					          log.Println("Reason : ",n)
					         %s
					          span.SetTag("error", true)
            span.LogEvent(fmt.Sprintf("%%s request at %%s, reason : %%s ", r.Method, r.URL.Path, n) )
					         	 
					         	w.WriteHeader(http.StatusInternalServerError)
							    w.Header().Set("Content-Type",  "text/html")
								pag,err := templates.LoadPage("%s")
											
								 if err != nil {
								        	log.Println(err.Error())
								        	callmet = true	        	
								        	return
								 }
								 	 pag.R = r
						         pag.Session = session	
								if pag.IsResource {
				        			w.Write(pag.Body)
						    	} 

								 callmet = true
					        }
						}()`, imp.Id, imp.Path, template.Name, template.Name, TraceOpt, template.ErrorPage)
				setv := strings.Split(imp.Method, "\n")
				for _, line := range setv {
					line = strings.TrimSpace(line)
					if len(line) > 0 {
						est += fmt.Sprintf("\nlastLine =  `%s`\n%s", line, line)
					}
				}

			} else {
				est = imp.Method
			}

			if imp.Type == "f" {
				var tlc = strings.Replace(fmt.Sprintf("%s%s", imp.Type, strings.Replace(strings.Title(strings.Replace(imp.Path, "/", " ", -1)), " ", "", -1)), "star", "", -1)

				est = fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package handlers

				import gosweb "github.com/cheikhshift/gos/web"

				import sessionStore "%s/api/sessions"

				import templates "%s/api/templates"

				import methods "%s/api/methods"

				import types "%s/types"
				import "%s/api/globals"



				func %s(w http.ResponseWriter, r *http.Request ,session *sessions.Session%s (response string, callmet bool){
					%s
					return
				}`, appname,appname,appname,appname, appname, tlc,  TraceinFunc, est)

				ioutil.WriteFile(fmt.Sprintf("./api/handlers/middleware_%s.go", tlc), []byte(est), 0700 )

				

				apiraw += fmt.Sprintf(` 
					if   strings.Contains(r.URL.Path, "%s")  { 
						response,callmet = %s(w,r,session%s				
					}
					`, imp.Path, tlc, strings.Replace(TraceinFunc, "opentracing.Span", "", -1) ) 
			}
		}

		apiraw += `if r.Method == "RESET" {
			return true
		}`
		for _, imp := range template.Endpoints.Endpoints {
			est := ``
			if !template.Prod && template.Type != "faas" {
				est = fmt.Sprintf(`	
					lastLine := ""
					var sp opentracing.Span
					    opName := fmt.Sprintf(" [%s]%%s %%s", r.Method,r.URL.Path)
					  
					  if true {
					   carrier := opentracing.HTTPHeadersCarrier(r.Header)
					wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier); if err != nil {
					        sp = opentracing.StartSpan(opName)
					    } else {
					        sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
					    }
						}
					  defer sp.Finish()
					defer func() {
					       if n := recover(); n != nil {
					          log.Println("Web request (%s) failed at line :",gosweb.GetLine("%s", lastLine),"Of file:%s :", strings.TrimSpace(lastLine))
					          log.Println("Reason : ",n)
					          %s
					       	 span.SetTag("error", true)
            span.LogEvent(fmt.Sprintf("%%s request at %%s, reason : %%s ", r.Method, r.URL.Path, n) )
					        	 w.WriteHeader(http.StatusInternalServerError)
							    w.Header().Set("Content-Type",  "text/html")
								pag,err := templates.LoadPage("%s")
				   
								 if err != nil {
								        	log.Println(err.Error())
								        	callmet = true	        	
								        	return 
								 }
								  pag.R = r
						         pag.Session = session	
								   if pag.IsResource {
							        	w.Write(pag.Body)
							    	} else {
							    		templates.RenderTemplate(w, pag%s //"s" 
							     
							    	}
					           callmet = true
					        
					        }
						}()
						`, imp.Id, imp.Path, template.Name, template.Name, TraceOpt, template.ErrorPage, TraceParam)
				setv := strings.Split(imp.Method, "\n")
				for _, line := range setv {
					line = strings.TrimSpace(line)
					if len(line) > 0 {
						est += fmt.Sprintf("\nlastLine =  `%s`\n%s", line, line)
					}
				}

			} else {
				est = strings.Replace(imp.Method, `&#38;`, `&`, -1)
			}
			if !strings.Contains(est, "w.Write") && !strings.Contains(est, "response") && imp.Type != "f" {
				color.Yellow(fmt.Sprintf("Warning : No response writing detected with endpoint : %s type : %s", imp.Path, imp.Type))
			}

			
			var tlc = strings.Replace( strings.Replace(fmt.Sprintf("%s%s", imp.Type, strings.Replace(strings.Title(strings.Replace(imp.Path, "/", " ", -1)), " ", "", -1)), "star", "", -1), "-", "", -1)


			if imp.Type == "star" || imp.Type != "f" {
				est = fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package handlers

				import gosweb "github.com/cheikhshift/gos/web"

				import sessionStore "%s/api/sessions"

				import templates "%s/api/templates"

				import methods "%s/api/methods"

				import types "%s/types"
				import "%s/api/globals"

				func %s(w http.ResponseWriter, r *http.Request, session *sessions.Session%s (response string,callmet bool){

					%s

					callmet = true
					return
				}`, appname,appname,appname,appname, appname, tlc, TraceinFunc,est)

				ioutil.WriteFile(fmt.Sprintf("./api/handlers/rest_%s.go", tlc), []byte(est), 0700 )
			}


			if imp.Type == "star" {
				apiraw += fmt.Sprintf(` else if  !callmet &&  gosweb.UrlAtZ(r.URL.Path, "%s")  { 
					

					response, callmet = %s(w, r, session%s			
					
				}`, imp.Path, tlc, strings.Replace(TraceinFunc, "opentracing.Span", "", -1) )
			} else if imp.Type != "f" {

				apiraw += fmt.Sprintf(` else if  isURL := (r.URL.Path == "%s" && r.Method == strings.ToUpper("%s") );!callmet && isURL{ 
				
					response, callmet = %s(w, r, session%s
					
				} `, imp.Path, imp.Type, tlc, strings.Replace(TraceinFunc, "opentracing.Span", "", -1) )
			}

		}

		timeline := ``
		for _, imp := range template.Timers.Timers {

			timeline += `
			` + imp.Name + ` := time.NewTicker(time.` + imp.Unit + ` * ` + imp.Interval + `)
					    go func() {
					        for _ = range ` + imp.Name + `.C {
					           ` + strings.Replace(imp.Method, `&#38;`, `&`, -1) + `
					        }
					    }()
    `
		}

		//log.Printf("APi Methods %v\n",api_methods)
		netMa := `template.FuncMap{"a":gosweb.Netadd,"s":gosweb.Netsubs,"m":gosweb.Netmultiply,"d":gosweb.Netdivided,"js" : gosweb.Netimportjs,"css" : gosweb.Netimportcss,"sd" : gosweb.NetsessionDelete,"sr" : gosweb.NetsessionRemove,"sc": gosweb.NetsessionKey,"ss" : gosweb.NetsessionSet,"sso": gosweb.NetsessionSetInt,"sgo" : gosweb.NetsessionGetInt,"sg" : gosweb.NetsessionGet,"form" : gosweb.Formval,"eq": gosweb.Equalz, "neq" : gosweb.Nequalz, "lte" : gosweb.Netlt`
		for _, imp := range available_methods {
			if !contains(api_methods, imp) && template.findMethod(imp).Keeplocal != "true" {
				netMa += fmt.Sprintf(`,"%s" : methods.%s`, imp, imp)
			}
		}


		for _, pkg := range template.Packages {
			netMa += fmt.Sprintf(`,"%s" : exported.Load%s`, pkg.Name, pkg.Name)
		}


		//int_lok := []string{}

		/*	for _,imp := range template.RootImports {
					//log.Println(imp)
				if strings.Contains(imp.Src,".gxml") {

					pathsplit := strings.Split(imp.Src,"/")
					gosName := pathsplit[len(pathsplit) - 1]
					pathsplit = pathsplit[:len(pathsplit)-1]
					if _, err := os.Stat(TrimSuffix(os.ExpandEnv("$GOPATH"), "/" ) + "/src/"  + strings.Join(pathsplit,"/")); os.IsNotExist(err){
							color.Red("Package not found")
							log.Println("∑ Downloading Package " + strings.Join(pathsplit,"/"))
							RunCmdSmart("go get -v " + strings.Join(pathsplit,"/"))
					}
					//split and replace last section
					log.Println("∑ Processing XML Yåå ", pathsplit)
					xmlPackageDir := TrimSuffix(os.ExpandEnv("$GOPATH"), "/" ) + "/src/" + strings.Join(pathsplit,"/") + "/"
						//copy gole with given path -
						log.Println("Installing Resources into project!")
						//delete prior to copy
					//	RemoveContents(r + "/" + web + "/" + xml_iter.Package)
					//	RemoveContents(r + "/" + tmpl + "/" + xml_iter.Package)
					//	CopyDir(xmlPackageDir + xml_iter.Web, r + "/" + web + "/" + xml_iter.Package)
					//	CopyDir(xmlPackageDir + xml_iter.Tmpl, r + "/" + tmpl + "/" + xml_iter.Package)
					//	template.MergeWith( xmlPackageDir + gosName)
					//	log.Println(template)

				}
			}
		*/
		for _, imp := range template.RootImports {
			if !strings.Contains(imp.Src, ".gxml") {
				//log.Println(TrimSuffix(os.ExpandEnv("$GOPATH"), "/" ) + "/src/" + imp.Src )
				if validimport := (!contains(Netimports, imp.Src) && imp.Src != "io/ioutil"); validimport {
					Netimports = append(Netimports, imp.Src)
				}
			}
		}

		/*
			enable if needed
			for _, imp := range template.Header.Objects {
				//struct return and function

				if !contains(int_lok, imp.Name) {
					int_lok = append(int_lok, imp.Name)
					netMa += `,"` + imp.Name + `" : Net` + imp.Name
				}
			} */

		for _, imp := range template.Templates.Templates {

			netMa += fmt.Sprintf(`,"%s" : Net%s`, imp.Name, imp.Name)
			netMa += fmt.Sprintf(`,"b%s" : Netb%s`, imp.Name, imp.Name)
			netMa += fmt.Sprintf(`,"c%s" : Netc%s`, imp.Name, imp.Name)
		}

		//	log.Println(template.Methods.Methods[0].Name)

		for _, imp := range Netimports {

			if hasQts := strings.Contains(imp, `"`); hasQts {
				local_string += fmt.Sprintf(`
			 %s`, imp)
			} else {
				local_string += fmt.Sprintf(`
			"%s"`, imp)
			}
		}
		var structs_string, struct_funcs string
		//Lets Do structs
		structs_string = ``
		for _, imp := range template.Header.Structs {
			if !contains(arch.objects, imp.Name) {
				log.Println("🔧 Processing Struct : ", imp.Name)
				arch.objects = append(arch.objects, imp.Name)

				structs_string += fmt.Sprintf(`
			type %s struct {`, imp.Name)
				structs_string += imp.Attributes
				structs_string += "}"

			struct_funcs += fmt.Sprintf(`

			// Assert first argument as struct
			// %s
			func  Netcast%s(args ...interface{}) *%s  {
				
				s := %s{}
				mapp := args[0].(db.O)
				if _, ok := mapp["_id"]; ok {
					mapp["Id"] = mapp["_id"]
				}
				data,_ := json.Marshal(&mapp)
				
				err := json.Unmarshal(data, &s) 
				if err != nil {
					log.Println(err.Error())
				}
				
				return &s
			}
			func Netstruct%s() *%s{ return &%s{} }`, imp.Name,imp.Name,imp.Name, imp.Name, imp.Name, imp.Name, imp.Name)



			netMa += fmt.Sprintf(`,"%s" : types.Netstruct%s`, imp.Name, imp.Name)
			netMa += fmt.Sprintf(`,"is%s" : types.Netcast%s`, imp.Name, imp.Name)

			}
		}

		netMa += `}`

		ReadyTemplate := "" //"func ReadyTemplate(body []byte) string { return strings.Replace(strings.Replace(strings.Replace(string(body), \"/{\", \"\\\"{\",-1),\"}/\", \"}\\\"\",-1 ) ,\"`\", \"\\\"\" ,-1) }"
		netmafuncs := netMa
		var CacheParam string
		var pkgStruct string


		if template.Type == "package" {
			pkgStruct = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package %s
			type PKG struct {}`, pk[len(pk) - 1] )

			ioutil.WriteFile("./entry.go", []byte(pkgStruct), 0700)
		}

		if template.Prod {
			CacheParam = `
			if lPage, ok := WebCache.Get(title); ok {
				return &lPage, nil
			}`
		}
		netMa = `TemplateFuncStore`
		local_string += fmt.Sprintf(`
		)`)

		globals := ""
		for _, imp := range template.Variables {
			globals += fmt.Sprintf(`
						var %s %s`, imp.Name, imp.Type)
		}

		globals = fmt.Sprintf(`
		// File generated by Gopher Sauce
		// DO NOT EDIT!!		
		package globals

			%s`, globals)

		ioutil.WriteFile("./api/globals/variables.go", []byte(globals), 0700)


		if template.Init_Func != "" {
			local_string += fmt.Sprintf(`
			func init(){
				%s
			}`, template.Init_Func)

		}

		handlersHtml := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package handlers

			import gosweb "github.com/cheikhshift/gos/web"
			import templates "%s/api/templates"
			import sessionStore "%s/api/sessions"
			import "%s/api/globals"

			// Access you .gxml's end tags with
				// this http.HandlerFunc.
				// Use MakeHandler(http.HandlerFunc) to serve your web
				// directory from memory.
				func MakeHandler(fn func (http.ResponseWriter, *http.Request%s) http.HandlerFunc {
				  return func(w http.ResponseWriter, r *http.Request) {
				  
				  	 %s
				  	
				  	if attmpt := ApiAttempt(w,r%s ;!attmpt {
				       fn(w, r%s
				  	} 
				  			  	
				  	
				  }
				} 

				

				
			func Handler(w http.ResponseWriter, r *http.Request%s {
				  var p *gosweb.Page
				  p,err := templates.LoadPage(r.URL.Path)
				  	var session *sessions.Session
				  	var er error
				  	if 	session, er = sessionStore.Store.Get(r, "session-"); er != nil {
						session,_ = sessionStore.Store.New(r, "session-")
					}
				  %s
				  if err != nil {	
				  		log.Println(err.Error())
				  		
				        w.WriteHeader(http.StatusNotFound)				  	
				       	%s
				        pag,err := templates.LoadPage("%s")
				        
				        if err != nil {
				        	log.Println(err.Error())
				        	//
				        	return
				        }
				         pag.R = r
						 pag.Session = session
						if p != nil {
						p.Session = nil
				  		p.Body = nil
				  		p.R = nil
				  		p = nil
				  		}
				  	
				        if pag.IsResource {
				        	w.Write(pag.Body)
				    	} else {
				    		templates.RenderTemplate(w, pag%s //"%s" 
				    	}
				    	session = nil
				    	
				        return
				  }

				   
				  if !p.IsResource {
				  		w.Header().Set("Content-Type",  "text/html")
				  		p.Session = session
				  		p.R = r
				      	templates.RenderTemplate(w, p%s //fmt.Sprintf("%s%%s", r.URL.Path)
				     	session.Save(r, w)
				     // log.Println(w)
				  } else {
				  		w.Header().Set("Cache-Control",  "public")
				  		if strings.Contains(r.URL.Path, ".css") {
				  	  		w.Header().Add("Content-Type",  "text/css")
				  	  	} else if strings.Contains(r.URL.Path, ".js") {
				  	  		w.Header().Add("Content-Type",  "application/javascript")
				  	  	} else {
				  	  	w.Header().Add("Content-Type",  http.DetectContentType(p.Body))
				  	  	}
				  	 
				  	 
				      w.Write(p.Body)
				  }

			  	 p.Session = nil
				 p.Body = nil
				 p.R = nil
				 p = nil
				 session = nil
				 
				 return
				}
				`, appname , appname, appname,TraCFt, TraceOpen, TraceParam, TraceParam, TraceinFunc, TraceGet, TraceError, template.NPage, TraceParam, template.ErrorPage, TraceParam, web)

		ioutil.WriteFile("./api/handlers/html.go", []byte(handlersHtml), 0700)

		apifinal := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package handlers

			import gosweb "github.com/cheikhshift/gos/web"

			import sessionStore "%s/api/sessions"

			import templates "%s/api/templates"

			import methods "%s/api/methods"

			import types "%s/types"
			import "%s/api/globals"

			var WebCache = gosweb.NewCache()

			func mResponse(v interface{}) string {
					data,_ := json.Marshal(&v)
					return string(data)
			}

			func ApiAttempt(w http.ResponseWriter, r *http.Request%s (callmet bool) {
					var response string
					var session *sessions.Session
					
				  	var er error

				  	if 	session, er = sessionStore.Store.Get(r, "session-"); er != nil {
						session,_ = sessionStore.Store.New(r, "session-")
					}
					
					%s
					if callmet {
						session.Save(r,w)
						session = nil
						if response != "" {
							//Unmarshal json
							//w.Header().Set("Access-Control-Allow-Origin", "*")
							w.Header().Set("Content-Type",  "application/json")
							w.Write([]byte(response))
						}
						return 
					}
					session = nil
					return
				}



			`, appname, appname, appname, appname, appname,TraceinFunc, apiraw)

		ioutil.WriteFile("./api/handlers/endpoints.go", []byte(apifinal), 0700)

		sessionStore := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package sessions 

			var Store = sessions.NewCookieStore([]byte("%s"))`, template.Key)

		ioutil.WriteFile("./api/sessions/store.go", []byte(sessionStore), 0700)

	

		for _, imp := range template.Header.Objects {
			structs_string = fmt.Sprintf(`type %s %s
				`, imp.Name, imp.Templ) + structs_string
		}

		structs_string = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package types

			%s`, structs_string)

		struct_funcs = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package types

			%s`, struct_funcs)

		ioutil.WriteFile("./types/struct_funcs.go", []byte(struct_funcs), 0700)
		ioutil.WriteFile("./types/structs.go", []byte(structs_string), 0700)

		var  methodsExported string 

		for _, imp := range available_methods {
			if !contains(int_methods, imp) && !contains(api_methods, imp) {
				log.Println("🚰 Processing : ", imp)
				var methodsString string

				meth := template.findMethod(imp)
				commentslice := strings.Split(string(meth.Comment), "\n")
				for i, val := range commentslice {
					commentslice[i] = strings.TrimSpace(val)
				}
				if len(commentslice) > 0 {
					methodsString += fmt.Sprintf(`
						// %s`, strings.Join(commentslice, `
						// `))

					splitAtComment := strings.Split(meth.Method, "-->") //at end of comment
					meth.Method = splitAtComment[len(splitAtComment)-1]
				}
				addedit := false
				if meth.Returntype == "" {
					meth.Returntype = "string"
					addedit = true
				}
				if meth.Man == "exp" {
					methodsString += fmt.Sprintf(`
						func %s(%s) %s {
							`, meth.Name, meth.Variables, meth.Returntype)
				} else {
					methodsString += fmt.Sprintf(`
						func %s(args ...interface{}) %s {
							`, meth.Name, meth.Returntype)
					for k, nam := range strings.Split(meth.Variables, ",") {
						if nam != "" {
							methodsString += fmt.Sprintf(`%s := args[%v]
								`, nam, k)
						}
					}
				}
				meth.Method = strings.Replace(meth.Method, "&lt;", "<", -1)
				est := ``
				if !template.Prod {
					est = fmt.Sprintf(`	
							lastLine := ""

							defer func() {
							       if n := recover(); n != nil {
							          log.Println("Pipeline failed at line :",gosweb.GetLine("%s", lastLine),"Of file:%s:", strings.TrimSpace(lastLine))
							          log.Println("Reason : ",n)

							        }
								}()`, template.Name, template.Name)
					setv := strings.Split(meth.Method, "\n")
					for _, line := range setv {
						line = strings.TrimSpace(line)
						if len(line) > 0 {
							est += fmt.Sprintf("\nlastLine = `%s`\n%s", line, line)
						}
					}

				} else {
					est = strings.Replace(meth.Method, `&#38;`, `&`, -1)
				}
				methodsString += est
				if addedit {
					methodsString += `
						 return ""
						 `
				}
				methodsString += `
						}`

				if template.Type == "package" {
					
					varChain := []string{}

					for _, nam := range strings.Split(meth.Variables, ",") {
						if nam != "" {
							declName := strings.Split(nam, " ")
							varChain = append(varChain,declName[0])
						}
					}

					if meth.Man == "exp" {

						methodsExported += fmt.Sprintf(`
							func (pkg PKG) %s(%s) %s {
								return methods.%s(%s) 
							}
							`, meth.Name, meth.Variables, meth.Returntype, meth.Name, strings.Join(varChain, ","))
					} else {
						methodsExported += fmt.Sprintf(`
							func (pkg PKG) %s(args ...interface{}) %s {
								return methods.%s(args...)
							}`, meth.Name, meth.Returntype, meth.Name)
							
					}


				}

					methodsString = fmt.Sprintf(`
					// File generated by Gopher Sauce
					// DO NOT EDIT!!
					package methods

					import gosweb "github.com/cheikhshift/gos/web"
					import "%s/types"
					import "%s/api/globals"

					%s`, appname,appname,methodsString)

					ioutil.WriteFile(fmt.Sprintf("./api/methods/method_%s.go", meth.Name ), []byte(methodsString), 0700)
			}

		

		}

		

		if template.Type == "package" {

			methodsExported = fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package %s

				import "%s/api/methods"

				%s`, pk[len(pk) - 1], appname, methodsExported)

			ioutil.WriteFile("./exported_methods.go", []byte(methodsExported), 0700)
		}

		var templatesFile,templateImport string

		for _, imp := range template.Templates.Templates {
			if imp.Struct == "" {
				imp.Struct = "gosweb.NoStruct"
			}

			imp.Struct = strings.Replace(imp.Struct,"*","", -1)

			commentslice := strings.Split(string(imp.Comment), "\n")
			var commentstring string
			for i, val := range commentslice {
				commentslice[i] = strings.TrimSpace(val)
			}
			if len(commentslice) > 0 {
				commentstring = strings.TrimSpace(fmt.Sprintf(`
						// %s`, strings.Join(commentslice, `
						// `)))

			}

			if !strings.Contains(imp.Struct, ".") {
				imp.Struct = fmt.Sprintf("types.%s", strings.Title(imp.Struct ) )
			}

			templateString := fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package templates
				
				import "%s/types"
				import "%s/api/assets"
				import gosweb "github.com/cheikhshift/gos/web"

				// Render HTML of template 
				// %s with struct %s
				func %s(d %s) string {
					return  Netb%s(d)
				}

				// recovery function used to log a 
				// panic.
				func templateFN%s(localid string, d interface{}) {
					    if n := recover(); n != nil {
					           	   color.Red(fmt.Sprintf("Error loading template in path (%s) : %%s" , localid ) )
					           	// log.Println(n)
					           		DebugTemplatePath(localid, d)	
						}
				}

				var  templateID%s = "%s/%s.tmpl"


				func  Net%s(args ...interface{}) string {
					
					localid := templateID%s
					var d *%s
					defer templateFN%s(localid, d)	
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, d)
					if err != nil {
						return err.Error()
					}
					} else {
						d = &%s{}
					}

					
    				
    				 output := new(bytes.Buffer) 
 	
    				 if _, ok := templateCache.Get(localid); !ok || !Prod {

    				 	body, er := assets.Asset(localid)
		    				if er != nil {
		    					return ""
		    			}
		    			var localtemplate =  template.New("%s")	  			
		    			localtemplate.Funcs(%s)
		    			var tmpstr = string(body)
				  		localtemplate.Parse(tmpstr)
	    				body = nil
	    				templateCache.Put(localid, localtemplate )
    				 }
					

					erro := templateCache.JGet(localid).Execute(output, d)
				    if erro != nil {
				   	color.Red(fmt.Sprintf("Error processing template %%s" , localid) )
				  	 DebugTemplatePath(localid, d)	
				    } 
				    var outps = output.String()
				    var outpescaped = html.UnescapeString(outps)
				    d = nil
				    output.Reset()
				    output = nil
				    args = nil
					return outpescaped
					
				}
				func b%s(d %s) string {
						return Netb%s(d)
				}

				%s

				func  Netb%s(d %s) string {
					localid := templateID%s
					defer templateFN%s(localid, d)
    				 output := new(bytes.Buffer) 
				  	
    				 if _, ok := templateCache.Get(localid); !ok || !Prod {

    				 	body, er := assets.Asset(localid)
		    				if er != nil {
		    					return ""
		    			}
		    			var localtemplate =  template.New("%s")	  			
		    			localtemplate.Funcs(%s)
		    			var tmpstr = string(body)
				  		localtemplate.Parse(tmpstr)
	    				body = nil
	    				templateCache.Put(localid, localtemplate )
    				 }
					

					erro := templateCache.JGet(localid).Execute(output, d)
				    if erro != nil {
				    log.Println(erro)
				    } 
					var outps = output.String()
				    var outpescaped = html.UnescapeString(outps)
				    d = %s{}
				    output.Reset()
				    output = nil
					return outpescaped
				}
				func  Netc%s(args ...interface{}) (d %s) {
					if len(args) > 0 {
					var jsonBlob = []byte(args[0].(string))
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						log.Println("error:", err)
						return 
					}
					} else {
						d = %s{}
					}
    				return
				}

				func  c%s(args ...interface{}) (d %s) {
					if len(args) > 0 {
						d = Netc%s(args[0])
					} else {
						d = Netc%s()
					}
    				return
				}

				

			
				`,  appname, appname , imp.Name, imp.Struct,strings.Title(imp.Name), imp.Struct, imp.Name ,imp.Name, imp.TemplateFile, imp.Name, tmpl, imp.TemplateFile, imp.Name, imp.Name, imp.Struct, imp.Name, imp.Struct, imp.Name, netMa, imp.Name, imp.Struct, imp.Name, commentstring, imp.Name, imp.Struct, imp.Name, imp.Name, imp.Name, netMa, imp.Struct, imp.Name, imp.Struct, imp.Struct, imp.Name, imp.Struct, imp.Name, imp.Name)
				
				

				if  template.Type == "package" {
					templateImport +=  fmt.Sprintf(`func (pkg PKG) %s(args ...interface{}) string {
					var result string

					if len(args) > 0 {
						result = templates.Netb%s(args[0].(%s))
					} else {
						result = templates.Net%s()
					}
    				return result
				}`, imp.Name, imp.Name, imp.Struct, imp.Name)

				}


				ioutil.WriteFile( fmt.Sprintf("./api/templates/template_%s.go", imp.Name), []byte(templateString), 0700 )
		}

		templatesFile = fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package templates

			import gosweb "github.com/cheikhshift/gos/web"
			
			import methods "%s/api/methods"
			import exported "%s/api/exported"
			import "%s/api/assets"

			import "%s/types"
			import "%s/api/globals"

			var TemplateFuncStore template.FuncMap
			
			var templateCache = gosweb.NewTemplateCache()
			var WebCache = gosweb.NewCache()

			%s 

			%s

			func LoadPage(title string) (*gosweb.Page,error) {
				   	
					%s

					var nPage = gosweb.Page{}
				    if roottitle := (title == "/"); roottitle  {
				    	webbase := "%s/"
					    	fname := fmt.Sprintf("%%s%%s", webbase, "index.html")
					    	body, err := assets.Asset(fname)
					    	if err != nil {
					    		fname = fmt.Sprintf("%%s%%s", webbase, "index.tmpl")
					    		body , err = assets.Asset(fname)
					    		if err != nil {
					    			return nil,err
					    		}
					    		nPage.Body = body
					    		WebCache.Put(title, nPage)
					    		body = nil
					    		return  &nPage, nil
					    	}
					    	nPage.Body = body
					    	nPage.IsResource = true
					    	WebCache.Put(title, nPage)
					    	body = nil
					    	return  &nPage, nil
					    		    		
				     } 
				     
				   filename := fmt.Sprintf("%s%%s.tmpl", title)

				   if body, err := assets.Asset(filename) ;err != nil {
				    	 filename = fmt.Sprintf("%s%%s.html", title) 
				    	
				    	if  body, err = assets.Asset(filename); err != nil {
				         filename = fmt.Sprintf("%s%%s", title) 
				         
				         if  body, err = assets.Asset(filename); err != nil {
				            return nil, err
				         } else {
				          if strings.Contains(title, ".tmpl")  {
				              return nil,nil
				          }
					    	nPage.Body = body
					    	nPage.IsResource = true
					    	WebCache.Put(title, nPage)
					    	body = nil
				            return &nPage, nil
				         }
				      } else {					    	
				      	nPage.Body = body
					    	nPage.IsResource = true
					    	WebCache.Put(title, nPage)
					    body = nil
				         return &nPage, nil
				      }
				    } else {
				    						    	nPage.Body = body
					    	WebCache.Put(title, nPage)
					    	body = nil
				    	  return &nPage, nil
				    }
 
				       %s
				  
				 } 
				
				   
			
				 %s`, appname, appname, appname, appname, appname,fmt.Sprintf(`func StoreNetfn () int {
				 	// List of pipelines linked to each template.
				 	TemplateFuncStore = %s
				 	return 0
				 	}
				 	var FuncStored = StoreNetfn()`, netmafuncs),templatesFile, CacheParam, web, web, web, web, TraceOpt, ReadyTemplate)



		ioutil.WriteFile("./api/templates/templates.go", []byte(templatesFile), 0700)

		loggerFile := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package templates

			import gosweb "github.com/cheikhshift/gos/web"

			import sessionStore "%s/api/sessions"
			import "%s/api/assets"

			var Prod = %v

			func DebugTemplate(w http.ResponseWriter,r *http.Request,tmpl string){
					lastline := 0
					linestring := ""
					defer func() {
					       if n := recover(); n != nil {
					           	log.Println()
					           	// log.Println(n)
					           			log.Println("Error on line :", lastline + 1 ,":" + strings.TrimSpace(linestring)) 
					           	 //http.Redirect(w,r,"%s",307)
					        }
					    }()	

					p,err := LoadPage(r.URL.Path)
					filename :=  tmpl  + ".tmpl"
				    body, err := assets.Asset(filename)
				    session, er := sessionStore.Store.Get(r, "session-")

				 	if er != nil {
				           session,er = sessionStore.Store.New(r,"session-")
				    }
				    p.Session = session
				    p.R = r
				    if err != nil {
				       	log.Print(err)
				    	
				    } else {
				    
				  
				   
				    lines := strings.Split(string(body), "\n")
				   // log.Println( lines )
				    linebuffer := ""
				    waitend := false
				    open := 0
				    for i, line := range lines {
				    	
				    	processd := false
				    	

				    	if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				    		linebuffer += line
				    		waitend = true
				    		
				    		endstr := ""
				    		processd = true
				    		if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}") ) {
				    
				    			open++;
					    		
				    		}
				    		for i := 0; i < open; i++ {
				    			endstr += "\n{{end}}"
				    		}
				    		//exec
				    		outp := new(bytes.Buffer)  
					    	t := template.New("PageWrapper")
					    	t = t.Funcs(%s)
					    	t, _ = t.Parse(string(body))
					    	lastline = i
					    	linestring =  line
					    	erro := t.Execute(outp, p)
						    if erro != nil {
						   		log.Println("Error on line :", i + 1,line,erro.Error())   
						    } 
				    	}
				   

				    	if waitend && !processd && !( strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") ) {
				    		linebuffer += line

				    		endstr := ""
				    		for i := 0; i < open; i++ {
				    			endstr += "\n{{end}}"
				    		}
				    		//exec
				    		outp := new(bytes.Buffer)  
					    	t := template.New("PageWrapper")
					    	t = t.Funcs(%s)
					    	t, _ = t.Parse(string(body) )
					    	lastline = i
					    	linestring =  line
					    	erro := t.Execute(outp, p)
						    if erro != nil {
						   		log.Println("Error on line :", i + 1,line,erro.Error())   
						    } 

				    	}



				    	if !waitend && !processd {
				    	outp := new(bytes.Buffer)  
				    	t := template.New("PageWrapper")
				    	t = t.Funcs(%s)
				    	t, _ = t.Parse(string(body) )
				    	lastline = i
				    	linestring = line
				    	erro := t.Execute(outp, p)
					    if erro != nil {
					   		log.Println("Error on line :", i + 1,line,erro.Error())   
					    }  
						}

						if  !processd && ( strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") ) {
							open--

							if open == 0 {
							waitend = false
				    		
							}
				    	}
				    }
				    
					
				    }

				}

			func DebugTemplatePath(tmpl string, intrf interface{}){
					lastline := 0
					linestring := ""
					defer func() {
					       if n := recover(); n != nil {
					         
					           			log.Println("Error on line :", lastline + 1,":" + strings.TrimSpace(linestring)) 
					           			log.Println(n)
					           	 //http.Redirect(w,r,"%s",307)
					        }
					    }()	

				
					filename :=  tmpl  
				    body, err := assets.Asset(filename)
				   
				    if err != nil {
				       	log.Print(err)
				    	
				    } else {
				    
				  
				   
				    lines := strings.Split(string(body), "\n")
				   // log.Println( lines )
				    linebuffer := ""
				    waitend := false
				    open := 0
				    for i, line := range lines {
				    	
				    	processd := false

				   		if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				    		linebuffer += line
				    		waitend = true
				    		
				    		endstr := ""
				    		if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}") ) {
				    
				    			open++;
					    		
				    		}

				    		for i := 0; i < open; i++ {
					    			endstr += "\n{{end}}"
					    	}
				    		//exec

				    		processd = true
				    		outp := new(bytes.Buffer)  
					    	t := template.New("PageWrapper")
					    	t = t.Funcs(%s)
					    	t, _ = t.Parse(string([]byte(fmt.Sprintf("%%s%%s",linebuffer, endstr))) )
					    	lastline = i
					    	linestring =  line	    	
					    	erro := t.Execute(outp, intrf)
						    if erro != nil {
						   		log.Println("Error on line :", i + 1,line,erro.Error())   
						    } 
				    	}



				    	if waitend && !processd && !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}") )  {
				    		linebuffer += line

				    		endstr := ""
				    		for i := 0; i < open; i++ {
				    			endstr += "\n{{end}}"
				    		}
				    		//exec
				    		outp := new(bytes.Buffer)  
					    	t := template.New("PageWrapper")
					    	t = t.Funcs(%s)
					    	t, _ = t.Parse(string([]byte(fmt.Sprintf("%%s%%s",linebuffer, endstr))) )
					    	lastline = i
					    	linestring =  line
					    	erro := t.Execute(outp, intrf)
						    if erro != nil {
						   		log.Println("Error on line :", i + 1,line,erro.Error())   
						    } 

				    	}



				    	if !waitend && !processd {
				    	outp := new(bytes.Buffer)  
				    	t := template.New("PageWrapper")
				    	t = t.Funcs(%s)
					    t, _ = t.Parse(string([]byte(fmt.Sprintf("%%s%%s",linebuffer))) )
				    	lastline = i
				    	linestring = line
				    	erro := t.Execute(outp, intrf)
					    if erro != nil {
					   		log.Println("Error on line :", i + 1,line,erro.Error())   
					    }  
						}

						if  !processd && (strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}") ) {
							open--

							if open == 0 {
							waitend = false
				    		
							}
				    	}
				    }
				    
					
				    }

				}` , appname, appname, template.Prod,template.ErrorPage, netMa, netMa, netMa, template.ErrorPage, netMa, netMa, netMa)

		ioutil.WriteFile("./api/templates/logger.go", []byte(loggerFile), 0700 )

		renderFile := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package templates

			import gosweb "github.com/cheikhshift/gos/web"

			func RenderTemplate(w http.ResponseWriter, p *gosweb.Page%s   {
				     defer func() {
					        if n := recover(); n != nil {
					           	 color.Red(fmt.Sprintf("Error loading template in path : %s%%s.tmpl reason : %%s", p.R.URL.Path,n)  )
					           	 
					           	 DebugTemplate( w,p.R , fmt.Sprintf("%s%%s", p.R.URL.Path) )
					           	 w.WriteHeader(http.StatusInternalServerError)
					           	 
						         pag,err := LoadPage("%s" )
						       
						    
						        if err != nil {
						        	log.Println(err.Error())	        	
						        	return
						        }

						         if pag.IsResource {
						        	w.Write(pag.Body)
						    	} else {
						    		pag.R = p.R
						         	pag.Session = p.Session
						    		RenderTemplate(w, pag%s //%s"
						     
						    	}
					        }
					    }()


				  	%s
				  
				    // %s
		
				 	if _,ok := templateCache.Get(p.R.URL.Path); !ok || !Prod {
				 		var tmpstr = string(p.Body)
				 		var localtemplate =  template.New(p.R.URL.Path)
				 		
				 		localtemplate.Funcs(TemplateFuncStore)
				 		localtemplate.Parse(tmpstr)
				 		templateCache.Put(p.R.URL.Path, localtemplate )
				 	}
	
				    outp := new(bytes.Buffer)
				    err := templateCache.JGet(p.R.URL.Path).Execute(outp, p)
				    if err != nil {
				        log.Println(err.Error())
				    	DebugTemplate( w,p.R , fmt.Sprintf("%s%%s", p.R.URL.Path))
				    	w.WriteHeader(http.StatusInternalServerError)
					    w.Header().Set("Content-Type",  "text/html")
						pag,err := LoadPage("%s" )
						 
						 if err != nil {
						        	log.Println(err.Error())	        	
						        	return
						 }
						 pag.R = p.R
						 pag.Session = p.Session
						    
						  if pag.IsResource {
				        	w.Write(pag.Body)
				    	} else {
				    		RenderTemplate(w, pag%s // "%s" 
				     
				    	}
				    	return
				    } 


				 	// p.Session.Save(p.R, w)

				    var outps = outp.String()
				    var outpescaped = html.UnescapeString(outps)
				    outp = nil
				    fmt.Fprintf(w, outpescaped )
					
				    
				}`, TraceinFunc, web, web, template.ErrorPage, TraceParam, template.ErrorPage, TraceTemplate, netMa, web, template.ErrorPage, TraceParam, template.ErrorPage )

				ioutil.WriteFile("./api/templates/render.go", []byte(renderFile), 0700 )

		if  template.Type == "package" {

			templateImport = fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package %s

				%s`, pk[len(pk) - 1], templateImport)

			ioutil.WriteFile("./exported_templates.go", []byte(templateImport), 0700)
		}


		//Methods have been added


		if template.Type == "faas" {

			log.Println("🔗 Saving file to ", fmt.Sprintf("%s%s%s", r, "/", template.Output))
			
			local_string += `
		
		func MakeFSHandle(root string){
			http.Handle("/dist/",  http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: root}))
		}
		`
			local_string = strings.Replace(local_string, "net_", "Net", -1)



			d1 := []byte(local_string)
			_ = ioutil.WriteFile(fmt.Sprintf("%s%s", r, template.Output), d1, 0700)

			
			if _, err := os.Stat(fmt.Sprintf("%s/src/func", os.ExpandEnv("$GOPATH"))); os.IsNotExist(err) {
				os.MkdirAll(fmt.Sprintf("%s/src/func", os.ExpandEnv("$GOPATH")), 0700)
			}

			os.Chdir(fmt.Sprintf("%s/src/func", os.ExpandEnv("$GOPATH")))
			if template.Gate == "" {
				template.Gate = "http://localhost:8080"
			}

			var sourcedep string
			for _, v := range template.Templates.Templates {
				tlc := fmt.Sprintf("%s", strings.ToLower(v.Name))
				os.RemoveAll(tlc)
				os.Mkdir(tlc, 0700)

				handlerTemp := fmt.Sprintf(`
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package function
// Handle a serverless request
import (
	app "%s/api/templates"
)

func Handle(req []byte) string {
	return app.Netb%s(string(req))
}
`, appname, v.Name)

				yamlTemp := fmt.Sprintf(`provider:
  name: faas
  gateway: %s

functions:
  %s:
    lang: go
    handler: ./%s
    image: %s

`, template.Gate, v.Name, tlc, tlc)

				_ = ioutil.WriteFile(fmt.Sprintf("%s/handler.go", tlc), []byte(handlerTemp), 0700)
				_ = ioutil.WriteFile(fmt.Sprintf("%s.yml", tlc), []byte(yamlTemp), 0700)

				os.Chdir(fmt.Sprintf("%s", tlc))
				chn := make(chan int)

				go DoSpin(chn)
				if sourcedep == "" {
					RunCmd("dep init")

					appath := strings.Replace(fmt.Sprintf("%s/src/%s/", os.ExpandEnv("$GOPATH"), appname), "//", "/", -1)
					vendpath := fmt.Sprintf("vendor/%s/", appname)
					os.RemoveAll(vendpath)
					os.MkdirAll(vendpath, 0700)
					CopyDir(appath, vendpath)
					if _, err := os.Stat("vendor/github.com/elazarl/go-bindata-assetfs/"); os.IsNotExist(err) {
						CopyDir(strings.Replace(fmt.Sprintf("%s/src/github.com/elazarl/go-bindata-assetfs/", os.ExpandEnv("$GOPATH")), "//", "/", -1), "vendor/github.com/elazarl/go-bindata-assetfs/")
					}
					os.RemoveAll("vendor/golang.org/x/tools")
					RunCmd(fmt.Sprintf("gofmt -s -w ../%s", tlc))
					sourcedep = fmt.Sprintf("%s/src/func/%s/", os.ExpandEnv("$GOPATH"), tlc)
				} else {
					CopyDir(fmt.Sprintf("%s/vendor/", sourcedep), "vendor/")
					CopyFile(fmt.Sprintf("%s/Gopkg.lock", sourcedep), "Gopkg.lock")
					CopyFile(fmt.Sprintf("%s/Gopkg.toml", sourcedep), "Gopkg.toml")
					RunCmd(fmt.Sprintf("gofmt -s -w ../%s", tlc))
				}
				chn <- 1
				close(chn)

				os.Chdir("../")
			}

			for _, v := range template.Endpoints.Endpoints {
				if v.Type != "f" {
					var tlc = strings.Replace(fmt.Sprintf("%s%s", v.Type, strings.Replace(strings.Title(strings.Replace(v.Path, "/", " ", -1)), " ", "", -1)), "star", "", -1)
					os.RemoveAll(tlc)
					os.Mkdir(tlc, 0700)

					handlerTemp := fmt.Sprintf(`
					// File generated by Gopher Sauce
					// DO NOT EDIT!!
					package function
// Handle a serverless request
import (
	app "%s/api/handlers"
	"bytes"
	"net/http"
	"net/http/httptest"
	"fmt"
)

func Handle(req []byte) string {

	readr := bytes.NewReader( req )
	request, err := http.NewRequest("%s", "%s", readr)
	if err != nil {
		return err.Error()
	 }
	rr := httptest.NewRecorder()
	handle := http.HandlerFunc(app.MakeHandler(app.Handler))
	handle.ServeHTTP(rr, request)
	rr.Flush()
	return fmt.Sprintf("%%s", rr.Body.String())
}
`, appname, v.Type, v.Path)

					yamlTemp := fmt.Sprintf(`provider:
  name: faas
  gateway: %s

functions:
  %s:
    lang: go
    handler: ./%s
    image: %s

`, template.Gate, tlc, tlc, strings.ToLower(tlc))

					_ = ioutil.WriteFile(fmt.Sprintf("%s/handler.go", tlc), []byte(handlerTemp), 0700)
					_ = ioutil.WriteFile(fmt.Sprintf("%s.yml", tlc), []byte(yamlTemp), 0700)

					os.Chdir(fmt.Sprintf("%s", tlc))
					chn := make(chan int)
					go DoSpin(chn)

					if sourcedep == "" {
						RunCmd("dep init")
						//cleanup bugs found with fmt
						appath := strings.Replace(fmt.Sprintf("%s/src/%s/", os.ExpandEnv("$GOPATH"), appname), "//", "/", -1)
						vendpath := fmt.Sprintf("vendor/%s/", appname)
						os.RemoveAll(vendpath)
						os.MkdirAll(vendpath, 0700)
						CopyDir(appath, vendpath)
						if _, err := os.Stat("vendor/github.com/elazarl/go-bindata-assetfs/"); os.IsNotExist(err) {
							CopyDir(strings.Replace(fmt.Sprintf("%s/src/github.com/elazarl/go-bindata-assetfs/", os.ExpandEnv("$GOPATH")), "//", "/", -1), "vendor/github.com/elazarl/go-bindata-assetfs/")
						}
						os.RemoveAll("vendor/golang.org/x/tools")
						RunCmd(fmt.Sprintf("gofmt -s -w ../%s", tlc))
						sourcedep = fmt.Sprintf("%s/src/func/%s/", os.ExpandEnv("$GOPATH"), tlc)
					} else {
						CopyDir(fmt.Sprintf("%s/vendor/", sourcedep), "vendor/")

						CopyFile(fmt.Sprintf("%s/Gopkg.lock", sourcedep), "Gopkg.lock")
						CopyFile(fmt.Sprintf("%s/Gopkg.toml", sourcedep), "Gopkg.toml")
						RunCmd(fmt.Sprintf("gofmt -s -w ../%s", tlc))
					}
					chn <- 1
					close(chn)

					os.Chdir("../")
				}
			}



			for _, v := range template.Templates.Templates {
				tlc := fmt.Sprintf("%s", strings.ToLower(v.Name))
				fmt.Printf("\r  \033[36mBuilding %s \033[m", tlc)
				chn := make(chan int)
				go DoSpin(chn)

				RunCmd(fmt.Sprintf("faas-cli build -f ./%s.yml", tlc))
				RunCmd(fmt.Sprintf("faas-cli deploy -f ./%s.yml", tlc))
				chn <- 1
				close(chn)
			}

			for _, v := range template.Endpoints.Endpoints {
				if v.Type != "f" {
					var tlc = strings.Replace(fmt.Sprintf("%s%s", v.Type, strings.Replace(strings.Title(strings.Replace(v.Path, "/", " ", -1)), " ", "", -1)), "star", "", -1)
					fmt.Printf("\r  \033[36mBuilding %s \033[m", tlc)
					chn := make(chan int)
					go DoSpin(chn)

					RunCmd(fmt.Sprintf("faas-cli build -f ./%s.yml", tlc))
					RunCmd(fmt.Sprintf("faas-cli deploy -f ./%s.yml", tlc))
					chn <- 1
					close(chn)
				}
			}

			os.Chdir(fmt.Sprintf("%s/src/%s", os.ExpandEnv("$GOPATH"), appname))

		} else if template.Type == "package" {

			log.Println("🔗 Saving file to ", fmt.Sprintf("%s%s%s", r, "/", template.Output))
			
			local_string += fmt.Sprintf(`
				func FileServer() http.Handler {
					return http.FileServer(&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, Prefix: "%s"})
				}`, web)

			local_string = strings.Replace(local_string, "net_", "Net", -1)
			d1 := []byte(local_string)

			_ = ioutil.WriteFile(fmt.Sprintf("%s%s", r, template.Output), d1, 0700)

		} else {

			local_string += `
			func main() {
				fmt.Fprintf(os.Stdout, "%v\n", os.Getpid())

				LaunchServer()
				`

			launcher := fmt.Sprintf(`
			// File generated by Gopher Sauce
			// DO NOT EDIT!!
			package main

			import "%s/api/globals"
			import "%s/api/methods"
			import "%s/api/templates"

			func LaunchServer(){
				%s
			}`, appname, appname, appname, template.Main)


			ioutil.WriteFile("launcher.go", []byte(launcher), 0700)

			if template.Prod {
				local_string += fmt.Sprintf(` sessionStore.Store.Options = &sessions.Options{
						    Path:     "/",
						    MaxAge:   86400 * 7,
						    HttpOnly: true,
						    Secure : true,
						    Domain : "%s",
						}`, template.Domain)

				//todo timeouts
			} else {
				os.Mkdir("api/tracer", 0700)

				appDashInit := `
				// File generated by Gopher Sauce
				// DO NOT EDIT!!
				package tracer

				import appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"

				func LoadTraceServer(){
					store := appdash.NewMemoryStore()

					// Listen on any available TCP port locally.
					l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
					if err != nil {
						log.Fatal(err)
					}
					collectorPort := l.Addr().(*net.TCPAddr).Port

					// Start an Appdash collection server that will listen for spans and
					// annotations and add them to the local collector (stored in-memory).
					cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
					go cs.Start()

					// Print the URL at which the web UI will be running.
					appdashPort := 8700
					appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
					appdashURL, err := url.Parse(appdashURLStr)
					if err != nil {
						log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
					}
					color.Red("✅ Important!")
					log.Println("To see your traces, go to ",  appdashURL )

					// Start the web UI in a separate goroutine.
					tapp, err := traceapp.New(nil, appdashURL)
					if err != nil {
						log.Fatal(err)
					}
					tapp.Store = store
					tapp.Queryer = store
					go func() {
						log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
					}()

					tracer := appdashot.NewTracer(appdash.NewRemoteCollector(fmt.Sprintf(":%d", collectorPort) ) )
					opentracing.InitGlobalTracer(tracer)
				}`

				local_string += "tracer.LoadTraceServer()"
				ioutil.WriteFile("api/tracer/appdash.go", []byte(appDashInit), 0700)


			}

			defhandler := ""
			if template.Cert == "" && template.HKey == "" {
				log.Println("Set gos gxml attributes `https-key` & `https-cert` to enable HTTPS")
				defhandler = `errgos := h.ListenAndServe()
					if errgos != nil {
						log.Fatal(errgos)
					}`
			} else {
				defhandler = fmt.Sprintf(`
	errgos := h.ListenAndServeTLS("%s", "%s")
		if errgos != nil {
			log.Fatal("ListenAndServe: ", errgos)
	}`,template.Cert, template.HKey)
			}
			local_string += fmt.Sprintf(`
					 %s
					 port := ":%s"
						if envport := os.ExpandEnv("$PORT"); envport != "" {
							port = fmt.Sprintf(":%%s", envport)
						}
					 log.Printf("Listenning on Port %%v\n", port)
					

					//+++extendgxmlmain+++
					
					stop := make(chan os.Signal, 1)

					signal.Notify(stop, os.Interrupt)
					http.Handle("/dist/", http.FileServer(&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, Prefix: "%s"}))
					http.HandleFunc("/", handlers.MakeHandler(handlers.Handler))

					h := &http.Server{Addr: port}

					go func(){
						%s
					}()

					<-stop

					log.Println("\nShutting down the server...")

					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

					h.Shutdown(ctx)

					Shutdown()

					log.Println("Server gracefully stopped")

					}

					//+++extendgxmlroot+++`, timeline, template.Port, web, defhandler)


					templateShutdown := fmt.Sprintf(`
					// File generated by Gopher Sauce
					// DO NOT EDIT!!
					package main

					func Shutdown(){
						%s
					}`, template.Shutdown)

					ioutil.WriteFile("shutdown.go", []byte(templateShutdown), 0700)

			var port string
			if !template.Prod  {			
				port = fmt.Sprintf("%s\nEXPOSE 8700", template.Port)
			} else {
				
				port = template.Port
			}
			dockerfile := fmt.Sprintf(`FROM golang:1.8
ENV WEBAPP /go/src/server
RUN mkdir -p ${WEBAPP}
COPY . ${WEBAPP}
ENV PORT=%s 
WORKDIR ${WEBAPP}
RUN go install
RUN rm -rf ${WEBAPP} 
EXPOSE %s
CMD server
`, template.Port, port)

			_ = ioutil.WriteFile(fmt.Sprintf("%s%s", r, "Dockerfile"), []byte(dockerfile), 0700)
			log.Println("🔗 Saving file to ", fmt.Sprintf("%s%s%s", r, "/", template.Output))
		

			local_string = strings.Replace(local_string, "net_", "Net", -1)
			d1 := []byte(local_string)

			outputPath := fmt.Sprintf("%s%s", r, template.Output)

			_ = ioutil.WriteFile(outputPath, d1, 0700)

			
		}


		RunCmd(fmt.Sprintf("goimports -w %s",r ) )

		var logfull string
		for _, sh := range template.PostCommand {
			logfull, _ = RunCmdSmart(sh)
			log.Println(logfull)
		}

	} 

	return
}

func (d *gos) ensureBackwardsCompatible(){


	for _, m := range d.Methods.Methods {

	
		newName := fmt.Sprintf("methods.%s(", strings.Title(m.Name) )
		newNameInternal := fmt.Sprintf("%s(", strings.Title(m.Name) )
		currentName := fmt.Sprintf("Net%s(", m.Name)

		

		for mE, e := range d.Endpoints.Endpoints {

			e.Method = ReplaceLegacy(e.Method, "net_", "Net")
			e.Method = ReplaceLegacy(e.Method,  currentName, newName)


			d.Endpoints.Endpoints[mE] = e

		}

		for mI, mM := range d.Methods.Methods {

			mM.Method = ReplaceLegacy(mM.Method, "net_", "Net")
			mM.Method = ReplaceLegacy(mM.Method,  currentName, newNameInternal)
			d.Methods.Methods[mI] = mM

		}

		

		d.Main = ReplaceLegacy(d.Main, "net_", "Net")
		d.Main = ReplaceLegacy(d.Main, currentName, newName)

		d.Shutdown = ReplaceLegacy(d.Shutdown, "net_", "Net")
		d.Shutdown = ReplaceLegacy(d.Shutdown, currentName, newName)



	}

	for i, m := range d.Methods.Methods {
		d.Methods.Methods[i].Name = strings.Title(m.Name) 
	}


	for i,g := range d.Variables {

		g.ExplicitName = strings.TrimSpace(strings.Split(g.Name,"=")[0])



		newName := fmt.Sprintf("globals.%s", strings.Title(g.ExplicitName) )
		currentName := fmt.Sprintf("%s", g.ExplicitName)


		for mE, e := range d.Endpoints.Endpoints {

			e.Method = ReplaceLegacy(e.Method,  currentName, newName)
			d.Endpoints.Endpoints[mE] = e

		}

		for mI, mM := range d.Methods.Methods {

			mM.Method = ReplaceLegacy(mM.Method,  currentName, newName)
			d.Methods.Methods[mI] = mM

		}

		

		d.Main = ReplaceLegacy(d.Main, currentName, newName)
		d.Shutdown = ReplaceLegacy(d.Shutdown, currentName, newName)
		d.Variables[i].Name = strings.Title(g.Name)

		if strings.Contains(g.Name,"=") {
			varComponents := strings.Split(g.Name,"=")
			d.Variables[i].Name = fmt.Sprintf("%s = %s", strings.Title(varComponents[0]), varComponents[1])
		}
	}


	// handle templates
	for _, t := range d.Templates.Templates {


		newName := fmt.Sprintf("templates.%s(", strings.Title(t.Name) )
		currentName := fmt.Sprintf("Netb%s(", t.Name)
		currentName2 := fmt.Sprintf("b%s(", t.Name)

		for mE, e := range d.Endpoints.Endpoints {

			e.Method = ReplaceLegacy(e.Method, "net_", "Net")
			e.Method = ReplaceLegacy(e.Method,  currentName, newName)
			e.Method = ReplaceLegacy(e.Method,  currentName2, newName)
			d.Endpoints.Endpoints[mE] = e

		}

		for mI, mM := range d.Methods.Methods {

			mM.Method = ReplaceLegacy(mM.Method, "net_", "Net")
			mM.Method = ReplaceLegacy(mM.Method,  currentName, newName)
			mM.Method = ReplaceLegacy(mM.Method,  currentName2, newName)
			d.Methods.Methods[mI] = mM

		}

		d.Main = ReplaceLegacy(d.Main, currentName, newName)
		d.Shutdown = ReplaceLegacy(d.Shutdown, currentName, newName)



	}

	// handle structs
	for i, s := range d.Header.Structs {

		original := s.Name
		s.Name = strings.Title(s.Name)

		newName := fmt.Sprintf("types.%s{", s.Name )
		newNameArr := fmt.Sprintf(" []types.%s", s.Name )

		returnName := fmt.Sprintf("types.%s", s.Name )

		assertType := fmt.Sprintf("(%s", s.Name )

		newAssertType := fmt.Sprintf("(types.%s", s.Name)

		currentName := fmt.Sprintf("%s{", original)
		currentNameArr := fmt.Sprintf("[]%s", original)

		d.Header.Structs[i] = s

		for mE, e := range d.Endpoints.Endpoints {

			e.Method = ReplaceLegacy(e.Method, "net_", "Net")
			e.Method = ReplaceLegacy(e.Method,  currentName, newName)
			e.Method = ReplaceLegacy(e.Method,  currentNameArr, newNameArr)

			e.Method = ReplaceLegacy(e.Method,  assertType, newAssertType)

			d.Endpoints.Endpoints[mE] = e

		}

		for mI, mM := range d.Methods.Methods {

			mM.Method = ReplaceLegacy(mM.Method, "net_", "Net")
			mM.Method = ReplaceLegacy(mM.Method,  assertType, newAssertType)
			mM.Method = ReplaceLegacy(mM.Method,  currentName, newName)
			mM.Method = ReplaceLegacy(mM.Method,  currentNameArr, newNameArr)
		
			mM.Returntype = ReplaceLegacy(mM.Returntype, s.Name, returnName)

			d.Methods.Methods[mI] = mM

		}



		d.Main = ReplaceLegacy(d.Main, assertType, newAssertType)

		d.Main = ReplaceLegacy(d.Main, currentName, newName)

		d.Shutdown = ReplaceLegacy(d.Shutdown, assertType, newAssertType)
		d.Shutdown = ReplaceLegacy(d.Shutdown, currentName, newName)

	}


}

func ReplaceLegacy(source, target, newExpr string) string {
	return strings.Replace(source, target, newExpr, -1)
}

func RunFile(root string, file string) {
	log.Println("∑ Running " + root + "/" + file)
	exe_cmd("go run " + root + "/" + file)
}

func RunCmd(cmd string) {
	exe_cmd(cmd)
}

func RunCmdString(cmd string) string {
	parts := strings.Fields(cmd)
	log.Println(cmd)
	var out *exec.Cmd
	if len(parts) == 5 {
		log.Println("Match")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4])
	} else if len(parts) == 9 {
		log.Println("Match GPGl")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8])
	} else if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0], parts[1])
	}

	var ou bytes.Buffer
	out.Stdout = &ou
	err := out.Run()
	if err != nil {
		log.Println("Error")
	}
	return ou.String()
}

func RunCmdByte(cmd string) []byte {
	parts := strings.Fields(cmd)
	log.Println(cmd)
	var out *exec.Cmd
	if len(parts) == 5 {
		log.Println("Match")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4])
	} else if len(parts) == 9 {
		log.Println("Match GPGl")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8])
	} else if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0], parts[1])
	}

	var ou bytes.Buffer
	out.Stdout = &ou
	err := out.Run()
	if err != nil {
		log.Println("Error")
	}
	return ou.Bytes()
}

func RunCmdSmartB(cmd string) ([]byte, error) {
	parts := strings.Fields(cmd)
	log.Println(parts[0], parts[1:])
	var out *exec.Cmd

	if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else {
		out = exec.Command(parts[0], parts[1:]...)
	}
	var ou, our bytes.Buffer
	out.Stdout = &ou
	out.Stderr = &our

	err := out.Run()
	log.Println(our.String())
	if err != nil {
		return our.Bytes(), err
	}
	return ou.Bytes(), nil
}

func RunCmdSmart(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	//log.Println(parts)
	var out *exec.Cmd

	if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else {
		out = exec.Command(parts[0], parts[1:]...)
	}

	var ou, our bytes.Buffer
	out.Stdout = &ou
	out.Stderr = &our

	//log.Println(BytesToString(our.Bytes()))
	err := out.Run()
	if err != nil {
		//	log.Println("%v", err.Error())
		return our.String(), err
	}
	return ou.String(), nil
}

func RunCmdSmarttwo(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	//	log.Println(parts[0],parts[1:])
	var out *exec.Cmd

	out = exec.Command(parts[0], parts[1:]...)

	var ou, our bytes.Buffer
	out.Stdout = &ou
	out.Stderr = &our

	log.Println(BytesToString(our.Bytes()))
	err := out.Run()
	if err != nil {
		//	log.Println("%v", err.Error())
		return our.String(), err
	}
	return ou.String(), nil
}

func RunCmdSmartZ(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	//	log.Println(parts[0],parts[1:])
	var out *exec.Cmd

	out = exec.Command(parts[0], parts[1])

	var ou, our bytes.Buffer
	out.Stdout = &ou
	out.Stderr = &our

	log.Println(BytesToString(our.Bytes()))
	err := out.Run()
	if err != nil {
		//	log.Println("%v", err.Error())
		return ou.String() + our.String(), err
	}
	return ou.String(), nil
}

func RunCmdSmartP(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	//	log.Println(parts[0],parts[1:])
	var out *exec.Cmd

	out = exec.Command(parts[0], parts[1], parts[2], "-benchmem")

	var ou, our bytes.Buffer
	out.Stdout = &ou
	out.Stderr = &our

	log.Println(BytesToString(our.Bytes()))
	err := out.Run()
	if err != nil {
		//	log.Println("%v", err.Error())
		return ou.String() + our.String(), err
	}
	return ou.String(), nil
}

func RunCmdSmartCmb(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	log.Println(parts[0], parts[1:])
	var out *exec.Cmd

	out = exec.Command(parts[0], parts[1:]...)

	ou, err := out.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return string(ou), nil
}

func BytesToString(b []byte) string {
	return string(b)
}

func RunCmdB(cmd string) {
	parts := strings.Fields(cmd)
	log.Println(cmd)
	var out *exec.Cmd
	if len(parts) == 5 {
		log.Println("Match")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4])
	} else if len(parts) == 9 {
		log.Println("Match GPGl")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8])
	} else if len(parts) == 8 {
		log.Println("Match decrypt GPGl")
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], "pass")

	} else if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0], parts[1], "2>&1")
	}

	var ou bytes.Buffer
	out.Stdout = &ou
	err := out.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println(ou.String())
}

func RunCmdA(cm string) error {

	parts := strings.Fields(cm)
	log.Println(parts)
	var cmd *exec.Cmd
	log.Println("Match decrypt GPGl")
	cmd = exec.Command(parts[0], parts[1], parts[2], "--no-tty", parts[3], parts[4], parts[5], parts[6])
	inpipe, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Pipe : ", err)
	}
	io.WriteString(inpipe, "")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	log.Println("Result: " + out.String())
	return nil
}

func exe_cmd(cmd string) {
	defer func() {
		if n := recover(); n != nil {

			log.Println(n)
		}
	}()
	parts := strings.Fields(cmd)
	log.Println(cmd)
	var out *exec.Cmd
	if len(parts) == 5 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3], parts[4])
	} else if len(parts) == 4 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0], parts[1])
	}
	stdoutStderr, err := out.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("%s\n", stdoutStderr)
}

func Exe_Stall(cmd string, chn chan bool) {
	//log.Println(cmd)
	parts := strings.Fields(cmd)
	var out *exec.Cmd
	if len(parts) > 3 {
		out = exec.Command(parts[0], parts[1], parts[2], parts[3])
	} else if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2], "2>&1")
	} else if len(parts) == 1 {
		out = exec.Command(parts[0], "2>&1")
	} else {
		out = exec.Command(parts[0], parts[1], "2>&1")
	}
	//stdout, err := out.StdoutPipe()

	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	out.Start()
	//r := bufio.NewReader(stdout)
	<-chn
	log.Println("💣 Killing proc.")
	out.Process.Kill()
	/* _, err = RunCmdSmart("kill -3 " + strconv.Itoa(out.Process.Pid) )

	if err != nil {
		log.Fatal(err)
	}
	*/
	runtime.Goexit()
}

func Exe_Stalll(cmd string) {
	log.Println(cmd)
	parts := strings.Fields(cmd)
	var out *exec.Cmd
	if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2], "2>&1")
	} else if len(parts) == 1 {
		out = exec.Command(parts[0], "2>&1")
	} else {
		out = exec.Command(parts[0], parts[1], "2>&1")
	}
	stdout, err := out.StdoutPipe()

	if err != nil {
		log.Println("error occurred")
		log.Printf("%s", err)
	}
	out.Start()
	r := bufio.NewReader(stdout)

	t := false
	for !t {
		line, _, _ := r.ReadLine()
		if string(line) != "" {
			log.Println(string(line))
		}

	}

}

func Exe_BG(cmd string) {
	log.Println(cmd)
	parts := strings.Fields(cmd)
	var out *exec.Cmd
	if len(parts) > 2 {
		out = exec.Command(parts[0], parts[1], parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0], parts[1])
	}
	stdout, err := out.StdoutPipe()
	if err != nil {
		log.Println("error occurred")
		log.Printf("%s", err)
	}
	out.Start()
	r := bufio.NewReader(stdout)
	t := false
	for !t {
		line, _, _ := r.ReadLine()
		if string(line) != "" {
			log.Println(string(line))
		}
	}
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func (d *gos) findStruct(name string) Struct {
	for _, imp := range d.Header.Structs {
		if imp.Name == name {
			return imp
		}
	}
	return Struct{Name: "000"}
}

func (d *gos) findMethod(name string) Method {
	for _, imp := range d.Methods.Methods {
		if imp.Name == name {
			return imp
		}
	}
	return Method{Name: "000"}
}

func (d *gos) PSaveGos(path string) {
	b, _ := xml.Marshal(d)
	ioutil.WriteFile(path, b, 0644)
}

func (d *gos) Delete(typ, id string) {
	if typ == "var" {
		temp := []GlobalVariables{}
		for _, v := range d.Variables {
			if v.Name != id {
				temp = append(temp, v)
			}
		}
		d.Variables = temp
	} else if typ == "import" {
		temp := []Import{}
		for _, v := range d.RootImports {
			if v.Src != id {
				temp = append(temp, v)
			}
		}
		d.RootImports = temp
	} else if typ == "timer" {

		temp := []Timer{}
		for _, v := range d.Timers.Timers {
			if v.Name != id {
				temp = append(temp, v)
			}
		}
		d.Timers.Timers = temp

	} else if typ == "end" {
		temp := []Endpoint{}
		for _, v := range d.Endpoints.Endpoints {
			if v.Path != id {
				temp = append(temp, v)
			}
		}
		d.Endpoints.Endpoints = temp
	} else if typ == "template" {

		temp := []Template{}
		for _, v := range d.Templates.Templates {
			if v.Name != id {
				temp = append(temp, v)
			}
		}
		d.Templates.Templates = temp

	} else if typ == "bundle" {

		temp := []Template{}
		for _, v := range d.Templates.Templates {
			if v.Bundle != id {
				temp = append(temp, v)
			}
		}
		d.Templates.Templates = temp

	}

}

func CreateVGos(path string) *VGos {
	v := &VGos{}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	d := xml.NewDecoder(bytes.NewReader(body))
	d.Entity = map[string]string{
		"&": "&",
	}
	err = d.Decode(&v)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	return v
}

func (d *gos) MStructs(ne []Struct) {
	d.Header.Structs = ne
}

func (d *gos) MObjects(ne []Object) {
	d.Header.Objects = ne
}

func (d *gos) MMethod(ne []Method) {
	d.Methods.Methods = ne
}

func PLoadGos(pathraw string) (*gos, error) {
	path := strings.Replace(pathraw, "\\", "/", -1)
	log.Println("🎁 loading " + path)
	v := &gos{}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//obj := Error{}
	//log.Println(obj);
	body = EscpaseGXML(body)
	d := xml.NewDecoder(bytes.NewReader(body))
	d.Strict = false
	err = d.Decode(&v)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, nil
	}

	for _, imp := range v.RootImports {
		//log.Println(imp.Src)
		if strings.Contains(imp.Src, ".gxml") {
			v.MergeWithV(os.ExpandEnv("$GOPATH") + "/" + strings.Trim(imp.Src, "/"))
			//copy files
		}
	}

	return v, nil
}

func (d *gos) Add(sec, typ, name string) {
	if sec == "var" {
		d.Variables = append(d.Variables, GlobalVariables{Name: name, Type: typ})
	} else if sec == "import" {
		d.RootImports = append(d.RootImports, Import{Src: name})
	} else if sec == "end" {
		d.Endpoints.Endpoints = append(d.Endpoints.Endpoints, Endpoint{Path: name, Id: NewID(15)})
	} else if sec == "timer" {
		d.Timers.Timers = append(d.Timers.Timers, Timer{Name: name})
	}
}

func (d *gos) AddS(sec string, typ interface{}) {
	if sec == "template" {
		d.Templates.Templates = append(d.Templates.Templates, typ.(Template))
	} else if sec == "import" {
		//d.RootImports = append(d.RootImports,Import{Src: name})
	}
}

func Decode64(decBuf, enc []byte) []byte {
	e64 := base64.StdEncoding
	maxDecLen := e64.DecodedLen(len(enc))
	if decBuf == nil || len(decBuf) < maxDecLen {
		decBuf = make([]byte, maxDecLen)
	}
	n, err := e64.Decode(decBuf, enc)
	_ = err
	return decBuf[0:n]
}

func (d *gos) UpdateMethod(id string, data string) {

	temp := []Endpoint{}
	for _, v := range d.Endpoints.Endpoints {
		if id == v.Id {
			v.Method = data
			temp = append(temp, v)
		} else {
			temp = append(temp, v)
		}
	}
	d.Endpoints.Endpoints = temp

}

func (d *gos) Update(sec, id string, update interface{}) {
	if sec == "var" {
		temp := []GlobalVariables{}
		for _, v := range d.Variables {
			if id == v.Name {
				temp = append(temp, update.(GlobalVariables))
			} else {
				temp = append(temp, v)
			}
		}
		d.Variables = temp
	} else if sec == "import" {
		//d.RootImports = append(d.RootImports,Import{Src: name})
		temp := []Import{}
		for _, v := range d.RootImports {
			if id == v.Src {
				temp = append(temp, update.(Import))
			} else {
				temp = append(temp, v)
			}
		}
		d.RootImports = temp
	} else if sec == "template" {
		temp := []Template{}
		for _, v := range d.Templates.Templates {
			if id == v.Name {
				v.Struct = update.(string)
				temp = append(temp, v)
			} else {
				temp = append(temp, v)
			}
		}
		d.Templates.Templates = temp
	} else if sec == "timer" {
		temp := []Timer{}
		for _, v := range d.Timers.Timers {
			if id == v.Name {

				temp = append(temp, update.(Timer))
			} else {
				temp = append(temp, v)
			}
		}
		d.Timers.Timers = temp
	} else if sec == "end" {
		temp := []Endpoint{}
		for _, v := range d.Endpoints.Endpoints {
			if id == v.Id {
				upd := update.(Endpoint)
				upd.Method = v.Method
				upd.Id = id
				temp = append(temp, upd)
			} else {
				temp = append(temp, v)
			}
		}
		d.Endpoints.Endpoints = temp
	}
}

func (d *gos) Set(attr, value string) {
	if attr == "app" {
		d.Type = value
	} else if attr == "port" {
		d.Port = value
	} else if attr == "key" {
		d.Key = value
	} else if attr == "erpage" {
		d.ErrorPage = value
	} else if attr == "fpage" {
		d.NPage = value
	} else if attr == "domain" {
		d.Domain = value
	}
}

func VLoadGos(pathraw string) (gos, error) {
	path := strings.Replace(pathraw, "\\", "/", -1)
	log.Println("🎁 loading " + path)
	v := gos{}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return v, err
	}

	//obj := Error{}
	//log.Println(obj);
	body = EscpaseGXML(body)
	d := xml.NewDecoder(bytes.NewReader(body))
	d.Strict = false
	err = d.Decode(&v)
	if err != nil {
		log.Printf("✊ error: %v", err)
		return v, nil
	}

	if v.FolderRoot == "" {
		ab := strings.Split(path, "/")
		v.FolderRoot = strings.Join(ab[:len(ab)-1], "/") + "/"
	}
	//process mergs
	for _, imp := range v.RootImports {
		//log.Println(imp.Src)
		if strings.Contains(imp.Src, ".gxml") {
			srcP := strings.Split(imp.Src, "/")
			dir := strings.Join(srcP[:len(srcP)-1], "/")
			if _, err := os.Stat(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + dir); os.IsNotExist(err) && strings.Contains(imp.Src, ".") {
				// path/to/whatever does not exist
				//log.Println("")
				color.Red("Package not found")
				log.Println("🎁 Downloading Package " + imp.Src)
				logg, err := RunCmdSmart("go get -v " + dir)
				if err != nil {
					log.Println("Error :", logg)
				} else {
					log.Println("go get ", dir, "Ok :", logg)
				}
			}
		} else {
			dir := TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + strings.Trim(imp.Src, "/")
			if strings.Contains(imp.Src, "\"") {
				tpset := strings.Split(imp.Src, "\"")
				dir = strings.Replace(dir, imp.Src, tpset[1], -1)
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) && strings.Contains(imp.Src, ".") {
				// path/to/whatever does not exist
				//log.Println("")
				color.Red("Package not found")
				log.Println("🎁 Downloading Package " + imp.Src)
				logg, err := RunCmdSmart("go get -v " + imp.Src)
				if err != nil {
					log.Println("Error :", logg)
				} else {
					log.Println("go get ", imp.Src, "Ok :", logg)
				}

			}
		}

		if strings.Contains(imp.Src, ".gxml") {
			v.MergeWith(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + strings.Trim(imp.Src, "/"))
			//copy files
		}
		//
	}

	return v, nil
}

func EscpaseGXML(data []byte) []byte {
	lines := strings.Split(string(data), "\n")

	
	str := strings.Join(lines, "\n")
	str = strings.Replace(str, "<func ", "<method m=\"exp\"", -1)
	str = strings.Replace(str, "</func", "</method", -1)
	return []byte(str)
}

func LoadGos(pathraw string) (*gos, error) {
	path := strings.Replace(pathraw, "\\", "/", -1)

	log.Println("🎁 loading " + path)
	v := gos{}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//obj := Error{}
	//log.Println(obj);
	body = EscpaseGXML(body)
	d := xml.NewDecoder(bytes.NewReader(body))
	d.Strict = false
	err = d.Decode(&v)
	if err != nil {
		log.Printf("✊ error: %v", err)
		return nil, nil
	}

	if v.FolderRoot == "" {
		ab := strings.Split(path, "/")
		pkgPath := strings.Join(ab[:len(ab)-1], "/")
		v.FolderRoot = pkgPath + "/"
	}
	//process mergs
	for _, imp := range v.RootImports {
		//log.Println(imp.Src)
		if strings.Contains(imp.Src, ".gxml") {
			srcP := strings.Split(imp.Src, "/")
			dir := strings.Join(srcP[:len(srcP)-1], "/")
			dir = fmt.Sprintf("%s%s", dir, "/")
			if _, err := os.Stat(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + dir); os.IsNotExist(err) && strings.Contains(imp.Src, ".") {
				// path/to/whatever does not exist
				//log.Println("")

				color.Red("Package not found")
				log.Println("🎁 Downloading Package " + imp.Src)

				logg, err := RunCmdSmart("go get -v " + dir)
				if err != nil {
					log.Println("Error :", logg)
				} else {
					log.Println("go get ", dir, "Ok :", logg)
				}
			}
		} else {
			dir := TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + strings.Trim(imp.Src, "/")
			if strings.Contains(imp.Src, "\"") {
				tpset := strings.Split(imp.Src, "\"")
				dir = strings.Replace(dir, imp.Src, tpset[1], -1)

			}
			dir = fmt.Sprintf("%s%s", dir, "/")
			if _, err := os.Stat(dir); os.IsNotExist(err) && strings.Contains(imp.Src, ".") {
				// path/to/whatever does not exist
				//log.Println("")
				color.Red("Package not found")
				log.Println("🎁 Ddownloading Package " + imp.Src)
				logg, err := RunCmdSmart("go get -v " + imp.Src)
				if err != nil {
					log.Println("Error :", logg)
				} else {
					log.Println("go get ", imp.Src, "Ok :", logg)
				}

			}
		}

		if strings.Contains(imp.Src, ".gxml") {
			v.MergeWith(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + strings.Trim(imp.Src, "/"))
			//copy files
		}
		//
	}

	v.ensureBackwardsCompatible()

	return &v, nil
}

func (d *gos) MergeWithV(target string) {
	log.Println("∑ Merging " + target)
	imp, err := LoadGos(target)
	if err != nil {
		log.Println(err)
	} else {

		for _, im := range imp.RootImports {
			if strings.Contains(im.Src, ".gxml") {
				imp.MergeWithV(os.ExpandEnv("$GOPATH") + "/" + strings.Trim(im.Src, "/"))
				//copy files
			} else {
				d.RootImports = append(d.RootImports, im)
			}
		}

		if imp.FolderRoot == "" {
			ab := strings.Split(target, "/")
			imp.FolderRoot = strings.Join(ab[:len(ab)-1], "/") + "/"
		}
		//d.RootImports = append(imp.RootImports,d.RootImports...)
		d.Header.Structs = append(imp.Header.Structs, d.Header.Structs...)
		d.Header.Objects = append(imp.Header.Objects, d.Header.Objects...)
		d.Methods.Methods = append(imp.Methods.Methods, d.Methods.Methods...)
		d.PostCommand = append(imp.PostCommand, d.PostCommand...)
		d.Timers.Timers = append(imp.Timers.Timers, d.Timers.Timers...)
		//Specialize method for templates
		//d.Variables = append(imp.Variables, d.Variables...)
		if imp.Package != "" && imp.Type == "package" {
			log.Println("Parsing Prefixes for " + imp.Package)
			for _, im := range imp.Templates.Templates {
				im.TemplateFile = imp.Package + "/" + im.TemplateFile
				d.Templates.Templates = append(d.Templates.Templates, im)
			}
		} else {
			d.Templates.Templates = append(imp.Templates.Templates, d.Templates.Templates...)
		}
		//copy files
		d.Endpoints.Endpoints = append(imp.Endpoints.Endpoints, d.Endpoints.Endpoints...)
	}

	d.Init_Func = d.Init_Func + ` 
	` + imp.Init_Func
}

func (d *gos) MergeWith(target string) {
	log.Println("∑ Merging " + target)
	imp, err := LoadGos(target)
	if err != nil {
		log.Println(err)
	} else {

		for _, im := range imp.RootImports {

			if strings.Contains(im.Src, ".gxml") {
				srcP := strings.Split(im.Src, "/")
				dir := strings.Join(srcP[:len(srcP)-1], "/")
				if _, err := os.Stat(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + dir); os.IsNotExist(err) && strings.Contains(im.Src, ".") {
					// path/to/whatever does not exist
					//log.Println("")

					color.Red("Package not found")
					log.Println("🎁 Downloading Package " + im.Src)

					logg, err := RunCmdSmart("go get -v " + dir)
					if err != nil {
						log.Println("Error :", logg)
					} else {
						log.Println("go get ", dir, "Ok :", logg)
					}
				}
			}
			if strings.Contains(im.Src, ".gxml") {
				imp.MergeWith(TrimSuffix(os.ExpandEnv("$GOPATH"), "/") + "/src/" + strings.Trim(im.Src, "/"))
				//copy files
			} else {
				d.RootImports = append(d.RootImports, im)
			}
		}


		

		srcP := strings.Split(target, "/")
		dir := strings.Join(srcP[:len(srcP)-1], "/")

		if _, err := os.Stat(filepath.Join( dir , "entry.go" ) ); os.IsNotExist(err) && strings.Contains(target, ".gxml") {

			d.Methods.Methods = append(imp.Methods.Methods, d.Methods.Methods...)
			d.Header.Structs = append(imp.Header.Structs, d.Header.Structs...)
			d.Header.Objects = append(imp.Header.Objects, d.Header.Objects...)

			if imp.Package != "" && imp.Type == "package" {
				log.Println("Parsing Prefixes for " + imp.Package)
				for _, im := range imp.Templates.Templates {
					im.TemplateFile = imp.Package + "/" + im.TemplateFile
					d.Templates.Templates = append(d.Templates.Templates, im)
				}
			} else {
				d.Templates.Templates = append(imp.Templates.Templates, d.Templates.Templates...)
			}

			if imp.Tmpl == "" {
				imp.Tmpl = "tmpl"
			}

			if imp.Web == "" {
				imp.Web = "web"
			}

			if d.Tmpl == "" {
				d.Tmpl = "tmpl"
			}

			if d.Web == "" {
				d.Web = "web"
			}

			os.MkdirAll(d.Tmpl+"/"+imp.Package, 0777)
			os.MkdirAll(d.Web+"/"+imp.Package, 0777)
			CopyDir(imp.FolderRoot+imp.Tmpl, d.Tmpl+"/"+imp.Package)
			CopyDir(imp.FolderRoot+imp.Web, d.Web+"/"+imp.Package)

		} else {
			targ := strings.Replace(target,"\\","/", -1)

			path := strings.Replace(targ, filepath.Join(os.ExpandEnv("$GOPATH"), "src") + "/", "" ,-1)
			path = strings.Replace(path, "/gos.gxml", "", -1)
	
			d.Packages = append(d.Packages , Package{ Name : imp.Package , Path : path} )
		}

		d.PostCommand = append(imp.PostCommand, d.PostCommand...)
		d.Timers.Timers = append(imp.Timers.Timers, d.Timers.Timers...)
		//Specialize method for templates
		d.Variables = append(imp.Variables, d.Variables...)

		
		d.Init_Func = d.Init_Func + ` 
	` + imp.Init_Func
		d.Main = d.Main + ` 
	` + imp.Main
		

		//copy files
		d.Endpoints.Endpoints = append(imp.Endpoints.Endpoints, d.Endpoints.Endpoints...)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func DoubleInput(p1 string, p2 string) (r1 string, r2 string) {
	log.Println(p1)
	fmt.Scanln(&r1)
	log.Println(p2)
	fmt.Scanln(&r2)
	return
}

func AskForConfirmation() bool {
	var response string
	log.Println("Please type yes or no and then press enter:")
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		log.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func getPy() string {
	return ``
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
