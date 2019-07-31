## Mon 06 Mar 2017 04:15:27 PM MST Version: 0.0.1
Initial commit

Multiple directory even monitoring works

## Wed 08 Mar 2017 04:44:43 PM MST Version: 0.0.2
Comments and code cleanup

## Fri 10 Mar 2017 06:39:59 PM MST Version: 0.0.3
Make RecursiveAdd() actually work

Refactor and cleanup

## Fri 10 Mar 2017 08:29:10 PM MST Version: 0.0.4
Add ListDescriptors() and RemoveDescriptor()

Update example to showcase the proper use of these functions

## Fri 10 Mar 2017 09:06:36 PM MST Version: 0.0.5
Forgot to add mutex locks in ListDescriptors() and RemoveDescriptor()

## Sat 11 Mar 2017 01:31:32 PM MST Version: 0.0.6
Forgot to add mutex lock to DescriptorExists()

Fix comment typo

## Sat 11 Mar 2017 03:34:30 PM MST Version: 0.0.7
Rename getWatchDescriptor() to GetDescriptorByWatch()

## Sat 11 Mar 2017 03:37:41 PM MST Version: 0.0.8
Added GetDescriptorByPath()

## Sat 11 Mar 2017 03:49:17 PM MST Version: 0.0.9

Added Stop() for stopping running watch descriptors

Added d.Running to check a descriptor's status

Added status checks to d.Stop() and d.Start()

## Sat 11 Mar 2017 04:08:11 PM MST Version: 0.0.10
Added w.StopAll() to stop all currently running WatchDescriptors

## Sun 12 Mar 2017 03:02:50 PM MDT Version: 0.0.11

Added WatchDescriptor.DoesPathExist() that returns true if a descriptor's path
exists, false otherwise

Fixed Watcher.RemoveDescriptor() to not try to remove an inotify watch
of a file that has already been deleted, since inotify removes watches
itself. So we just need to handle our own bookkeeping.

## Sun 12 Mar 2017 03:17:26 PM MDT Version: 0.0.12

Refactor GetDescriptorByPath() to be a little less dumb

## Sun 12 Mar 2017 03:35:42 PM MDT Version: 0.0.13

Added Name string to FsEvent structure. This is the name of the actual file
associated with an event.

## Sun 12 Mar 2017 05:46:32 PM MDT Version: 0.0.14

Updated example/example.go

Added some new higher level functions for getting event types:

IsDirEvent(): Returns true if the event happened to a directory

IsDirCreated(): Returns true if the event was a directory creation or
                a directory was moved into the root watch directory

IsDirRemoved(): Returns true if the event was a direction deletion or
                a directory was moved outside of the root watch directory

IsFileCreated(): Returns true if the event was a file being created or moved
                 into the root watch directory

IsFileRemoved(): Returns true if the event was a file being deleted or moved
                 outside of the root watch directory

## Sun 12 Mar 2017 06:46:48 PM MDT Version: 0.0.15

Refactored event mask checking in helper functions using new CheckMask() function

Added IsRootMoved() and IsRootDeletion() helper functions

Added custom predefined event flags that work with the new event checking functions

## Thu 30 Mar 2017 11:42:58 AM MDT Version: 0.0.16

Fixed fmt.Errorf() typo in fsevents.go

## Thu 30 Mar 2017 02:58:10 PM MDT Version: 0.0.17
Add check to AddDescriptor() to make sure dirPath exists
before trying to create a watch for it

## Thu 30 Mar 2017 03:33:31 PM MDT Version: 0.0.18
Renamed/removed/added some error variables

## Thu 30 Mar 2017 03:56:37 PM MDT Version: 0.0.19
Make IsRootDeletion() and IsRootMoved() compare the pathname of the event, and the root
path of the watcher

## Thu 20 Apr 2017 04:07:58 PM MDT Version: 0.0.20

Added fsevents_test.go

## Thu 20 Apr 2017 04:20:18 PM MDT Version: 0.0.21

Formatting in fsevents.go
Logging output in example/example.go

## Sat 29 Apr 2017 09:52:43 AM MDT Version: 0.0.22
Return error from w.StopAll()

## Sat 29 Apr 2017 10:16:11 AM MDT Version: 0.0.23
Added scripts/test/cover.sh and updated .gitignore

## Mon 18 Mar 2019 01:31:09 PM MDT Version: 0.0.24
Added Event counter and helper functions. Simplified Watch().

Added ReadSingleEvent() with simplified event reading from inotify descriptor

## Sat 23 Mar 2019 10:35:33 AM MDT Version: 0.0.25
Added tests for (most) event masks

## Sat 23 Mar 2019 05:54:24 PM MDT Version: 0.0.26
Updated mask variables to uint32 to avoid conversion

Added EventHandle interface, RegisterEventHandle UnregisterEventHandle, getEventHandle and WatchAndHandle.

Updated example code to show new EventHandle interface usage

## Sun 24 Mar 2019 05:42:51 PM MDT Version: 0.0.27
Refactored Stop and Start WatchDescriptor methods to not require the inotify file descriptor from the Watcher

Added ID field to FsEvent to allow for differentiation between events and discerning what event came when.

Split example code into examples/handle and examples/loop to show different ways to use go-fsevents

## Sun 24 Mar 2019 06:02:03 PM MDT Version: 0.0.28

Added time.Time timestamp to FsEvent object

## Sun 24 Mar 2019 06:05:54 PM MDT Version: 0.0.29

Changed FsEvent timestamp to UTC instead of local

## Mon 25 Mar 2019 09:03:07 PM MDT Version: 0.0.30

Make UnregisterEventHandle return an error if we try to remove a handler that doesn't exit
