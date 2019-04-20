// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import wire "gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"

// Committee is an autogenerated mock type for the Committee type
type Committee struct {
	mock.Mock
}

// GetVotingCommittee provides a mock function with given fields: _a0, _a1
func (_m *Committee) GetVotingCommittee(_a0 uint64, _a1 uint8) (map[string]uint8, error) {
	ret := _m.Called(_a0, _a1)

	var r0 map[string]uint8
	if rf, ok := ret.Get(0).(func(uint64, uint8) map[string]uint8); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]uint8)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint8) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsMember provides a mock function with given fields: _a0
func (_m *Committee) IsMember(_a0 []byte) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Priority provides a mock function with given fields: _a0, _a1
func (_m *Committee) Priority(_a0 wire.Event, _a1 wire.Event) wire.Event {
	ret := _m.Called(_a0, _a1)

	var r0 wire.Event
	if rf, ok := ret.Get(0).(func(wire.Event, wire.Event) wire.Event); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(wire.Event)
		}
	}

	return r0
}

// Quorum provides a mock function with given fields:
func (_m *Committee) Quorum() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// VerifyVoteSet provides a mock function with given fields: voteSet, hash, round, step
func (_m *Committee) VerifyVoteSet(voteSet []wire.Event, hash []byte, round uint64, step uint8) error {
	ret := _m.Called(voteSet, hash, round, step)

	var r0 error
	if rf, ok := ret.Get(0).(func([]wire.Event, []byte, uint64, uint8) error); ok {
		r0 = rf(voteSet, hash, round, step)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
