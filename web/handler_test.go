package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatedResponse(t *testing.T) {
	handler := &WebHandler{}
	writer := httptest.NewRecorder()
	handler.Created(writer)

	if writer.Code != http.StatusCreated {
		t.Errorf("Expected status %v to equal %v", writer.Code, http.StatusCreated)
	}
}

func TestOKResponse(t *testing.T) {
	handler := &WebHandler{}
	writer := httptest.NewRecorder()
	handler.Ok(writer)

	if writer.Code != http.StatusOK {
		t.Errorf("Expected status %v to equal %v", writer.Code, http.StatusOK)
	}
}

func TestBadRequestResponse(t *testing.T) {
	handler := &WebHandler{}
	writer := httptest.NewRecorder()
	expectedString := "Bad request error"
	handler.BadRequest(writer, expectedString)

	if writer.Code != http.StatusBadRequest {
		t.Errorf("Expected status %v to equal %v", writer.Code, http.StatusBadRequest)
	}

	if !bytes.Equal([]byte(expectedString), writer.Body.Bytes()) {
		t.Errorf("Expected body \"%v\" to equal \"%v\"", writer.Body.String(), expectedString)
	}
}

func TestInternalResponse(t *testing.T) {
	handler := &WebHandler{}
	writer := httptest.NewRecorder()
	expectedString := "Internal error"
	handler.Internal(writer, expectedString)

	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %v to equal %v", writer.Code, http.StatusInternalServerError)
	}

	if !bytes.Equal([]byte(expectedString), writer.Body.Bytes()) {
		t.Errorf("Expected body \"%v\" to equal \"%v\"", writer.Body.String(), expectedString)
	}
}

func TestNotFoundResponse(t *testing.T) {
	handler := &WebHandler{}
	writer := httptest.NewRecorder()
	expectedString := "Not found"
	handler.NotFound(writer, expectedString)

	if writer.Code != http.StatusNotFound {
		t.Errorf("Expected status %v to equal %v", writer.Code, http.StatusNotFound)
	}

	if !bytes.Equal([]byte(expectedString), writer.Body.Bytes()) {
		t.Errorf("Expected body \"%v\" to equal \"%v\"", writer.Body.String(), expectedString)
	}
}
