//
//  FlowTissue.m
//  GoTetst2
//
//  Created by OrendaCapital on 12/29/15.
//  Copyright © 2015 Cheikh Seck LLC. All rights reserved.
//

#import "FlowTissue.h"



@implementation FlowTissue


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
