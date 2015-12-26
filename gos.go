package main

import (
	"github.com/cheikhshift/gos/core"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	//"time"
	"unicode"
)

var webroot string
var template_root string
var gos_root string
var GOHOME string


func LowerInitial(str string) string {
    for i, v := range str {
        return string(unicode.ToLower(v)) + str[i+1:]
    }
    return ""
  }

  func UpperInitial(str string) string {
    for i, v := range str {
        return string(unicode.ToUpper(v)) + str[i+1:]
    }
    return ""
  }

func prepBindForMobile(body []byte,pkg string) []byte {
	data := string(body)
	finds := []string{"AssetDir","AssetInfo","AssetNames"}

	for _,e := range finds {
		data = strings.Replace(data,e,LowerInitial(e), -1)		
	}

	data = strings.Replace(data,"package main","package " + pkg, -1)

	return []byte(data)
}

func writeLocalProtocol(pack string){
	cObjFile := `//
					//  FlowProtocol.m
					//  FlowCode
					//
					//  Created by Cheikh Seck on 4/2/15.
					//  Copyright (c) 2015 Gopher Sauce LLC. All rights reserved.
					//

					#import "FlowProtocol.h"
					#import "` + UpperInitial(pack) + `/` + UpperInitial(pack)  +`.h"

					@implementation FlowProtocol



					+ (BOOL)canInitWithRequest:(NSURLRequest*)theRequest
					{
					    if ([theRequest.URL.host caseInsensitiveCompare:@"localhost"] == NSOrderedSame) {
					        return YES;
					    }
					    return NO;
					}

					+ (NSURLRequest*)canonicalRequestForRequest:(NSURLRequest*)theRequest
					{
					    return theRequest;
					}

					- (void)startLoading
					{
					  
					    NSString *process = [self.request.URL.absoluteString stringByReplacingOccurrencesOfString:@"http://localhost" withString:@""];
					    //check here
					    NSString *GetString;
					   //NSLog(@"%@", self.request.HTTPBody );
					    if([process rangeOfString:@"?"].location != NSNotFound){
					        if([process componentsSeparatedByString:@"?"].count > 1 )
					        GetString = [[process componentsSeparatedByString:@"?"] objectAtIndex:1];
					        process = [[process componentsSeparatedByString:@"?"] objectAtIndex:0];
					    }


                        if([self.request HTTPBody] != nil){
                            GetString = [GetString stringByAppendingString:@"&"];
                            GetString = [GetString stringByAppendingString:[NSString stringWithUTF8String:[self.request.HTTPBody bytes] ]];
                        }
					    
					    CFStringRef fileExtension = (__bridge CFStringRef)[process pathExtension];
					    CFStringRef UTI = UTTypeCreatePreferredIdentifierForTag(kUTTagClassFilenameExtension, fileExtension, NULL);
					    CFStringRef MIMEType = UTTypeCopyPreferredTagWithClass(UTI, kUTTagClassMIMEType);
					    CFRelease(UTI);
					    NSString *MIMETypeString = (__bridge_transfer NSString *)MIMEType;
					    NSURLResponse *response = [[NSURLResponse alloc] initWithURL:[self.request URL]
					                                                        MIMEType:MIMETypeString
					                                           expectedContentLength:-1
					                                                textEncodingName:nil];
					    
					      dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
					          
					    //NSLog(@"%@", self.request.HTTPBody );
					   
					          
					  
					    [[self client] URLProtocol:self didReceiveResponse:response cacheStoragePolicy:NSURLCacheStorageNotAllowed];
					    [[self client] URLProtocol:self didLoadData:Go` + UpperInitial(pack) +`LoadUrl(process, [self parseParams:GetString], self.request.HTTPMethod)];
					    [[self client] URLProtocolDidFinishLoading:self];
					      });
					   
					}

					- (NSData *) parseParams: (NSString *) input {
					    if(![input isEqualToString:@""]){
					    NSArray *pieces = [input componentsSeparatedByString:@"&"];
					    NSDictionary *payload = [NSMutableDictionary new];
					    
					    
					    
					    for (int i = 0; i < pieces.count; i++) {
					        NSString * param  = [pieces objectAtIndex:i];
					        if(![param isEqualToString:@""]){
					         
					            NSArray *keyset = [param componentsSeparatedByString:@"="];
					            [payload setValue:[self urlDecode:[keyset objectAtIndex:1] ] forKey:[self urlDecode:[keyset objectAtIndex:0]] ];
					            
					        }
					    }
					    NSError *error;
					    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:payload
					                                                       options:NSJSONWritingPrettyPrinted // Pass 0 if you don't care about the readability of the generated string
					                                                         error:&error];
					    
					    if (! jsonData) {
					        NSLog(@"Got an error: %@", error);
					        return nil;
					    } else {
					        NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
					        return [jsonString dataUsingEncoding:NSUTF8StringEncoding];
					    }
					    }
					    return nil;
					    
					}

					- (NSString *) urlDecode :(NSString *) input {
					    return [[input stringByReplacingOccurrencesOfString:@"+" withString:@" "]
					            stringByReplacingPercentEscapesUsingEncoding:NSUTF8StringEncoding];
					}
	

					- (void) stopLoading {
					    
					}

					@end
`

	ioutil.WriteFile(os.ExpandEnv("$GOPATH") + "/src/github.com/cheikhshift/gos/iosClasses/FlowProtocol.m",[]byte(cObjFile), 0644)
}

var gosTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<gos>
	<!--Stating the deployment type GoS should compile -->
	<!-- Curent valid types are webapp,shell and bind -->
	<!-- Shell = cli, sort of a GoS(Ghost) in the Shell -->
	<deploy>SETTHIS</deploy>
	<package>mymobile</package>
	
	<!-- Using import within different tags will have different results -->
	<!-- We going to make the goPkg Mongo Db Driver available to our application -->
	<!-- Using <import/> within the <go/> tag is similar to using the import call within a .go file -->
	<!-- To be less dramating, GoS will skip packages that it has already imported -->
	
	<!-- Go File output name -->
	<output>server_out.go</output>
	<!-- exported session fields available to Session -->


	<key>a very very very very secret key</key>
	<!-- Declare global variables -->
	<!-- Contains interfaces and structs
	 that will be used by the GoS application -->
	<header> 
			<!-- remember to Jumpline when stating methods or different struct attributes, it is vital for our parser \n trick -->
	</header>
	<methods>
		<!-- Vars are defined as usual except within the var attribute for example :  -->
		<!-- If there is a basic go function : func hackfmt(data string, data2 string) -->
		<!-- the attribute names would be called as such var="data string,data2 string" -->
		<!-- Similar to a go function decleration-->
		<!--  if a method matches the criteria for an  interface it will be used as an interface method -->
		<!-- To prevent that use the autoface attribute and set it to "false" By default it is true -->
		<!-- Use the keep-local="true" attribute to limit a method within a Go file -->	
		<!-- Sometimes your method will return data  -->
		<!-- And to do so we will need to add a return var list by using the return attribute  -->
		<!-- Sometimes the autointerface will reuse the wrong the function, or your interface methods need a bit more distinction -->
		<!-- Vis a  vis which object types are used in generating these mutating methods -->
		<!--Use the limit attribute to narrow down the applicable structs for this method -->
		<!-- Use the object attribute to determine the name of the local variable name to be mutated within the function. By default GoS will assume object is the variable name  -->
	</methods>

	<templates>
 		<!-- Template libraries are useful for expediting page creation and reuse common website elements within this GoS application -->
 		<!-- Templates are nested and customized with the template function instead of using the normal {{template "Name"}} call you can now use {{Button &{Color:"#fff"}& }}
 		{{Modal &{Color:"#fff"}& }}  -->
 		<!-- *Notice that special braces are used to initialize the parameters of the struct '&{' and '}&' -->
 		
 		<!-- <template name="Bootstrap_alert" tmpl="bootstrap/alert" struct="Bootstrap_alert" /> -->
 		
	</templates>
	<endpoints>
      <!-- Depending on your build type the usage of this tag will vary. -->
      <!-- For WebServers it will override any request for a given path and run the specified method. No vars or return types are needed for  -->
      <!-- methods linked to an API call, please keep in mind that you may use w for http.ResponseWriter and r for http.Request . Additional available function variables is params and session. If a function is api listed it will not be used anywhere else.-->
      <!-- <end /> is the endpoint tag and has the variables path,method, -->
      <!-- Happy trails!!! -->
      <!-- <end path="/index/api" method="login" type="POST" ></end> -->
	</endpoints>
