// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_gcs

import (
	fmt "fmt"
	gcs "github.com/jacobsa/gcloud/gcs"
	oglemock "github.com/jacobsa/oglemock"
	context "golang.org/x/net/context"
	storage "google.golang.org/cloud/storage"
	io "io"
	runtime "runtime"
	unsafe "unsafe"
)

type MockBucket interface {
	gcs.Bucket
	oglemock.MockObject
}

type mockBucket struct {
	controller  oglemock.Controller
	description string
}

func NewMockBucket(
	c oglemock.Controller,
	desc string) MockBucket {
	return &mockBucket{
		controller:  c,
		description: desc,
	}
}

func (m *mockBucket) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockBucket) Oglemock_Description() string {
	return m.description
}

func (m *mockBucket) CreateObject(p0 context.Context, p1 *gcs.CreateObjectRequest) (o0 *storage.Object, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"CreateObject",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.CreateObject: invalid return values: %v", retVals))
	}

	// o0 *storage.Object
	if retVals[0] != nil {
		o0 = retVals[0].(*storage.Object)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockBucket) DeleteObject(p0 context.Context, p1 string) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"DeleteObject",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockBucket.DeleteObject: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockBucket) ListObjects(p0 context.Context, p1 *storage.Query) (o0 *storage.Objects, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"ListObjects",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.ListObjects: invalid return values: %v", retVals))
	}

	// o0 *storage.Objects
	if retVals[0] != nil {
		o0 = retVals[0].(*storage.Objects)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockBucket) Name() (o0 string) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Name",
		file,
		line,
		[]interface{}{})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockBucket.Name: invalid return values: %v", retVals))
	}

	// o0 string
	if retVals[0] != nil {
		o0 = retVals[0].(string)
	}

	return
}

func (m *mockBucket) NewReader(p0 context.Context, p1 string) (o0 io.ReadCloser, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"NewReader",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.NewReader: invalid return values: %v", retVals))
	}

	// o0 io.ReadCloser
	if retVals[0] != nil {
		o0 = retVals[0].(io.ReadCloser)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockBucket) StatObject(p0 context.Context, p1 *gcs.StatObjectRequest) (o0 *storage.Object, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"StatObject",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.StatObject: invalid return values: %v", retVals))
	}

	// o0 *storage.Object
	if retVals[0] != nil {
		o0 = retVals[0].(*storage.Object)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockBucket) UpdateObject(p0 context.Context, p1 *gcs.UpdateObjectRequest) (o0 *storage.Object, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"UpdateObject",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.UpdateObject: invalid return values: %v", retVals))
	}

	// o0 *storage.Object
	if retVals[0] != nil {
		o0 = retVals[0].(*storage.Object)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}