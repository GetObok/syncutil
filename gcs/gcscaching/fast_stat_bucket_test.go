// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcscaching_test

import (
	"errors"
	"testing"
	"time"

	"github.com/googlecloudplatform/gcsfuse/timeutil"
	"github.com/jacobsa/gcloud/gcs"
	"github.com/jacobsa/gcloud/gcs/gcscaching"
	"github.com/jacobsa/gcloud/gcs/gcscaching/mock_gcscaching"
	"github.com/jacobsa/gcloud/gcs/mock_gcs"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
)

func TestFastStatBucket(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////

const ttl = time.Second

type fastStatBucketTest struct {
	cache   mock_gcscaching.MockStatCache
	clock   timeutil.SimulatedClock
	wrapped mock_gcs.MockBucket

	bucket gcs.Bucket
}

func (t *fastStatBucketTest) SetUp(ti *TestInfo) {
	// Set up a fixed, non-zero time.
	t.clock.SetTime(time.Date(2015, 4, 5, 2, 15, 0, 0, time.Local))

	// Set up dependencies.
	t.cache = mock_gcscaching.NewMockStatCache(ti.MockController, "cache")
	t.wrapped = mock_gcs.NewMockBucket(ti.MockController, "wrapped")

	t.bucket = gcscaching.NewFastStatBucket(
		ttl,
		t.cache,
		&t.clock,
		t.wrapped)
}

////////////////////////////////////////////////////////////////////////
// CreateObject
////////////////////////////////////////////////////////////////////////

type CreateObjectTest struct {
	fastStatBucketTest
}

func init() { RegisterTestSuite(&CreateObjectTest{}) }

func (t *CreateObjectTest) CallsEraseAndWrapped() {
	const name = "taco"

	// Erase
	ExpectCall(t.cache, "Erase")(name)

	// Wrapped
	var wrappedReq *gcs.CreateObjectRequest
	ExpectCall(t.wrapped, "CreateObject")(Any(), Any()).
		WillOnce(DoAll(SaveArg(1, &wrappedReq), Return(nil, errors.New(""))))

	// Call
	req := &gcs.CreateObjectRequest{
		Name: name,
	}

	_, _ = t.bucket.CreateObject(nil, req)

	AssertNe(nil, wrappedReq)
	ExpectEq(req, wrappedReq)
}

func (t *CreateObjectTest) WrappedFails() {
	const name = ""
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	ExpectCall(t.wrapped, "CreateObject")(Any(), Any()).
		WillOnce(Return(nil, errors.New("taco")))

	// Call
	_, err = t.bucket.CreateObject(nil, &gcs.CreateObjectRequest{})

	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *CreateObjectTest) WrappedSucceeds() {
	const name = "taco"
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	obj := &gcs.Object{
		Name:       name,
		Generation: 1234,
	}

	ExpectCall(t.wrapped, "CreateObject")(Any(), Any()).
		WillOnce(Return(obj, nil))

	// Insert
	ExpectCall(t.cache, "Insert")(obj, timeutil.TimeEq(t.clock.Now().Add(ttl)))

	// Call
	o, err := t.bucket.CreateObject(nil, &gcs.CreateObjectRequest{})

	AssertEq(nil, err)
	ExpectEq(obj, o)
}

////////////////////////////////////////////////////////////////////////
// StatObject
////////////////////////////////////////////////////////////////////////

type StatObjectTest struct {
	fastStatBucketTest
}

func init() { RegisterTestSuite(&StatObjectTest{}) }

func (t *StatObjectTest) CallsCache() {
	const name = "taco"

	// LookUp
	ExpectCall(t.cache, "LookUp")(name, timeutil.TimeEq(t.clock.Now())).
		WillOnce(Return(true, &gcs.Object{}))

	// Call
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	_, _ = t.bucket.StatObject(nil, req)
}

func (t *StatObjectTest) CacheHit_Positive() {
	const name = "taco"

	// LookUp
	obj := &gcs.Object{
		Name: name,
	}

	ExpectCall(t.cache, "LookUp")(Any(), Any()).
		WillOnce(Return(true, obj))

	// Call
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	o, err := t.bucket.StatObject(nil, req)
	AssertEq(nil, err)
	ExpectEq(obj, o)
}

func (t *StatObjectTest) CacheHit_Negative() {
	const name = "taco"

	// LookUp
	ExpectCall(t.cache, "LookUp")(Any(), Any()).
		WillOnce(Return(true, nil))

	// Call
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	_, err := t.bucket.StatObject(nil, req)
	ExpectThat(err, HasSameTypeAs(&gcs.NotFoundError{}))
}

func (t *StatObjectTest) CallsWrapped() {
	const name = ""
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	// LookUp
	ExpectCall(t.cache, "LookUp")(Any(), Any()).
		WillOnce(Return(false, nil))

	// Wrapped
	ExpectCall(t.wrapped, "StatObject")(Any(), req).
		WillOnce(Return(nil, errors.New("")))

	// Call
	_, _ = t.bucket.StatObject(nil, req)
}

func (t *StatObjectTest) WrappedFails() {
	const name = ""

	// LookUp
	ExpectCall(t.cache, "LookUp")(Any(), Any()).
		WillOnce(Return(false, nil))

	// Wrapped
	ExpectCall(t.wrapped, "StatObject")(Any(), Any()).
		WillOnce(Return(nil, errors.New("taco")))

	// Call
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	_, err := t.bucket.StatObject(nil, req)
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *StatObjectTest) WrappedSaysNotFound() {
	AssertTrue(false, "TODO")
}

func (t *StatObjectTest) WrappedSucceeds() {
	const name = "taco"

	// LookUp
	ExpectCall(t.cache, "LookUp")(Any(), Any()).
		WillOnce(Return(false, nil))

	// Wrapped
	obj := &gcs.Object{
		Name: name,
	}

	ExpectCall(t.wrapped, "StatObject")(Any(), Any()).
		WillOnce(Return(obj, nil))

	// Insert
	ExpectCall(t.cache, "Insert")(obj, timeutil.TimeEq(t.clock.Now().Add(ttl)))

	// Call
	req := &gcs.StatObjectRequest{
		Name: name,
	}

	o, err := t.bucket.StatObject(nil, req)
	AssertEq(nil, err)
	ExpectEq(obj, o)
}

////////////////////////////////////////////////////////////////////////
// ListObjects
////////////////////////////////////////////////////////////////////////

type ListObjectsTest struct {
	fastStatBucketTest
}

func init() { RegisterTestSuite(&ListObjectsTest{}) }

func (t *ListObjectsTest) WrappedFails() {
	// Wrapped
	ExpectCall(t.wrapped, "ListObjects")(Any(), Any()).
		WillOnce(Return(nil, errors.New("taco")))

	// Call
	_, err := t.bucket.ListObjects(nil, &gcs.ListObjectsRequest{})
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ListObjectsTest) EmptyListing() {
	// Wrapped
	expected := &gcs.Listing{}

	ExpectCall(t.wrapped, "ListObjects")(Any(), Any()).
		WillOnce(Return(expected, nil))

	// Call
	listing, err := t.bucket.ListObjects(nil, &gcs.ListObjectsRequest{})

	AssertEq(nil, err)
	ExpectEq(expected, listing)
}

func (t *ListObjectsTest) NonEmptyListing() {
	// Wrapped
	o0 := &gcs.Object{Name: "taco"}
	o1 := &gcs.Object{Name: "burrito"}

	expected := &gcs.Listing{
		Objects: []*gcs.Object{o0, o1},
	}

	ExpectCall(t.wrapped, "ListObjects")(Any(), Any()).
		WillOnce(Return(expected, nil))

	// Insert
	ExpectCall(t.cache, "Insert")(o0, timeutil.TimeEq(t.clock.Now().Add(ttl)))
	ExpectCall(t.cache, "Insert")(o1, timeutil.TimeEq(t.clock.Now().Add(ttl)))

	// Call
	listing, err := t.bucket.ListObjects(nil, &gcs.ListObjectsRequest{})

	AssertEq(nil, err)
	ExpectEq(expected, listing)
}

////////////////////////////////////////////////////////////////////////
// UpdateObject
////////////////////////////////////////////////////////////////////////

type UpdateObjectTest struct {
	fastStatBucketTest
}

func init() { RegisterTestSuite(&UpdateObjectTest{}) }

func (t *UpdateObjectTest) CallsEraseAndWrapped() {
	const name = "taco"

	// Erase
	ExpectCall(t.cache, "Erase")(name)

	// Wrapped
	var wrappedReq *gcs.UpdateObjectRequest
	ExpectCall(t.wrapped, "UpdateObject")(Any(), Any()).
		WillOnce(DoAll(SaveArg(1, &wrappedReq), Return(nil, errors.New(""))))

	// Call
	req := &gcs.UpdateObjectRequest{
		Name: name,
	}

	_, _ = t.bucket.UpdateObject(nil, req)

	AssertNe(nil, wrappedReq)
	ExpectEq(req, wrappedReq)
}

func (t *UpdateObjectTest) WrappedFails() {
	const name = ""
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	ExpectCall(t.wrapped, "UpdateObject")(Any(), Any()).
		WillOnce(Return(nil, errors.New("taco")))

	// Call
	_, err = t.bucket.UpdateObject(nil, &gcs.UpdateObjectRequest{})

	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *UpdateObjectTest) WrappedSucceeds() {
	const name = "taco"
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	obj := &gcs.Object{
		Name:       name,
		Generation: 1234,
	}

	ExpectCall(t.wrapped, "UpdateObject")(Any(), Any()).
		WillOnce(Return(obj, nil))

	// Insert
	ExpectCall(t.cache, "Insert")(obj, timeutil.TimeEq(t.clock.Now().Add(ttl)))

	// Call
	o, err := t.bucket.UpdateObject(nil, &gcs.UpdateObjectRequest{})

	AssertEq(nil, err)
	ExpectEq(obj, o)
}

////////////////////////////////////////////////////////////////////////
// DeleteObject
////////////////////////////////////////////////////////////////////////

type DeleteObjectTest struct {
	fastStatBucketTest
}

func init() { RegisterTestSuite(&DeleteObjectTest{}) }

func (t *DeleteObjectTest) CallsEraseAndWrapped() {
	const name = "taco"

	// Erase
	ExpectCall(t.cache, "Erase")(name)

	// Wrapped
	ExpectCall(t.wrapped, "DeleteObject")(Any(), name).
		WillOnce(Return(errors.New("")))

	// Call
	_ = t.bucket.DeleteObject(nil, name)
}

func (t *DeleteObjectTest) WrappedFails() {
	const name = ""
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	ExpectCall(t.wrapped, "DeleteObject")(Any(), Any()).
		WillOnce(Return(errors.New("taco")))

	// Call
	err = t.bucket.DeleteObject(nil, name)

	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *DeleteObjectTest) WrappedSucceeds() {
	const name = ""
	var err error

	// Erase
	ExpectCall(t.cache, "Erase")(Any())

	// Wrapped
	ExpectCall(t.wrapped, "DeleteObject")(Any(), Any()).
		WillOnce(Return(nil))

	// Call
	err = t.bucket.DeleteObject(nil, name)
	AssertEq(nil, err)
}
