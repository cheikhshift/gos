package main

import (
	"bufio"
	"fmt"
	"github.com/0xAX/notificator"
	"github.com/cheikhshift/gos/core"
	"github.com/fatih/color"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

var notify *notificator.Notificator
var webroot string
var templateroot string
var gos_root string
var appout string
var GOHOME string
var serverconfig string
var Type string

func LowerInitial(str string) (reslt string) {
	for _, v := range str {
		reslt = fmt.Sprintf("%s%s", string(unicode.ToLower(v)), str[1:])
		return
	}
	return
}

func UpperInitial(str string) (reslt string) {
	for _, v := range str {
		reslt = fmt.Sprintf("%s%s", string(unicode.ToUpper(v)), str[1:])
		return
	}
	return
}

func prepBindForMobile(body []byte, pkg string) []byte {
	data := string(body)
	finds := []string{"AssetDir", "AssetInfo", "AssetNames"}

	for _, e := range finds {
		data = strings.Replace(data, e, LowerInitial(e), -1)
	}

	data = strings.Replace(data, "package main", "package "+pkg, -1)

	return []byte(data)
}

func writeLocalProtocol(pack string) {
	cTissueHeader := `
			//
			//  FlowTissue.h
			//  GoTetst2
			//
			//  Created by OrendaCapital on 12/29/15.
			//  Copyright © 2015 Cheikh Seck LLC. All rights reserved.
			//

			#import <Foundation/Foundation.h>
			#import <AVFoundation/AVFoundation.h>
			#import <CoreLocation/CoreLocation.h>
			#import "` + UpperInitial(pack) + `/` + UpperInitial(pack) + `.h"
			#import "ViewController.h"
			#import "FlowThreadManager.h"


			@interface FlowTissue : NSObject  <Go` + UpperInitial(pack) + `Flow> {
			    
			}

			+ (void) handleRequest:(NSString *) endpoint;
			@end

	`

	cTissueClass := ` //
//  FlowTissue.m
//  GoTetst2
//
//  Created by OrendaCapital on 12/29/15.
//  Copyright © 2015 Cheikh Seck LLC. All rights reserved.
//

#import "FlowTissue.h"
#import "FlowBluetooth.h"
#import "FlowAccellerometer.h"



@implementation FlowTissue



- (void) trackMotion {
    [[UIAccelerometer sharedAccelerometer] setDelegate:[FlowThreadManager instance]];
    NSLog(@"Watching movements");

}

+ (void) handleRequest:(NSString *) endpoint {
	Go` + UpperInitial(pack) + `LoadUrl(endpoint, nil, @"GET",[FlowThreadManager tissue]);
}

- (void) stopMotion {
    [[UIAccelerometer sharedAccelerometer] setDelegate:nil];
}


- (void) notify:(NSString *)title message:(NSString *)message {
    UILocalNotification* localNotification = [[UILocalNotification alloc] init];
    localNotification.fireDate = [NSDate dateWithTimeIntervalSinceNow:0];
    localNotification.alertBody = message;
    localNotification.alertTitle = title;
    localNotification.timeZone = [NSTimeZone defaultTimeZone];
    [[UIApplication sharedApplication] scheduleLocalNotification:localNotification];
}

/*
    Flow Tissue Core Comm between Go and native langs to reach hardware specs
    Sound, touch scan, app links, GPS and files...
*/

- (int) device {
    return 0;
}

- (void) createPictureNamed:(NSString *)name {
    //take picture and save to specified name....
       dispatch_async(dispatch_get_main_queue(), ^{
           [FlowThreadManager takePicture:name];
       });
    
}

//sound
- (void) play:(NSString *)path {
    
    NSError *error = nil;
    FlowThreadManager *stream = [FlowThreadManager instance];
    NSData *fileData = [NSData dataWithContentsOfFile:[[FlowTissue applicationDocumentsDirectory] stringByAppendingString:path] ];
    
    if (stream.audioPlayer != nil) {
        if (stream.isPlaying){
            [stream.audioPlayer stop];
        }
    }
    
    stream.audioPlayer = [[AVAudioPlayer alloc] initWithData:fileData error:&error];
    
    [stream.audioPlayer prepareToPlay];
    [stream.audioPlayer play];
    if (error == nil)
    stream.isPlaying = YES;
    else stream.isPlaying = NO;
}

- (void) playFromWebRoot:(NSString *)path {
    NSError *error = nil;
    FlowThreadManager *stream = [FlowThreadManager instance];
    NSData *fileData = Go` + UpperInitial(pack) + `LoadUrl(path, nil, @"GET", nil);
    
    if (stream.audioPlayer != nil) {
        if (stream.isPlaying){
            [stream.audioPlayer stop];
        }
    }
    
    stream.audioPlayer = [[AVAudioPlayer alloc] initWithData:fileData error:&error];
    
    [stream.audioPlayer prepareToPlay];
    [stream.audioPlayer play];
    
    if (error == nil)
    stream.isPlaying = YES;
    else stream.isPlaying = NO;
    
}

- (void) setVolume:(int)power {
    FlowThreadManager *stream = [FlowThreadManager instance];
    [stream.audioPlayer setVolume: (float) (power/100) ];
}

- (int) getVolume {
    FlowThreadManager *stream = [FlowThreadManager instance];
    //[stream.audioPlayer setVolume: (float) (power/100) ];
    return 100*stream.audioPlayer.volume;
}

- (void) stop {
    FlowThreadManager *stream = [FlowThreadManager instance];
    stream.isPlaying = NO;
    [stream.audioPlayer stop];
}

- (BOOL) isPlaying {
    FlowThreadManager *stream = [FlowThreadManager instance];
    return stream.isPlaying;
}


//Applinks
- (void) openAppLink:(NSString *)url {
        //process applinkios
    dispatch_async(dispatch_get_main_queue(), ^{
    UIApplication *ourApplication = [UIApplication sharedApplication];
    NSString *URLEncodedText = [url stringByAddingPercentEscapesUsingEncoding:NSUTF8StringEncoding];
    NSString *ourPath =URLEncodedText;
    NSURL *ourURL = [NSURL URLWithString:ourPath];
    if ([ourApplication canOpenURL:ourURL]) {
        [ourApplication openURL:ourURL];
    }
    });
}

//GPS
- (void) requestLocation {
    //[[FlowThreadManager getGPS] requestWhenInUseAuthorization];
    //[[FlowThreadManager getGPS] requestLocation];
}

- (void) showLoad {
    dispatch_async(dispatch_get_main_queue(), ^{
    [FlowThreadManager loadScreen:YES usingMessage:@""];
    });
}

- (void) hideLoad {
    [FlowThreadManager loadScreen:NO usingMessage:@""];
}

- (void) runJS:(NSString *)line {
    dispatch_async(dispatch_get_main_queue(), ^{
    [FlowThreadManager runJS:line];
    });
}




//files
- (NSString *) absolutePath:(NSString *)file {
    return [[FlowTissue applicationDocumentsDirectory] stringByAppendingString:file];
}

- (BOOL) download:(NSString *)url target:(NSString *)target {
    
    //NSString *stringURL = @"http://www.somewhere.com/thefile.png";
    NSURL  *urll = [NSURL URLWithString:url];
    NSData *urlData = [NSData dataWithContentsOfURL:urll];
    if ( urlData )
    {
        NSString  *filePath = [self absolutePath:target];
        [urlData writeToFile:filePath atomically:YES];
        return YES;
    }
    
    return NO;
}

- (void) downloadLarge:(NSString *)url target:(NSString *)target {
    
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{

    //NSString *stringURL = @"http://www.somewhere.com/thefile.png";
    NSURL  *urll = [NSURL URLWithString:url];
    NSData *urlData = [NSData dataWithContentsOfURL:urll];
    if ( urlData )
    {
        NSString  *filePath = [self absolutePath:target];
        dispatch_async(dispatch_get_main_queue(), ^{
        [urlData writeToFile:filePath atomically:YES];
        });
       
    }
        
    });
   
}

- (NSString *) base64String:(NSString *)target {
    return [[self getBytes:target] base64EncodedStringWithOptions:0];
}

- (NSData *) getBytes:(NSString *)target {
    return [NSData dataWithContentsOfFile:[self absolutePath:target]];
}

- (NSData *) getBytesFromUrl:(NSString *)target {
    return [NSData dataWithContentsOfURL:[NSURL URLWithString:[self absolutePath:target]]];
}


- (BOOL) deleteDirectory:(NSString *)path {
    return [[NSFileManager defaultManager] removeItemAtPath:[self absolutePath:path] error:nil];

}

- (BOOL) deleteFile:(NSString *)path {
    return [self deleteDirectory:path];
}






+ (NSString *) applicationDocumentsDirectory
{
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *basePath = paths.firstObject;
    return basePath;
}


- (double) width {
    CGFloat width = [UIScreen mainScreen].bounds.size.width;
    return (double) width;
}

- (double) height {
    CGFloat height = [UIScreen mainScreen].bounds.size.height;
    return (double) height;
}


- (void)pushView:(NSString *)url {
        dispatch_async(dispatch_get_main_queue(), ^{
          [FlowThreadManager pulseView:url];
        });
    
    NSLog(@"Openning view %@", url);
}

- (void) dismissView {
    dispatch_async(dispatch_get_main_queue(), ^{
    UINavigationController *navcontroller = (UINavigationController *)[UIApplication sharedApplication].keyWindow.rootViewController;
    // Replace the current view controller
    NSMutableArray *viewControllers = [NSMutableArray arrayWithArray:[navcontroller viewControllers]];
    
    [viewControllers removeLastObject];
    
    [navcontroller setViewControllers:viewControllers animated:YES];
    });
}

- (void) dismissViewatInt:(int)index {
     dispatch_async(dispatch_get_main_queue(), ^{
    UINavigationController *navcontroller = (UINavigationController *)[UIApplication sharedApplication].keyWindow.rootViewController;
    // Replace the current view controller
    NSMutableArray *viewControllers = [NSMutableArray arrayWithArray:[navcontroller viewControllers]];
    
    [viewControllers removeObjectAtIndex:index];
    
    [navcontroller setViewControllers:viewControllers animated:YES];
         
    });
}


@end
`

	cObjFile := `//
					//  FlowProtocol.m
					//  FlowCode
					//
					//  Created by Cheikh Seck on 4/2/15.
					//  Copyright (c) 2015 Gopher Sauce LLC. All rights reserved.
					//

					#import "FlowProtocol.h"
					#import "FlowTissue.h"
					#import "` + UpperInitial(pack) + `/` + UpperInitial(pack) + `.h"

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


                        if([self.request HTTPBody] != nil && [self.request.HTTPBody length] > 0){
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
					   
					    [[self client] URLProtocol:self didLoadData:Go` + UpperInitial(pack) + `LoadUrl(process, [self parseParams:GetString], self.request.HTTPMethod,[FlowThreadManager tissue])];
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

	ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/github.com/cheikhshift/gos/iosClasses/FlowProtocol.m", []byte(cObjFile), 0644)
	ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/github.com/cheikhshift/gos/iosClasses/FlowTissue.h", []byte(cTissueHeader), 0644)
	ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/github.com/cheikhshift/gos/iosClasses/FlowTissue.m", []byte(cTissueClass), 0644)
}

var htmlTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    	<title>Blank page</title>
  </head>
  <body>
    <h1>Hello, world!</h1>

    <!-- jQuery first, then Tether, then Bootstrap JS. -->
    <script src="https://code.jquery.com/jquery-3.1.1.slim.min.js" integrity="sha384-A7FZj7v+d/sdmMqp/nOQwliLvUsJfDHW+k9Omg/a/EheAdgtzNs3hpfag6Ed950n" crossorigin="anonymous"></script>
   
  </body>
</html>`
var gosTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<gos>
	<!-- xml docs : http://golangserver.com/docs/markup.html -->
	<deploy>webapp</deploy>
	<port>8080</port>
	<package>if-package-is-library</package>
	<not_found>/your-404-page</not_found>
	<error>/your-500-page</error>

	<output>application.go</output>
	<domain></domain><!-- Cookie domain -->
	<main>	
	//psss go code here : func main()
	</main>


	<key>a very very very very secret key</key>
	
	<header> 
	</header>
	<methods>
		
	</methods>

	<templates>

	</templates>
	<endpoints>
  
	</endpoints>
</gos>
`

func GetLine(fname string, match string) int {
	intx := 0
	file, err := os.Open(fname)
	if err != nil {
		color.Red("Could not find a source file")
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		intx = intx + 1
		if strings.Contains(scanner.Text(), strings.TrimSpace(match)) {

			return intx
		}

	}

	return -1
}

func VmdOne() {

	log.Println(">>")
	var args_c string
	cmd_set := os.Args[2:]

	if len(cmd_set) < 2 {

		color.Red("List of commands : ")
		color.Red("Test a GoS method (func) : m <method name> args...(Can use golang statements as well)")
		color.Red("Test a template : t <template name> <json_of_interface(optional)>")
		color.Red("Test a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
		color.Red("Test a func in current main package : f <func name> args...(Use golang statements as well)")

		color.Green("Help needed with Event keys.")

		return
	}

	if len(cmd_set) > 2 {
		args_c = strings.Join(cmd_set[2:], ",")
	} else {
		args_c = ``
	}

	if cmd_set[0] == "m" {
		//[2:]

		templat := `package main

			import "testing"

			func Testnet_` + cmd_set[1] + `(t *testing.T){
				usr := net_` + cmd_set[1] + ` (` + args_c + `)
				if net_` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "f" {
		//[2:]

		templat := `package main

			import "testing"

			func Test` + cmd_set[1] + `(t *testing.T){
				usr := ` + cmd_set[1] + ` (` + args_c + `)
				if ` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "t" {

		if args_c != "" {
			args_c = `"` + args_c + `"`
		}

		templat := `package main

			import "testing"

			func Testnet_` + cmd_set[1] + `(t *testing.T){
				usr := net_` + cmd_set[1] + ` (` + args_c + `)
				if net_` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "p" {

		var method = "GET"
		var path = cmd_set[1]
		var params = "nil"

		if len(cmd_set) > 3 {
			method = cmd_set[3]
		}

		templat := `
			package main

			import (
			    "net/http"
			    "net/http/httptest"
			    "testing"`

		if len(cmd_set) > 2 {
			templat += `"bytes"`
			params = ``
		}

		templat += `
			)

		

			func Test(t *testing.T) {
			    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			    // pass 'nil' as the third parameter.
			    req, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        t.Fatal(err)
			    }

			      reqtwo, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        t.Fatal(err)
			    }

			    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			    rr := httptest.NewRecorder()
			    handle := http.HandlerFunc(makeHandler(handler))

			    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
			    // directly and pass in our Request and ResponseRecorder.
			 
			   		rrtwo := httptest.NewRecorder()
			   		handle.ServeHTTP(rrtwo, reqtwo) 
			   		expected := rrtwo.Body.String()
			   	

			   	handle.ServeHTTP(rr, req)

			    // Check the status code is what we expect.
			    if status := rr.Code; status != http.StatusOK {
			        t.Errorf("handler returned wrong status code: got %v want %v",
			            status, http.StatusOK)
			    }

			    // Check the response body is what we expect.
			 
			  
			    if rr.Body.String() != expected {
			        t.Errorf("handler returned unexpected body: got %v want %v",
			            rr.Body.String(), expected)
			    }
			}`

		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed! " + err.Error())
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	}

}

func Vmd() {

	log.Println(">>")
	var args_c string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	cmd := scanner.Text()
	cmd_set := strings.Split(cmd, " ")

	if len(cmd_set) < 2 {

		color.Red("List of commands : ")
		color.Red("Test a GoS method (func) : m <method name> args...(Can use golang statements as well)")
		color.Red("Test a template : t <template name> <json_of_interface(optional)>")
		color.Red("Test a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
		color.Red("Test a func in current main package : f <func name> args...(Use golang statements as well)")

		color.Green("Help needed with Event keys.")
		Vmd()
		return
	}

	if len(cmd_set) > 2 {
		args_c = strings.Join(cmd_set[2:], ",")
	} else {
		args_c = ``
	}

	if cmd_set[0] == "m" {
		//[2:]

		templat := `package main

			import "testing"

			func Testnet_` + cmd_set[1] + `(t *testing.T){
				usr := net_` + cmd_set[1] + ` (` + args_c + `)
				if net_` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "f" {
		//[2:]

		templat := `package main

			import "testing"

			func Test` + cmd_set[1] + `(t *testing.T){
				usr := ` + cmd_set[1] + ` (` + args_c + `)
				if ` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "t" {

		if args_c != "" {
			args_c = `"` + args_c + `"`
		}

		templat := `package main

			import "testing"

			func Testnet_` + cmd_set[1] + `(t *testing.T){
				usr := net_` + cmd_set[1] + ` (` + args_c + `)
				if net_` + cmd_set[1] + `(` + args_c + `) != usr {
					t.Error("...")
				}
			}`
		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	} else if cmd_set[0] == "p" {

		var method = "GET"
		var path = cmd_set[1]
		var params = "nil"

		if len(cmd_set) > 3 {
			method = cmd_set[3]
		}

		templat := `
			package main

			import (
			    "net/http"
			    "net/http/httptest"
			    "testing"`

		if len(cmd_set) > 2 {
			templat += `"bytes"`
			params = `bytes.NewReader( []byte("` + cmd_set[2] + `") )`
		}

		templat += `
			)

		

			func Test(t *testing.T) {
			    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			    // pass 'nil' as the third parameter.
			    req, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        t.Fatal(err)
			    }

			      reqtwo, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        t.Fatal(err)
			    }

			    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			    rr := httptest.NewRecorder()
			    handle := http.HandlerFunc(makeHandler(handler))

			    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
			    // directly and pass in our Request and ResponseRecorder.
			 
			   		rrtwo := httptest.NewRecorder()
			   		handle.ServeHTTP(rrtwo, reqtwo) 
			   		expected := rrtwo.Body.String()
			   	

			   	handle.ServeHTTP(rr, req)

			    // Check the status code is what we expect.
			    if status := rr.Code; status != http.StatusOK {
			        t.Errorf("handler returned wrong status code: got %v want %v",
			            status, http.StatusOK)
			    }

			    // Check the response body is what we expect.
			 
			  
			    if rr.Body.String() != expected {
			        t.Errorf("handler returned unexpected body: got %v want %v",
			            rr.Body.String(), expected)
			    }
			}`

		ioutil.WriteFile("test_internal_test.go", []byte(templat), 0777)
		color.Magenta("Running test...")
		log_build, err := core.RunCmdSmartZ("go test")
		if err != nil {
			color.Red("Test failed! " + err.Error())
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_internal_test.go")

	}

	Vmd()
}

func VmT() {

	log.Println("Trace >>")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	cmd := scanner.Text()
	cmd_set := strings.Split(`p `+cmd, " ")

	if len(cmd_set) < 2 {

		color.Red("List of commands : ")
		color.Red("Trace a server path (API | page) : </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")

		color.Green("Help needed with Event keys.")
		VmT()
		return
	}

	if cmd_set[0] == "p" {

		var method = "GET"
		var path = cmd_set[1]
		var params = "nil"

		if len(cmd_set) > 3 {
			method = cmd_set[3]
		}

		templat := `
			package main

			import (
			    "net/http"
			    "net/http/httptest"
			    "testing"
			    `

		if len(cmd_set) > 2 {
			templat += `"bytes"`
			params = `bytes.NewReader( []byte("` + cmd_set[2] + `") )`
		}

		templat += `
			)
			var result int

			func GWeb(b *testing.B){
				  req, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        b.Fatal(err)
			    }

			   

			    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			    rr := httptest.NewRecorder()
			    handle := http.HandlerFunc(makeHandler(handler))

			    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
			    // directly and pass in our Request and ResponseRecorder.
			 
			   	handle.ServeHTTP(rr, req)

			    // Check the status code is what we expect.
			    if status := rr.Code; status != http.StatusOK {
			        b.Errorf("handler returned wrong status code: got %v want %v",
			            status, http.StatusOK)
			    }

			    // Check the response body is what we expect.
			 
			}

			func BenchmarkGWeb(b *testing.B) {
			    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			    // pass 'nil' as the third parameter.

			  
				GWeb(b)
				
			}`

		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed! " + err.Error())
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")
		process := make(chan bool)
		go core.Exe_Stall("go tool trace __heap", process)
		color.Yellow("Hit enter to trace another path.")
		scanner.Scan() // use `for scanner.Scan()` to keep reading
		cmd = scanner.Text()
		process <- true

	}

	VmT()
}

func VmP() {

	log.Println("Bench >>")
	var args_c string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	cmd := scanner.Text()
	cmd_set := strings.Split(cmd, " ")

	if len(cmd_set) < 2 {

		color.Red("List of commands : ")
		color.Red("Benchmark a GoS method (func) : m <method name> args...(Can use golang statements as well)")
		color.Red("Benchmark a template : t <template name> <json_of_interface(optional)>")
		color.Red("Benchmark a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
		color.Red("Benchmark a func in current main package : f <func name> args...(Use golang statements as well)")

		color.Green("Help needed with Event keys.")

		return
	}

	if len(cmd_set) > 2 {
		args_c = strings.Join(cmd_set[2:], ",")
	} else {
		args_c = ``
	}

	if cmd_set[0] == "m" {
		//[2:]

		templat := `package main

import "testing"

var result int

func BenchmarkNet_` + cmd_set[1] + `(b *testing.B) {  
	var r int
	for n := 0; n < b.N; n++ {
		
		r = 0
		net_` + cmd_set[1] + `(` + args_c + `)
	}
	
	result = r
}

`

		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "f" {
		//[2:]

		templat := `package main

			import "testing"

			var result int

			func Benchmark` + cmd_set[1] + `(b *testing.B) {  
				var r int
				for n := 0; n < b.N; n++ {
					
					r = 0
					` + cmd_set[1] + `(` + args_c + `)
				}
				
				result = r
			}

			`
		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "t" {

		if args_c != "" {
			args_c = `"` + args_c + `"`
		}

		templat := `package main

			import "testing"

			var result int

			func BenchmarkNet_` + cmd_set[1] + `(b *testing.B) {  
				var r int
				for n := 0; n < b.N; n++ {
					
					r = 0
					net_` + cmd_set[1] + `(` + args_c + `)
				}
				
				result = r
			}

			`
		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "p" {

		var method = "GET"
		var path = cmd_set[1]
		var params = "nil"

		if len(cmd_set) > 3 {
			method = cmd_set[3]
		}

		templat := `
			package main

			import (
			    "net/http"
			    "net/http/httptest"
			    "testing"`

		if len(cmd_set) > 2 {
			templat += `"bytes"`
			params = `bytes.NewReader( []byte("` + cmd_set[2] + `") )`
		}

		templat += `
			)
			var result int

			func GWeb(b *testing.B){
				  req, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        b.Fatal(err)
			    }

			   

			    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			    rr := httptest.NewRecorder()
			    handle := http.HandlerFunc(makeHandler(handler))

			    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
			    // directly and pass in our Request and ResponseRecorder.
			 
			   	handle.ServeHTTP(rr, req)

			    // Check the status code is what we expect.
			    if status := rr.Code; status != http.StatusOK {
			        b.Errorf("handler returned wrong status code: got %v want %v",
			            status, http.StatusOK)
			    }

			    // Check the response body is what we expect.
			 
			}

			func BenchmarkGWeb(b *testing.B) {
			    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			    // pass 'nil' as the third parameter.

			    var r int
				
					for n := 0; n < b.N; n++ {
					r = 0
			  
					GWeb(b)
			   
				}
				result = r
			}`

		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed! " + err.Error())
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	}

	VmP()
}

func VmPOne() {

	log.Println("Bench >>")
	var args_c string

	cmd_set := os.Args[2:]

	if len(cmd_set) < 2 {

		color.Red("List of commands : ")
		color.Red("Benchmark a GoS method (func) : m <method name> args...(Can use golang statements as well)")
		color.Red("Benchmark a template : t <template name> <json_of_interface(optional)>")
		color.Red("Benchmark a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
		color.Red("Benchmark a func in current main package : f <func name> args...(Use golang statements as well)")

		color.Green("Help needed with Event keys.")

		return
	}

	if len(cmd_set) > 2 {
		args_c = strings.Join(cmd_set[2:], ",")
	} else {
		args_c = ``
	}

	if cmd_set[0] == "m" {
		//[2:]

		templat := `package main

import "testing"

var result int

func BenchmarkNet_` + cmd_set[1] + `(b *testing.B) {  
	var r int
	for n := 0; n < b.N; n++ {
		
		r = 0
		net_` + cmd_set[1] + `(` + args_c + `)
	}
	
	result = r
}

`

		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "f" {
		//[2:]

		templat := `package main

			import "testing"

			var result int

			func Benchmark` + cmd_set[1] + `(b *testing.B) {  
				var r int
				for n := 0; n < b.N; n++ {
					
					r = 0
					` + cmd_set[1] + `(` + args_c + `)
				}
				
				result = r
			}

			`
		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "t" {

		if args_c != "" {
			args_c = `"` + args_c + `"`
		}

		templat := `package main

			import "testing"

			var result int

			func BenchmarkNet_` + cmd_set[1] + `(b *testing.B) {  
				var r int
				for n := 0; n < b.N; n++ {
					
					r = 0
					net_` + cmd_set[1] + `(` + args_c + `)
				}
				
				result = r
			}

			`
		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed!")
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	} else if cmd_set[0] == "p" {

		var method = "GET"
		var path = cmd_set[1]
		var params = "nil"

		if len(cmd_set) > 3 {
			method = cmd_set[3]
		}

		templat := `
			package main

			import (
			    "net/http"
			    "net/http/httptest"
			    "testing"`

		if len(cmd_set) > 2 {
			templat += `"bytes"`
			params = `bytes.NewReader( []byte("` + cmd_set[2] + `") )`
		}

		templat += `
			)
			var result int

			func GWeb(b *testing.B){
				  req, err := http.NewRequest("` + method + `", "` + path + `", ` + params + `)
			    if err != nil {
			        b.Fatal(err)
			    }

			   

			    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			    rr := httptest.NewRecorder()
			    handle := http.HandlerFunc(makeHandler(handler))

			    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
			    // directly and pass in our Request and ResponseRecorder.
			 
			   	handle.ServeHTTP(rr, req)

			    // Check the status code is what we expect.
			    if status := rr.Code; status != http.StatusOK {
			        b.Errorf("handler returned wrong status code: got %v want %v",
			            status, http.StatusOK)
			    }

			    // Check the response body is what we expect.
			 
			}

			func BenchmarkGWeb(b *testing.B) {
			    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			    // pass 'nil' as the third parameter.

			    var r int
				
					for n := 0; n < b.N; n++ {
					r = 0
			  
					GWeb(b)
			   
				}
				result = r
			}`

		ioutil.WriteFile("test_test.go", []byte(templat), 0777)
		color.Magenta("Running benchmark...")
		log_build, err := core.RunCmdSmartP("go test -bench=.")
		if err != nil {
			color.Red("Test failed! " + err.Error())
		} else {
			color.Green("Success")

		}
		log.Println(log_build)
		os.Remove("test_test.go")

	}

}

func WatchForUpdate(path string) {

	done := make(chan bool)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	reloading := false

	go func() {
		for {
			select {
			case ev := <-watcher.Event:

				if !reloading {
					fname := fmt.Sprintf("%v", ev)
					if TriggerType(fname) && !strings.Contains(fname, "bindata.go") && !strings.Contains(fname, appout) {
						//Build( GOHOME + "/" + serverconfig )
						reloading = true
						done <- true
						break
					}
				}

			}
		}
	}()
	err = watcher.Watch(fmt.Sprintf("%s/", GOHOME))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting for new updates to source.")
	<-done
	defer close(done)
	watcher.RemoveWatch(path)
	watcher.Close()
	log.Println("🤔 Refreshing")
	core.RunCmd("gos --t")
	JBuild(path, appout)
}

func TriggerType(typ string) (is bool) {
	is = strings.Contains(typ, ".go") || strings.Contains(typ, ".css") || strings.Contains(typ, ".tmpl") || strings.Contains(typ, ".gxml") || strings.Contains(typ, ".html") || strings.Contains(typ, ".js")
	return
}

func JBuild(path string, out string) {
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
	//time.Sleep(time.Second*100 )
	//core.RunFile(GOHOME, coreTemplate.Output)

	if Type == "webapp" {
		chh := make(chan int)
		go core.DoSpin(chh)
		log_build, err := core.RunCmdSmart("go build")
		chh <- 1
		close(chh)
		if err != nil {
			//log.Println(err.Error())
			color.Red("Your build failed, Here is why :>")
			lines := strings.Split(log_build, "\n")
			for i, line := range lines {
				if i > 0 {
					if strings.Contains(line, "imported and") {
						line_part := strings.Split(line, ":")
						color.Red(strings.Join(line_part[2:], " - "))
					} else {
						if line != "" {
							line_part := strings.Split(line, ":")
							if len(line_part) > 3 && strings.Contains(line_part[0], out) {
								lnumber, _ := strconv.Atoi(line_part[1])

								file, err := os.Open(out)
								if err != nil {
									color.Red("Could not find a source file")
									return
								}

								//log.Println(line_part[len(line_part) - 1])
								scanner := bufio.NewScanner(file)
								inm := 0
								for scanner.Scan() {
									inm++
									//log.Println("%+V", inm)
									lin := scanner.Text()
									if inm == lnumber {
										acT_line := GetLine(serverconfig, lin)
										if acT_line > -1 {
											color.Magenta(fmt.Sprintf("Verify your file %s on line : %v | %s", serverconfig, acT_line, strings.Join(line_part[2:], " - ")))

										} else {
											color.Blue(fmt.Sprintf("Verify your golang WebApp libraries (linked libraries) file : %s line : %s", line_part[0], line_part[1]))

										}
									}
									// log.Println("data : " + scanner.Text())
								}

								if err := scanner.Err(); err != nil {
									color.Red("Could not find a source file")
									return
								}

								file.Close()
							} else {
								color.Blue(fmt.Sprintf("Verify your golang WebApp libraries (linked libraries) file : %s line : %s", line_part[0], line_part[1]))
							}
						}
					}
				}
			}
			notify.Push("Build failed!", fmt.Sprintf("Your project %s failed to build!", pk[len(pk)-1]), "", notificator.UR_CRITICAL)
			color.Red("Full compiler build log : ")
			log.Println(log_build)
			WatchForUpdate(path)
			return
		}

		log.Println("Use Ctrl + C to quit")

	}

	process := make(chan bool)
	done := make(chan bool)
	//	log_console := make(chan string)

	//always delete add on folders prior
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	reloading := false

	go func() {
		for {
			select {
			case ev := <-watcher.Event:

				if !reloading {
					fname := fmt.Sprintf("%v", ev)

					if TriggerType(fname) && !strings.Contains(fname, "bindata.go") && !strings.Contains(fname, appout) {
						//Build( GOHOME + "/" + serverconfig )
						reloading = true
						if !strings.Contains(Type, "faas") {
							process <- true
						}
						done <- true
						break
					}
				}

			}
		}
	}()

	err = watcher.Watch(fmt.Sprintf("%s/", GOHOME))
	if err != nil {
		log.Fatal(err)
	}
	notify.Push("Build Passed!", fmt.Sprintf("Your project %s is running!", pk[len(pk)-1]), fmt.Sprintf("%s/src/github.com/cheikhshift/gos/icon.png", os.ExpandEnv("$GOPATH")), notificator.UR_NORMAL)
	log.Println("Ready!")
	if !strings.Contains(Type, "faas") {
		go core.Exe_Stall(fmt.Sprintf("./%s", pk[len(pk)-1]), process)
	}
	<-done
	defer close(process)
	defer close(done)
	watcher.RemoveWatch(path)
	watcher.Close()
	log.Println("🤔 Refreshing")
	log.Println("📦 Invoking go-bindata")

	if strings.Contains(Type, "faas") {

		core.RunCmd(fmt.Sprintf("go-bindata -pkg=%s %s/... %s/...", pk[len(pk)-1], webroot, templateroot))

	} else {
		core.RunCmd(fmt.Sprintf("go-bindata -debug %s/... %s/...", webroot, templateroot))
	}

	chn := make(chan int)
	go core.DoSpin(chn)
	core.RunCmd("gos --t")
	chn <- 1
	close(chn)
	JBuild(path, out)
}

func Build(path string) {

	color.Green("😎 Loading project!")
	coreTemplate, err := core.LoadGos(path)
	if err != nil {
		log.Fatal(err)
	}

	if coreTemplate == nil {
		return
	}

	Type = coreTemplate.Type
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
	if !strings.Contains(os.Args[1], "export") {
		log.Println("📦 Invoking go-bindata ")

		if strings.Contains(Type, "faas") {
			core.RunCmd(fmt.Sprintf("go-bindata -pkg=%s %s/... %s/...", pk[len(pk)-1], webroot, templateroot))

		} else {
			core.RunCmd(fmt.Sprintf("go-bindata -debug %s/... %s/...", webroot, templateroot))
		}
	} else {

		if strings.Contains(Type, "faas") {

			core.RunCmd(fmt.Sprintf("go-bindata -pkg=%s %s/... %s/...", pk[len(pk)-1], webroot, templateroot))

		} else {
			core.RunCmd(fmt.Sprintf("go-bindata %s/... %s/...", webroot, templateroot))
		}
	}
	appout = coreTemplate.Output
	//log.Println(coreTemplate.Methods.Methods)
	coreTemplate.WriteOut = false
	coreTemplate.Name = path
	//log.Println(coreTemplate)
	if os.Args[1] == "export" || os.Args[1] == "export-sub" || os.Args[1] == "--export" {
		coreTemplate.Prod = true
	}

	if os.Args[1] == "--trace" {
		coreTemplate.Debug = "on"
	}

	core.Process(coreTemplate, GOHOME, webroot, templateroot)

	cwd, er := os.Getwd()
	if er != nil {
		log.Println(er)
		os.Exit(1)
	}
	pkgpath := strings.Split(strings.Trim(cwd, "/"), "/")

	if isWin := strings.Contains(runtime.GOOS, "indows"); isWin {
		core.RunCmd(fmt.Sprintf("gofmt -w %s", pkgpath[len(pkgpath)-1]))
	} else {
		core.RunCmd(fmt.Sprintf("gofmt -w ../%s", pkgpath[len(pkgpath)-1]))
	}

	if os.Args[1] == "--t" {
		return
	}

	if coreTemplate.Type == "webapp" || coreTemplate.Type == "faas" {

		if os.Args[1] == "run" || os.Args[1] == "run-sub" || os.Args[1] == "--run" || os.Args[1] == "--serv" {
			//
			if !strings.Contains(os.Args[1], "run-") && !strings.Contains(os.Args[1], "--") {
				os.Chdir(GOHOME)
			}

			//time.Sleep(time.Second*100 )
			//core.RunFile(GOHOME, coreTemplate.Output)

			if coreTemplate.Type == "webapp" {
				//fmt.Printf("\r  \033[36mBuilding\033[m %s ", s.Next())
				log_build, err := core.RunCmdSmart("go build")
				if err != nil {
					//log.Println(err.Error())
					color.Red("😎 Your build failed, Here is why :>")
					lines := strings.Split(log_build, "\n")
					for i, line := range lines {
						if i > 0 {
							if strings.Contains(line, "imported and") {
								line_part := strings.Split(line, ":")
								color.Red(strings.Join(line_part[2:], " - "))
							} else {
								if line != "" {
									line_part := strings.Split(line, ":")

									if len(line_part) > 3 && strings.Contains(line_part[0], appout) {
										lnumber, _ := strconv.Atoi(line_part[1])

										file, err := os.Open(appout)
										if err != nil {
											color.Red("Could not find a source file")
											return
										}

										//log.Println(line_part[len(line_part) - 1])
										scanner := bufio.NewScanner(file)
										inm := 0
										for scanner.Scan() {
											inm++
											//log.Println("%+V", inm)
											lin := scanner.Text()
											if inm == lnumber {
												acT_line := GetLine(serverconfig, lin)
												if acT_line > -1 {
													color.Magenta(fmt.Sprintf("Verify your file %s on line : %v | %s", serverconfig, acT_line, strings.Join(line_part[2:], " - ")))

												}
											}
											// log.Println("data : " + scanner.Text())
										}

										if err := scanner.Err(); err != nil {
											color.Red("Could not find a source file")
											return
										}

										file.Close()
									} else if len(line_part) > 3 {
										color.Blue(fmt.Sprintf("Verify your golang WebApp libraries (linked libraries) file : %s line : %s reason : %s", line_part[0], line_part[1], line_part[2:]))
									}
								}
							}
						}
					}
					notify.Push("Build failed!", fmt.Sprintf("Your project %s failed to build!", pk[len(pk)-1]), "", notificator.UR_CRITICAL)
					color.Red("Full compiler build log : ")
					log.Println(log_build)
					WatchForUpdate(path)
					return
				}

				if len(os.Args) > 2 {
					if os.Args[2] == "--buildcheck" {
						return
					}
				}
				log.Println("🤓 Use Ctrl + C to quit")
			}
			process := make(chan bool)
			done := make(chan bool)
			//log_console := make(chan string)

			//always delete add on folders prior
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}

			//brokep := false
			reloading := false
			// Process events
			go func() {
				for {
					select {

					case ev := <-watcher.Event:
						if !reloading {
							fname := fmt.Sprintf("%v", ev)

							if TriggerType(fname) && !strings.Contains(fname, "bindata.go") && !strings.Contains(fname, appout) {
								//Build( GOHOME + "/" + serverconfig )
								reloading = true
								if !strings.Contains(Type, "faas") {
									process <- true
								}
								//	done <- true
								done <- true
								break
							}
						} /* else if !brokep {

						}*/
					}
				}
			}()

			err = watcher.Watch(fmt.Sprintf("%s/", GOHOME))
			if err != nil {
				log.Fatal(err)
			}
			notify.Push("Build Passed!", fmt.Sprintf("Your project %s is running!", pk[len(pk)-1]), strings.Replace(fmt.Sprintf("%s/src/github.com/cheikhshift/gos/icon.png", os.ExpandEnv("$GOPATH")), "//", "/", -1), notificator.UR_NORMAL)
			log.Println("Ready!")
			if !strings.Contains(Type, "faas") {
				go core.Exe_Stall(fmt.Sprintf("./%s", pk[len(pk)-1]), process)
			}
			//process <- false
			<-done

			defer close(process)
			defer close(done)
			watcher.RemoveWatch(path)
			watcher.Close()
			log.Println("🤔 Refreshing")
			log.Println("📦 Invoking go-bindata")

			if strings.Contains(Type, "faas") {

				core.RunCmd(fmt.Sprintf("go-bindata  -pkg=%s %s/... %s/...", pk[len(pk)-1], webroot, templateroot))

			} else {
				core.RunCmd(fmt.Sprintf("go-bindata -debug %s/... %s/...", webroot, templateroot))
			}
			chn := make(chan int)
			go core.DoSpin(chn)
			core.RunCmd("gos --t")
			chn <- 1
			close(chn)
			JBuild(path, coreTemplate.Output)
			return

		}

		if os.Args[1] == "--trace" {
			color.Red("List of commands : ")
			color.Red("Trace a server path (API | page) : </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")

			color.Green("Help needed with Event keys.")
			VmT()

		}

		if os.Args[1] == "--test" {
			//test console
			log.Println("📦 Invoking go-bindata")
			core.RunCmd("go-bindata -debug " + webroot + "/... " + templateroot + "/...")
			color.Magenta("Welcome to the Gopher Sauce test console.")
			color.Red("List of commands : ")
			color.Red("Test a GoS method (func) : m <method name> args...(Can use golang statements as well)")
			color.Red("Test a template : t <template name> <json_of_interface(optional)>")
			color.Red("Test a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
			color.Red("Test a func in current main package : f <func name> args...(Use golang statements as well)")

			color.Green("Help needed with Event keys.")
			Vmd()
		}

		if os.Args[1] == "--test-f" {
			//test console
			log.Println("📦 Invoking go-bindata")
			core.RunCmd("go-bindata -debug " + webroot + "/... " + templateroot + "/...")
			VmdOne()
		}

		if os.Args[1] == "--bench" {
			//test console
			log.Println("📦 Invoking go-bindata")
			core.RunCmd("go-bindata -debug " + webroot + "/... " + templateroot + "/...")
			color.Magenta("Welcome to the Gopher Sauce benchmark console.")
			color.Red("List of commands : ")
			color.Red("Benchmark a GoS method (func) : m <method name> args...(Can use golang statements as well)")
			color.Red("Benchmark a template : t <template name> <json_of_interface(optional)>")
			color.Red("Benchmark a server path (API | page) : p </path/to/resource/without/hostname/> <json_of_request(optional)> <method_of_request(optional)> ")
			color.Red("Benchmark a func in current main package : f <func name> args...(Use golang statements as well)")

			color.Green("Help needed with Event keys.")
			VmP()
		}

		if os.Args[1] == "--bench-f" {
			//test console
			log.Println("📦 Invoking go-bindata")
			core.RunCmd("go-bindata -debug " + webroot + "/... " + templateroot + "/...")

			VmPOne()
		}

		if os.Args[1] == "export" || os.Args[1] == "export-sub" || os.Args[1] == "--export" {
			log.Println("Generating Export Program")
			if !strings.Contains(os.Args[1], "export-") && !strings.Contains(os.Args[1], "--export") {
				os.Chdir(GOHOME)
			}
			//create both zips
			core.RunCmd("go build")
		}
	} else if coreTemplate.Type == "bind" {

		//check for directory gomobile
		if os.Args[1] == "export" {
			log.Println("Generating Export Program")
			os.Chdir(GOHOME)
			//create both zips
			log.Println("📦 Invoking go-bindata")
			core.RunCmd(os.ExpandEnv("$GOPATH") + `/bin/go-bindata ` + webroot + "/... " + templateroot + "/...")
			body, er := ioutil.ReadFile(GOHOME + "/bindata.go")
			if er != nil {
				log.Println(er)
				return
			}
			writeLocalProtocol(coreTemplate.Package)
			log.Println("Preparing Bindata for framework conversion...")
			ioutil.WriteFile("bindata.go", prepBindForMobile(body, coreTemplate.Package), 0644)
			core.RunCmd(os.ExpandEnv("$GOPATH") + "/bin/gomobile bind -target=ios")
			//edit plist file
			subp := "/sub.check"
			_, error := ioutil.ReadFile(subp)
			if error != nil {
				ioutil.WriteFile(subp, []byte("StubCompletion"), 0644)
				pathSp := strings.Split(os.Args[6], "/")
				finalSub := ""
				if len(pathSp) > 1 {
					finalSub = pathSp[len(pathSp)-1]
				} else {
					finalSub = os.Args[6]
				}
				plistPath := os.Args[6] + "/" + finalSub + "/Info.plist"
				plist, erro := ioutil.ReadFile(plistPath)
				if erro != nil {
					log.Println("Please check your project's folder for the Info.plist (Info.plisp chuckles...) file")
					return
				}

				ioutil.WriteFile(plistPath, []byte(strings.Replace(string(plist), `<key>UIMainStoryboardFile</key>
							<string>Main</string>`, `<key>UIBackgroundModes</key>
						<array>
						    <string>fetch</string>
						</array>`, -1)), 0644)

				core.RunCmd("python " + os.ExpandEnv("$GOPATH") + "/src/github.com/cheikhshift/gos/core/addFlow.py " + strings.Trim(os.Args[2], "/") + " " + os.Args[6] + " " + UpperInitial(coreTemplate.Package))
				//if project does not exist create it and link this framework

			} else {
				log.Println("Subexists no need for Linkage :o")
			}
		}

	}
}

func main() {
	if tpath := os.ExpandEnv("$USERPROFILE"); tpath != "" && os.ExpandEnv("$GOPATH") == "" {
		os.Setenv("GOPATH", tpath+"/go")
	}
	GOHOME = fmt.Sprintf("%s/src/", os.ExpandEnv("$GOPATH"))

	if GOHOME == "/src/" {
		log.Println("Please set a GOPATH prior to using this program.")
		return
	}
	notify = notificator.New(notificator.Options{
		DefaultIcon: strings.Replace(fmt.Sprintf("%s/src/github.com/cheikhshift/gos/icon.png", os.ExpandEnv("$GOPATH")), "//", "/", -1),
		AppName:     "Gopher Sauce",
	})
	//log.Println( os.Args)
	if len(os.Args) > 1 {
		//args := os.Args[1:]
		if os.Args[1] == "dependencies" || os.Args[1] == "deps" {
			log.Println("∑ Getting GoS dependencies")
			core.RunCmd("go get -u github.com/jteeuwen/go-bindata/...")
			core.RunCmd("go get github.com/gorilla/sessions")
			core.RunCmd("go get github.com/elazarl/go-bindata-assetfs")
			core.RunCmd("go get -u github.com/golang/dep/cmd/dep")
			core.RunCmd("go get github.com/gorilla/context")
			core.RunCmd("go get gopkg.in/mgo.v2")
			core.RunCmd("go get github.com/asaskevich/govalidator")
			core.RunCmd("go get sourcegraph.com/sourcegraph/appdash")
			core.RunCmd("go get sourcegraph.com/sourcegraph/appdash-data")
			core.RunCmd("go get github.com/fatih/color")
			core.RunCmd("go get github.com/opentracing/opentracing-go")
			core.RunCmd("go get github.com/gorilla/mux")
			core.RunCmd("go get github.com/opentracing/basictracer-go")
			core.RunCmd("go get github.com/cheikhshift/db")
			core.RunCmd("go get gopkg.in/ldap.v2")

			//core.RunCmd("")
			//log.Println("ChDir " + os.ExpandEnv("$GOPATH") + "/src/github.com/kronenthaler/mod-pbxproj")
			//os.Chdir(os.ExpandEnv("$GOPATH") + "/src/github.com/kronenthaler/mod-pbxproj")
			//core.RunCmd("python setup.py install" )
			//time.Sleep(time.Second *120)
			log.Println("Done")
			return
		}

		if os.Args[1] == "--kill" {

			if strings.Contains(runtime.GOOS, "indows") {
				log.Println("This flag is not support on windows yet...")
				return
			}
			wdir, _ := os.Getwd()
			processname := strings.Split(wdir, "/")

			core.RunCmdSmart(fmt.Sprintf("killall -3 %s", processname[len(processname)-1]))
			core.RunCmdSmart("killall -3 gos")
			return
		}

		if os.Args[1] == "make" {
			//2 is project folder

			os.MkdirAll(os.ExpandEnv("$GOPATH")+"/src/"+strings.Trim(os.Args[2], "/")+"/web", 0777)
			os.MkdirAll(os.ExpandEnv("$GOPATH")+"/src/"+strings.Trim(os.Args[2], "/")+"/tmpl", 0777)
			ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/"+strings.Trim(os.Args[2], "/")+"/gos.gxml", []byte(gosTemplate), 0777)
			ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/"+strings.Trim(os.Args[2], "/")+"/web/your-404-page.tmpl", []byte(htmlTemplate), 0777)
			ioutil.WriteFile(os.ExpandEnv("$GOPATH")+"/src/"+strings.Trim(os.Args[2], "/")+"/web/your-500-page.tmpl", []byte(htmlTemplate), 0777)

			return
		}

		if os.Args[1] == "makesublime" || os.Args[1] == "--make" {
			//2 is project folder

			os.MkdirAll("web", 0777)
			os.MkdirAll("tmpl", 0777)
			ioutil.WriteFile("gos.gxml", []byte(gosTemplate), 0777)
			ioutil.WriteFile("web/your-404-page.tmpl", []byte(htmlTemplate), 0777)
			ioutil.WriteFile("web/your-500-page.tmpl", []byte(htmlTemplate), 0777)
			return
		}

		if strings.Contains(os.Args[1], "sub") {
			GOHOME = "./"

		}

		if strings.Contains(os.Args[1], "--") {
			GOHOME = "./"

		} else {
			GOHOME = GOHOME + strings.Trim(os.Args[2], "/")
		}

		if strings.Contains(os.Args[1], "--") {
			webroot = "web"
			templateroot = "tmpl"
			serverconfig = "gos.gxml"
		} else {
			webroot = os.Args[4]
			templateroot = os.Args[5]
			serverconfig = os.Args[3]
		}

		Build(fmt.Sprintf("%s/%s", GOHOME, serverconfig))

	} /* else {

		log.Println("To begin please tell us a bit about the gos project you wish to compile");
		log.Printf("We need the GoS package folder relative to your $GOPATH/src (%v)\n", GOHOME)
	   	gosProject := ""
	   	serverconfig := ""

	   	fmt.Scanln(&gosProject)
	   	GOHOME = GOHOME  + strings.Trim(gosProject,"/")
	   	log.Printf("We need your Gos Project config source (%v)\n", GOHOME)
	   	fmt.Scanln(&serverconfig)
	    //log.Println(GOHOME)
		webroot,templateroot = core.DoubleInput("What is the name of your webroot's folder ?", "What is the name of your template folder? ")
			log.Println("Are you ready to begin? ");
			if core.AskForConfirmation() {
				log.Println("ΩΩ Operation Started!!");
				coreTemplate,err := core.LoadGos( GOHOME + "/" + serverconfig );
				if err != nil {
					log.Println(err)
					return
				}

				coreTemplate.WriteOut = false
				core.Process(coreTemplate,GOHOME, webroot,templateroot);
				log.Println("One moment...")
				core.RunCmd("go get -u github.com/jteeuwen/go-bindata/...")
	    	    core.RunCmd("go get github.com/gorilla/sessions")
	    		core.RunCmd("go get github.com/elazarl/go-bindata-assetfs")
				log.Println("Would you like to just run this application [y,n]")

				if core.AskForConfirmation() {
					os.Chdir(GOHOME)
					log.Println("Invoking go-bindata");
					core.RunCmd("go-bindata -debug " + webroot +"/... " + templateroot + "/...")
					//time.Sleep(time.Second*100 )
					//core.RunFile(GOHOME, coreTemplate.Output)
					core.RunCmd("go build")
					pk := strings.Split(strings.Trim(gosProject,"/"), "/")
					log.Println("Use Ctrl + C to quit")
					core.Exe_Stall("./" + pk[len(pk) - 1] )

				} else {
					log.Println("Is this a Mobile application [y,n]")

					if !core.AskForConfirmation() {
					log.Println("Would you like to create an export release [y,n]")

					if core.AskForConfirmation() {
						log.Println("Generating Export Program")
						os.Chdir(GOHOME)
						//create both zips
						log.Println("Invoking go-bindata");
						core.RunCmd(  "go-bindata  " + webroot +"/... " + templateroot + "/...")
						core.RunCmd("go build")

					}
					} else {
						//create mobile export here
						log.Println("Generating Export Program")
							os.Chdir(GOHOME)
							//create both zips
							 log.Println("Invoking go-bindata");
							 core.RunCmd( os.ExpandEnv("$GOPATH") + `/bin/go-bindata `  + webroot +"/... " + templateroot + "/...")
							 body,er := ioutil.ReadFile(GOHOME + "/bindata.go")
							 if er != nil {
							 	log.Println(er)
							 	return
							 }
							 log.Println("Preparing Bindata for framework conversion...")
							 ioutil.WriteFile("bindata.go", prepBindForMobile(body, coreTemplate.Package)  ,0644)
							 core.RunCmd( os.ExpandEnv("$GOPATH")  + "/bin/gomobile bind -target=ios")
							 //edit plist file
							 subp := "sub.check"

							 log.Println("What is the folder name of your IOS application?")
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
							 	log.Println("Please check your project's folder for the Info.plit file")
							 	return
							 }
							 writeLocalProtocol(coreTemplate.Package)

							 ioutil.WriteFile(plistPath, []byte(strings.Replace(string(plist), `<key>UIMainStoryboardFile</key>
		<string>Main</string>`,``,-1)),0644 )

							 core.RunCmd("python " + os.ExpandEnv("$GOPATH") + "/src/github.com/cheikhshift/gos/core/addFlow.py " + strings.Trim(gosProject,"/") +" " + folderX + " " + UpperInitial(coreTemplate.Package))
							 //if project does not exist create it and link this framework

							} else {
								log.Println("Subexists no need for Linkage :o")
							}
							log.Println("Your file is ready, go on your default IDE and run your application :)")

					}
				}


			} else {
				log.Println("Operation Cancelled!!")
			}
		} */

}
