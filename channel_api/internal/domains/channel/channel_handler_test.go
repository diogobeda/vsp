package channel

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ChannelRepoMock struct {
	channels                []Channel
	channel                 Channel
	getChannelsError        string
	createChannelError      string
	getChannelByHandleError string
	removeChannelError      string
	updateChannelError      string
}

func (r ChannelRepoMock) GetChannels() ([]Channel, error) {
	if r.getChannelsError != "" {
		return []Channel{}, errors.New(r.getChannelsError)
	}

	return r.channels, nil
}

func (r ChannelRepoMock) CreateChannel(channel Channel) error {
	if r.createChannelError != "" {
		return errors.New(r.createChannelError)
	}

	return nil
}

func (r ChannelRepoMock) GetChannelByHandle(handle string) (Channel, error) {
	if r.getChannelByHandleError != "" {
		return r.channel, errors.New(r.getChannelByHandleError)
	}

	return r.channel, nil
}

func (r ChannelRepoMock) RemoveChannel(id string) error {
	if r.removeChannelError != "" {
		return errors.New(r.removeChannelError)
	}

	return nil
}

func (r ChannelRepoMock) UpdateChannel(channel Channel) error {
	if r.updateChannelError != "" {
		return errors.New(r.updateChannelError)
	}

	return nil
}

func setUpHandler(repo iChannelRepo) *ChannelHandler {
	handler := &ChannelHandler{
		channelRepo: repo,
	}

	return handler
}

