// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	redisearch "github.com/RediSearch/redisearch-go/redisearch"
)

// RedisSearchClient is an autogenerated mock type for the RedisSearchClient type
type RedisSearchClient struct {
	mock.Mock
}

// Delete provides a mock function with given fields: docID
func (_m *RedisSearchClient) Delete(docID string) error {
	ret := _m.Called(docID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(docID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Index provides a mock function with given fields: docs
func (_m *RedisSearchClient) Index(docs ...redisearch.Document) error {
	_va := make([]interface{}, len(docs))
	for _i := range docs {
		_va[_i] = docs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...redisearch.Document) error); ok {
		r0 = rf(docs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: text, offset, limit
func (_m *RedisSearchClient) Search(text string, offset int, limit int) ([]redisearch.Document, int, error) {
	ret := _m.Called(text, offset, limit)

	var r0 []redisearch.Document
	if rf, ok := ret.Get(0).(func(string, int, int) []redisearch.Document); ok {
		r0 = rf(text, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]redisearch.Document)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(string, int, int) int); ok {
		r1 = rf(text, offset, limit)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, int, int) error); ok {
		r2 = rf(text, offset, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