</gos>
` 


func main() {
	GOHOME = os.ExpandEnv("$GOPATH") + "/src/"
    	//fmt.Println( os.Args)
    if len(os.Args) > 1 {
    //args := os.Args[1:]
    		if os.Args[1] == "dependencies" {
    			fmt.Println("∑ Getting GoS dependencies")
    			core.RunCmd("go get -u github.com/jteeuwen/go-bindata/...")
    			core.RunCmd("go get github.com/gorilla/sessions")
    			core.RunCmd("go get github.com/elazarl/go-bindata-assetfs")
    			core.RunCmd("go get github.com/kronenthaler/mod-pbxproj")
    			fmt.Println("ChDir " + os.ExpandEnv("$GOPATH") + "/src/github.com/kronenthaler/mod-pbxproj")
    			os.Chdir(os.ExpandEnv("$GOPATH") + "/src/github.com/kronenthaler/mod-pbxproj")
    			core.RunCmd("python setup.py install" )
    			//time.Sleep(time.Second *120)
    			fmt.Println("Done")
    			return
    		}

    		if os.Args[1] == "make" {
    		//2 is project folder
    		
    		    os.MkdirAll(os.ExpandEnv("$GOPATH") + "/src/" + strings.Trim(os.Args[2], "/") + "/web", 0777 )
    			os.MkdirAll(os.ExpandEnv("$GOPATH") + "/src/" + strings.Trim(os.Args[2], "/") + "/tmpl",0777 )
    			ioutil.WriteFile(os.ExpandEnv("$GOPATH") + "/src/" + strings.Trim(os.Args[2], "/") + "/gos.xml", []byte(gosTemplate), 0777)	
    			return
    		}
    
    		GOHOME = GOHOME   + strings.Trim(os.Args[2],"/")
    		serverconfig := os.Args[3]
    		webroot = os.Args[4]
    		template_root = os.Args[5]
    		fmt.Println("∑ GoS Speed compiler ");
    		coreTemplate,err := core.LoadGos( GOHOME + "/" + serverconfig ); 
			if err != nil {
				fmt.Println(err)
				return 
			}

			//fmt.Println(coreTemplate.Methods.Methods)
			coreTemplate.WriteOut = false

			//always delete add on folders prior
		
			core.Process(coreTemplate,GOHOME, webroot,template_root);

			if coreTemplate.Type == "webapp" {


					if os.Args[1] == "run" {
						os.Chdir(GOHOME)
						fmt.Println("Invoking go-bindata");
						core.RunCmd(os.ExpandEnv("$GOPATH") + "/bin/go-bindata -debug " + webroot +"/... " + template_root + "/...")
						//time.Sleep(time.Second*100 )
						//core.RunFile(GOHOME, coreTemplate.Output)
						core.RunCmd("go build")
						pk := strings.Split(strings.Trim(os.Args[2],"/"), "/")
						fmt.Println("Use Ctrl + C to quit")
						core.Exe_Stall("./" + pk[len(pk) - 1] )
					}

					if os.Args[1] == "export" {
						fmt.Println("Generating Export Program")
						os.Chdir(GOHOME)		
						//create both zips
						fmt.Println("Invoking go-bindata");
						core.RunCmd(  os.ExpandEnv("$GOPATH") + "/bin/go-bindata  " + webroot +"/... " + template_root + "/...")
						core.RunCmd("go build")
					}
			} else if coreTemplate.Type == "bind" {

				//check for directory gomobile
				if os.Args[1] == "export" {
						fmt.Println("Generating Export Program")
						os.Chdir(GOHOME)		
						//create both zips
						 fmt.Println("Invoking go-bindata");
						 core.RunCmd( os.ExpandEnv("$GOPATH") + `/bin/go-bindata `  + webroot +"/... " + template_root + "/...")
						 body,er := ioutil.ReadFile(GOHOME + "/bindata.go")
						 if er != nil {
						 	fmt.Println(er)
						 	return
						 }
						 writeLocalProtocol(coreTemplate.Package)
						 fmt.Println("Preparing Bindata for framework conversion...")
						 ioutil.WriteFile("bindata.go", prepBindForMobile(body, coreTemplate.Package)  ,0644)
						 core.RunCmd( os.ExpandEnv("$GOPATH")  + "/bin/gomobile bind -target=ios")
						 //edit plist file
						 subp := "/sub.check"
						 _,error := ioutil.ReadFile(subp)	
						 if error != nil {
						 ioutil.WriteFile(subp,[]byte("StubCompletion"),0644)
						 pathSp := strings.Split(os.Args[6],"/")
						 finalSub := ""
						 if len(pathSp) > 1 {
						 	finalSub = pathSp[len(pathSp) - 1]
						 } else {
						 	finalSub = os.Args[6]
						 }
						 plistPath := os.Args[6] + "/" + finalSub + "/Info.plist"
						 plist,erro := ioutil.ReadFile(plistPath)
						 if erro != nil {
						 	fmt.Println("Please check your project's folder for the Info.plit file")
						 	return
						 }

						 ioutil.WriteFile(plistPath, []byte(strings.Replace(string(plist), `<key>UIMainStoryboardFile</key>
	<string>Main</string>`,``,-1)),0644 )

						 core.RunCmd("python " + os.ExpandEnv("$GOPATH") + "/src/github.com/cheikhshift/gos/core/addFlow.py " + strings.Trim(os.Args[2],"/") +" " + os.Args[6] + " " + UpperInitial(coreTemplate.Package))
						 //if project does not exist create it and link this framework

						} else {
							fmt.Println("Subexists no need for Linkage :o")
						}
					}

			}


    	

	} else { 
	
    fmt.Println("∑ Welcome to Gos v1.0")
	fmt.Println("To begin please tell us a bit about the gos project you wish to compile");
	fmt.Printf("We need the GoS package folder relative to your $GOPATH/src (%v)\n", GOHOME)
   	gosProject := "" 
   	serverconfig := ""

   	fmt.Scanln(&gosProject)
   	GOHOME = GOHOME  + strings.Trim(gosProject,"/")
   	fmt.Printf("We need your Gos Project config source (%v)\n", GOHOME)
   	fmt.Scanln(&serverconfig)
    //fmt.Println(GOHOME)
	webroot,template_root = core.DoubleInput("What is the name of your webroot's folder ?", "What is the name of your template folder? ") 
		fmt.Println("Are you ready to begin? ");
		if core.AskForConfirmation() {
			fmt.Println("ΩΩ Operation Started!!");
			coreTemplate,err := core.LoadGos( GOHOME + "/" + serverconfig ); 
			if err != nil {
				fmt.Println(err)
				return 
			}

			coreTemplate.WriteOut = false
			core.Process(coreTemplate,GOHOME, webroot,template_root);
			fmt.Println("One moment...")
			core.RunCmd("go get -u github.com/jteeuwen/go-bindata/...")
    	    core.RunCmd("go get github.com/gorilla/sessions")
    		core.RunCmd("go get github.com/elazarl/go-bindata-assetfs")
			fmt.Println("Would you like to just run this application [y,n]")

			if core.AskForConfirmation() {
				os.Chdir(GOHOME)
				fmt.Println("Invoking go-bindata");
				core.RunCmd(os.ExpandEnv("$GOPATH") + "/bin/go-bindata -debug " + webroot +"/... " + template_root + "/...")
				//time.Sleep(time.Second*100 )
				//core.RunFile(GOHOME, coreTemplate.Output)
				core.RunCmd("go build")
				pk := strings.Split(strings.Trim(gosProject,"/"), "/")
				fmt.Println("Use Ctrl + C to quit")
				core.Exe_Stall("./" + pk[len(pk) - 1] )

			} else {
				fmt.Println("Is this a Mobile application [y,n]")

				if !core.AskForConfirmation() {
				fmt.Println("Would you like to create an export release [y,n]")

				if core.AskForConfirmation() {
					fmt.Println("Generating Export Program")
					os.Chdir(GOHOME)		
					//create both zips
					fmt.Println("Invoking go-bindata");
					core.RunCmd(  os.ExpandEnv("$GOPATH") + "/bin/go-bindata  " + webroot +"/... " + template_root + "/...")
					core.RunCmd("go build")
				
				}
				} else {
					//create mobile export here
					fmt.Println("Generating Export Program")
						os.Chdir(GOHOME)		
						//create both zips
						 fmt.Println("Invoking go-bindata");
						 core.RunCmd( os.ExpandEnv("$GOPATH") + `/bin/go-bindata `  + webroot +"/... " + template_root + "/...")
						 body,er := ioutil.ReadFile(GOHOME + "/bindata.go")
						 if er != nil {
						 	fmt.Println(er)
						 	return
						 }
						 fmt.Println("Preparing Bindata for framework conversion...")
						 ioutil.WriteFile("bindata.go", prepBindForMobile(body, coreTemplate.Package)  ,0644)
						 core.RunCmd( os.ExpandEnv("$GOPATH")  + "/bin/gomobile bind -target=ios")
						 //edit plist file
						 subp := "sub.check"

						 fmt.Println("What is the folder name of your IOS application?")
						 folderX := ""
						 fmt.Scanln(&folderX)
						 _,error := ioutil.ReadFile(subp)	
						 if error != nil {
						 ioutil.WriteFile(subp,[]byte("StubCompletion"),0644)
						 pathSp := strings.Split(folderX,"/")
						 finalSub := ""
						 if len(pathSp) > 1 {
						 	finalSub = pathSp[len(pathSp) - 1]
						 } else {
						 	finalSub = folderX
						 }
						 plistPath := folderX + "/" + finalSub + "/Info.plist"
						 plist,erro := ioutil.ReadFile(plistPath)
						 if erro != nil {
						 	fmt.Println("Please check your project's folder for the Info.plit file")
						 	return
						 }
						 writeLocalProtocol(coreTemplate.Package)

						 ioutil.WriteFile(plistPath, []byte(strings.Replace(string(plist), `<key>UIMainStoryboardFile</key>
	<string>Main</string>`,``,-1)),0644 )

						 core.RunCmd("python " + os.ExpandEnv("$GOPATH") + "/src/github.com/cheikhshift/gos/core/addFlow.py " + strings.Trim(gosProject,"/") +" " + folderX + " " + UpperInitial(coreTemplate.Package))
						 //if project does not exist create it and link this framework

						} else {
							fmt.Println("Subexists no need for Linkage :o")
						}
						fmt.Println("Your file is ready, go on your default IDE and run your application :)")

				}
			}
			

		} else {
			fmt.Println("Operation Cancelled!!")
		}
	}

}