func TestGetChannels(t *testing.T) {
	rootRequest := httptest.NewRequest("GET", "/channels", nil)
	channels := make([]Channel, 2)
	successBody, _ := json.Marshal(channels)

	testCases := []struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		channelRepo    *ChannelRepoMock
		expectedStatus int
		expectedBody   []byte
	}{
		{
			w: httptest.NewRecorder(),
			r: rootRequest,
			channelRepo: &ChannelRepoMock{
				channels: channels,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   append(successBody, '\n'),
		},
		{
			w: httptest.NewRecorder(),
			r: rootRequest,
			channelRepo: &ChannelRepoMock{
				getChannelsError: "Error when getting channels",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   []byte("There was a problem when fetching the channels"),
		},
	}

	for _, testCase := range testCases {
		handler := setUpHandler(testCase.channelRepo)
		handler.GetChannels(testCase.w, testCase.r)

		if testCase.w.Code != testCase.expectedStatus {
			t.Errorf("Expect status %v to equal %v", testCase.w.Code, testCase.expectedStatus)
		}

		if !bytes.Equal(testCase.w.Body.Bytes(), testCase.expectedBody) {
			t.Errorf("Expected body \"%+v\" to equal \"%+v\"", testCase.w.Body.String(), string(testCase.expectedBody))
		}
	}
}

func TestCreateChannel(t *testing.T) {
	channelInput := Channel{}
	body, _ := json.Marshal(channelInput)

	testCases := []struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		channelRepo    *ChannelRepoMock
		expectedStatus int
		expectedBody   []byte
	}{
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("POST", "/channels", bytes.NewReader(body)),
			channelRepo: &ChannelRepoMock{
				channel:                 channelInput,
				getChannelByHandleError: "Channel not found",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   []byte(""),
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("POST", "/channels", strings.NewReader("nonsense string")),
			channelRepo: &ChannelRepoMock{
				channel: channelInput,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("There was a problem creating the channel"),
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("POST", "/channels", bytes.NewReader(body)),
			channelRepo: &ChannelRepoMock{
				channel: channelInput,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("A channel with that url already exists"),
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("POST", "/channels", bytes.NewReader(body)),
			channelRepo: &ChannelRepoMock{
				channel:                 channelInput,
				getChannelByHandleError: "Channel not found",
				createChannelError:      "Error when creating channel",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("There was a problem creating the channel"),
		},
	}

	for _, testCase := range testCases {
		handler := setUpHandler(testCase.channelRepo)
		handler.CraetChannel(testCase.w, testCase.r)

		if testCase.w.Code != testCase.expectedStatus {
			t.Errorf("Expect status %v to equal %v", testCase.w.Code, testCase.expectedStatus)
		}

		if !bytes.Equal(testCase.w.Body.Bytes(), testCase.expectedBody) {
			t.Errorf("Expected body \"%+v\" to equal \"%+v\"", testCase.w.Body.String(), string(testCase.expectedBody))
		}
	}
}

func TestGetChannel(t *testing.T) {
	rootRequest := httptest.NewRequest("GET", "/channels/123", nil)
	channel := Channel{}
	successBody, _ := json.Marshal(channel)

	testCases := []struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		channelRepo    *ChannelRepoMock
		expectedStatus int
		expectedBody   []byte
	}{
		{
			w: httptest.NewRecorder(),
			r: rootRequest,
			channelRepo: &ChannelRepoMock{
				channel: channel,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   append(successBody, '\n'),
		},
		{
			w: httptest.NewRecorder(),
			r: rootRequest,
			channelRepo: &ChannelRepoMock{
				getChannelByHandleError: "Channel does not exist",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   []byte("There isn't a channel with that url"),
		},
	}

	for _, testCase := range testCases {
		handler := setUpHandler(testCase.channelRepo)
		handler.GetChannel(testCase.w, testCase.r)

		if testCase.w.Code != testCase.expectedStatus {
			t.Errorf("Expect status %v to equal %v", testCase.w.Code, testCase.expectedStatus)
		}

		if !bytes.Equal(testCase.w.Body.Bytes(), testCase.expectedBody) {
			t.Errorf("Expected body \"%+v\" to equal \"%+v\"", testCase.w.Body.String(), string(testCase.expectedBody))
		}
	}
}

func TestRemoveChannel(t *testing.T) {
	rootRequest := httptest.NewRequest("DELETE", "/channels/123", nil)

	testCases := []struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		channelRepo    *ChannelRepoMock
		expectedStatus int
		expectedBody   []byte
	}{
		{
			w:              httptest.NewRecorder(),
			r:              rootRequest,
			channelRepo:    &ChannelRepoMock{},
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte(""),
		},
		{
			w: httptest.NewRecorder(),
			r: rootRequest,
			channelRepo: &ChannelRepoMock{
				removeChannelError: "Channel does not exist",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("There was a problem when removing the channel"),
		},
	}

	for _, testCase := range testCases {
		handler := setUpHandler(testCase.channelRepo)
		handler.RemoveChannel(testCase.w, testCase.r)

		if testCase.w.Code != testCase.expectedStatus {
			t.Errorf("Expect status %v to equal %v", testCase.w.Code, testCase.expectedStatus)
		}

		if !bytes.Equal(testCase.w.Body.Bytes(), testCase.expectedBody) {
			t.Errorf("Expected body \"%+v\" to equal \"%+v\"", testCase.w.Body.String(), string(testCase.expectedBody))
		}
	}
}

func TestUpdateChannel(t *testing.T) {
	channelInput := Channel{}
	body, _ := json.Marshal(channelInput)

	testCases := []struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		channelRepo    *ChannelRepoMock
		expectedStatus int
		expectedBody   []byte
	}{
		{
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("PUT", "/channels/123", bytes.NewReader(body)),
			channelRepo:    &ChannelRepoMock{},
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(""),
		},
		{
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("PUT", "/channels/123", strings.NewReader("nonsense string")),
			channelRepo:    &ChannelRepoMock{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("There was a problem when updating the channel"),
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("POST", "/channels", bytes.NewReader(body)),
			channelRepo: &ChannelRepoMock{
				updateChannelError: "Error updating the channel",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("There was a problem when updating the channel"),
		},
	}

	for _, testCase := range testCases {
		handler := setUpHandler(testCase.channelRepo)
		handler.UpdateChannel(testCase.w, testCase.r)

		if testCase.w.Code != testCase.expectedStatus {
			t.Errorf("Expect status %v to equal %v", testCase.w.Code, testCase.expectedStatus)
		}

		if !bytes.Equal(testCase.w.Body.Bytes(), testCase.expectedBody) {
			t.Errorf("Expected body \"%+v\" to equal \"%+v\"", testCase.w.Body.String(), string(testCase.expectedBody))
		}
	}
}
