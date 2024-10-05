#import <Foundation/Foundation.h>
#include "ApplicationSupportDirectory.h"

const char* GetApplicationSupportDirectory() {
    NSArray* paths = NSSearchPathForDirectoriesInDomains(NSApplicationSupportDirectory, NSUserDomainMask, YES);
    NSString* appSupportDir = [paths firstObject];

    // Ensure the directory exists
    NSError *error = nil;
    if (![[NSFileManager defaultManager] createDirectoryAtPath:appSupportDir withIntermediateDirectories:YES attributes:nil error:&error]) {
        NSLog(@"Error creating Application Support directory: %@", error);
    }

    return [appSupportDir UTF8String];
}
