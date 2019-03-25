//Package watchdir does watch directories
// +build windows
package watchdir

import (
	"log"
	"syscall"
	"unsafe"
)

//ModificationEvent is an event notification
type ModificationEvent struct {
	watcher  int
	filePath string
}

//StartWatching returns a EventNotification when triggered
func StartWatching(id int, path string, events chan ModificationEvent) {
	//Pointer to utf8 string
	sptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		log.Fatalln("UTF16PtrFromString failed: ", err, syscall.GetLastError())
	}
	//reference to a file or a directory to watch
	fhnd, err := syscall.CreateFile(
		sptr,
		syscall.FILE_LIST_DIRECTORY, //dwDesiredAccess
		syscall.FILE_SHARE_DELETE|syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, //dwShareMode
		nil,                   //lpSecurityAttributes
		syscall.OPEN_EXISTING, //dwCreationDisposition
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OVERLAPPED, //dwFlagsAndAttributes
		0)
	if err != nil {
		log.Fatalln("CreteFile failed: ", err, syscall.GetLastError())
	}

	log.Println("Getting information on: ", path)

	var (
		prt, _                       = syscall.CreateIoCompletionPort(syscall.InvalidHandle, 0, 0, 0)
		fileprt, _                   = syscall.CreateIoCompletionPort(fhnd, prt, 0, 0)
		resultbfr                    = make([]byte, 1, 1024*2)
		ovrlpd                       syscall.Overlapped
		ovrlpdptr                    *syscall.Overlapped
		bytesRead, resultOffset, key uint32
	)
	for {

		syscall.ReadDirectoryChanges(fhnd, //hDirectory
			&resultbfr[0],  //lpBuffer
			uint32(1024*2), //nBufferLength
			true,           //bWatchSubtree
			syscall.FILE_NOTIFY_CHANGE_DIR_NAME|syscall.FILE_NOTIFY_CHANGE_FILE_NAME|syscall.FILE_NOTIFY_CHANGE_LAST_WRITE|syscall.FILE_NOTIFY_CHANGE_SIZE, //dwNotifyFilter
			nil, &ovrlpd, 0)

		err = syscall.GetQueuedCompletionStatus(fileprt, &bytesRead, &key, &ovrlpdptr, syscall.INFINITE)

		event := (*syscall.FileNotifyInformation)(unsafe.Pointer(&resultbfr[resultOffset]))
		buf := (*[128]uint16)(unsafe.Pointer(&event.FileName))
		name := syscall.UTF16ToString(buf[:event.FileNameLength/2])
		events <- ModificationEvent{watcher: id, filePath: name}
	}
}
