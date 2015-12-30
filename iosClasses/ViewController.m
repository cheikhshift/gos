//
//  ViewController.m
//  FlowCode
//
//  Created by Cheikh Seck on 4/1/15.
//  Copyright (c) 2015 Orkiv LLC. All rights reserved.
//

#import "ViewController.h"

@interface ViewController()

@end


@implementation ViewController


- (void) viewDidAppear:(BOOL)animated {
    // NSLog(@"appe");
    [FlowThreadManager runJS:@"viewdidappear()"];
    // [[NSUserDefaults standardUserDefaults]
}
- (void)viewDidLoad {
    [super viewDidLoad];
    [NSURLProtocol registerClass:[FlowProtocol class]];
    // Do any additional setup after loading the view, typically from a nib.
    self.automaticallyAdjustsScrollViewInsets = NO;
    //   NSString *htmlFile = [[NSBundle mainBundle] pathForResource:@"root" ofType:@"html" inDirectory:@"SharedCode/Views/"];
    
    self.webView.delegate = [FlowThreadManager instance];
    self.webView.scrollView.bounces = NO;
    self.webView.mediaPlaybackRequiresUserAction = NO;
    [FlowThreadManager webviewCompletion:^(void){
        
    }];
    if(self.override){
        [self.webView loadRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:self.viewurl]] ];
    }
    else {
        [self.webView loadRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:@"http://localhost/index?test=foo"]] ];
    }
    
    [self setNeedsStatusBarAppearanceUpdate];
    
}

-(UIStatusBarStyle)preferredStatusBarStyle{
    return UIStatusBarStyleLightContent
    ;
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

@end
