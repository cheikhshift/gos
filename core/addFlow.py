from mod_pbxproj import XcodeProject

import sys
import os
import shutil

print 'Changing to ' + os.path.expandvars("$GOPATH") + "/src/" + sys.argv[1]
os.chdir(os.path.expandvars("$GOPATH") + "/src/" + sys.argv[1])
project = XcodeProject.Load( sys.argv[2] + "/" + sys.argv[2]  +'.xcodeproj/project.pbxproj')

iosclasses =  os.path.expandvars("$GOPATH") + "/src/github.com/cheikhshift/gos/iosClasses/"
group = project.get_or_create_group(sys.argv[2] )
frameworks = project.get_or_create_group("frameworks")
dest = sys.argv[2] + "/" + sys.argv[2] + "/Base.lproj/"

## By default delete all of the target names
project.remove_file_by_path('AppDelegate.m')
project.remove_file_by_path('ViewController.m')
project.remove_file_by_path('ViewController.h')
project.remove_file_by_path('FlowProtocol.m')
project.remove_file_by_path('LaunchImage.png')

project.add_file_if_doesnt_exist(iosclasses + "FlowProtocol.m",group)

shutil.copy2(iosclasses + "StoryBoards/Main.storyboard", dest + 'Main.storyboard')
shutil.copy2(iosclasses + "StoryBoards/LaunchScreen.storyboard", dest + 'LaunchScreen.storyboard')
shutil.copy2(iosclasses + "StoryBoards/Main_iPad.storyboard", dest + "Main_iPad.storyboard")

project.add_file_if_doesnt_exist(iosclasses + "AppDelegate.m",group)
project.add_file_if_doesnt_exist(iosclasses + "FlowProtocol.h",group)
project.add_file_if_doesnt_exist(iosclasses + "FlowProtocol.h",group)
project.add_file_if_doesnt_exist(iosclasses + "FlowTissue.h",group)
project.add_file_if_doesnt_exist(iosclasses + "FlowTissue.m",group)

project.add_file_if_doesnt_exist(iosclasses + "DejalActivityView.h",group)
project.add_file_if_doesnt_exist(iosclasses + "DejalActivityView.m",group)

##project.add_file_if_doesnt_exist(iosclasses + "FlowAccellerometer.h",group)
##project.add_file_if_doesnt_exist(iosclasses + "FlowAccellerometer.m",group)
##project.add_file_if_doesnt_exist(iosclasses + "FlowBluetooth.h",group)
##project.add_file_if_doesnt_exist(iosclasses + "FlowBluetooth.m",group)
project.add_file_if_doesnt_exist(iosclasses + "FlowThreadManager.h",group)

project.add_file_if_doesnt_exist(iosclasses + "FlowThreadManager.m",group)
project.add_file_if_doesnt_exist(iosclasses + "ViewController.m",group)
project.add_file_if_doesnt_exist(iosclasses + "ViewController.h",group)

##launch image linkage
project.add_file_if_doesnt_exist(os.path.expandvars("$GOPATH") + "/src/" + sys.argv[1] + "/LaunchImage.png"  , group)


## Until we figure out the framework bitcode compilation...
##if len(sys.argv) > 4 :
print 'Setting bitcode Off!!'
project.add_single_valued_flag('ENABLE_BITCODE', 'NO')
 # else :
	# project.add_single_valued_flag('ENABLE_BITCODE', 'YES')
project.add_file_if_doesnt_exist(os.path.expandvars("$GOPATH") + "/src/" + sys.argv[1] + "/" + sys.argv[3] + ".framework" , group, weak=True)

project.add_file_if_doesnt_exist('System/Library/Frameworks/AVFoundation.framework', parent=frameworks, weak=True, tree='SDKROOT')
project.add_file_if_doesnt_exist('System/Library/Frameworks/CoreLocation.framework', parent=frameworks, weak=True, tree='SDKROOT')

project.save()