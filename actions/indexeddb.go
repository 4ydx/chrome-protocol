package actions

import (
	"github.com/4ydx/cdp/protocol/indexeddb"
	"github.com/4ydx/chrome-protocol"
	"log"
	"time"
)

// RequestDatabaseNames returns a list of the databases controlled by the given security origin.
func RequestDatabaseNames(frame *cdp.Frame, securityOrigin string, timeout time.Duration) (*indexeddb.RequestDatabaseNamesReply, error) {
	action := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: indexeddb.CommandIndexedDBRequestDatabaseNames, Params: &indexeddb.RequestDatabaseNamesArgs{SecurityOrigin: securityOrigin}, Reply: &indexeddb.RequestDatabaseNamesReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*indexeddb.RequestDatabaseNamesReply), nil
}

// RequestDatabase returns the database with object stores.
func RequestDatabase(frame *cdp.Frame, securityOrigin, databaseName string, timeout time.Duration) (*indexeddb.RequestDatabaseReply, error) {
	action := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: indexeddb.CommandIndexedDBRequestDatabase, Params: &indexeddb.RequestDatabaseArgs{SecurityOrigin: securityOrigin, DatabaseName: databaseName}, Reply: &indexeddb.RequestDatabaseReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*indexeddb.RequestDatabaseReply), nil
}

// RequestData returns the requested data.
func RequestData(frame *cdp.Frame, securityOrigin, databaseName, objectStoreName, indexName string, skipCount, pageSize int, keyRange *indexeddb.KeyRange, timeout time.Duration) (*indexeddb.RequestDataReply, error) {
	args := &indexeddb.RequestDataArgs{
		SecurityOrigin:  securityOrigin,
		DatabaseName:    databaseName,
		ObjectStoreName: objectStoreName,
		IndexName:       indexName,
		SkipCount:       skipCount,
		PageSize:        pageSize,
		KeyRange:        keyRange,
	}
	action := cdp.NewAction(frame,
		[]cdp.Event{},
		[]cdp.Command{
			cdp.Command{ID: frame.RequestID.GetNext(), Method: indexeddb.CommandIndexedDBRequestData, Params: args, Reply: &indexeddb.RequestDataReply{}, Timeout: timeout},
		})
	err := action.Run(frame)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return action.Commands[0].Reply.(*indexeddb.RequestDataReply), nil
}
