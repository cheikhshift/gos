//
//  FlowThreadManager.h
//  FlowCode
//
//  Created by Cheikh Seck on 4/1/15.
//  Copyright (c) 2015 Gopher Sauce LLC. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <UIKit/UIKit.h>
#import "ViewController.h"

@interface FlowThreadManager : NSObject <UIWebViewDelegate>


typedef void(^Completion)(void);
typedef void(^CodeProcessCompletion)(NSString *result);

+ (void) runJS:(NSString *) function;
+ (FlowThreadManager *) instance;
+ (UIWebView *) currentFlow;
+ (NSString *) flowjs:(NSString *)function withData:(NSArray *) args;
+ (void) createFlowLayer;
+ (void) process: (NSString *) mvc completion:(CodeProcessCompletion) finished;
+ (id) getobject:(NSString *) name;
+ (BOOL) saveobject:(id)object withName:(NSString *) key;
+ (void) webviewCompletion:(Completion) finished;
+ (void) userDidCancelPayment;
+ (void) pinLogin;
@end
