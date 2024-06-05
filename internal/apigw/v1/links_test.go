package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/pb"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockLinksClient struct {
	mock.Mock
}

func (m *MockLinksClient) CreateLink(ctx context.Context, in *pb.CreateLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Empty), args.Error(1)
}

func (m *MockLinksClient) GetLink(ctx context.Context, in *pb.GetLinkRequest, opts ...grpc.CallOption) (*pb.Link, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Link), args.Error(1)
}

func (m *MockLinksClient) GetLinkByUserID(ctx context.Context, in *pb.GetLinksByUserId, opts ...grpc.CallOption) (*pb.ListLinkResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ListLinkResponse), args.Error(1)
}

func (m *MockLinksClient) UpdateLink(ctx context.Context, in *pb.UpdateLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Empty), args.Error(1)
}

func (m *MockLinksClient) DeleteLink(ctx context.Context, in *pb.DeleteLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Empty), args.Error(1)
}

func (m *MockLinksClient) ListLinks(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.ListLinkResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ListLinkResponse), args.Error(1)
}

func TestGetLinks(t *testing.T) {
	mockLinksClient := new(MockLinksClient)

	links := make([]*pb.Link, 0, 2)
	links = append(
		links, &pb.Link{
			CreatedAt: time.Now().String(),
			Id:        "test_id_1",
			Title:     "test_url_1",
			UpdatedAt: time.Now().String(),
			Url:       "test_url_1",
			UserId:    "test_user_1",
		},
	)
	links = append(
		links, &pb.Link{
			CreatedAt: time.Now().String(),
			Id:        "test_id_2",
			Title:     "test_url_2",
			UpdatedAt: time.Now().String(),
			Url:       "test_url_2",
			UserId:    "test_user_2",
		},
	)
	linkResponse := new(pb.ListLinkResponse)
	linkResponse.Links = links

	mockLinksClient.On("ListLinks", mock.Anything, mock.Anything).Return(linkResponse, nil)

	handler := New(nil, mockLinksClient)

	req, err := http.NewRequest("GET", "/v1/links", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler.GetLinks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var respLinks []apiv1.Link
	err = json.NewDecoder(rr.Body).Decode(&respLinks)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(respLinks))
}

func TestDeleteLinksId(t *testing.T) {
	mockLinksClient := new(MockLinksClient)
	handler := New(nil, mockLinksClient)

	testId := "test_id"
	mockLinksClient.On("DeleteLink", mock.Anything, mock.Anything).Return(new(pb.Empty), nil).Once()

	req, err := http.NewRequest("DELETE", "/v1/links/"+testId, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.DeleteLinksId(rr, req, testId)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	mockLinksClient.AssertCalled(t, "DeleteLink", mock.Anything, mock.Anything)
}

func TestPostLinks(t *testing.T) {
	mockLinksClient := new(MockLinksClient)
	handler := New(nil, mockLinksClient)

	mockLinksClient.On("CreateLink", mock.Anything, mock.Anything).Return(new(pb.Empty), nil).Once()

	requestBody := []byte(`{"id": "test_id", "user_id": "test_user", "url": "test_url", "title": "test_title"}`)

	req, err := http.NewRequest("POST", "/v1/links", bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.PostLinks(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	mockLinksClient.AssertCalled(t, "CreateLink", mock.Anything, mock.Anything)
}

func TestPutLinksId(t *testing.T) {
	mockLinksClient := new(MockLinksClient)
	handler := New(nil, mockLinksClient)

	mockLinksClient.On("UpdateLink", mock.Anything, mock.Anything).Return(new(pb.Empty), nil).Once()

	requestBody := []byte(`{"id": "test_id", "user_id": "test_user", "url": "test_url", "title": "test_title"}`)

	testId := "test_id"
	req, err := http.NewRequest("PUT", "/v1/links"+testId, bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.PutLinksId(rr, req, testId)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	mockLinksClient.AssertCalled(t, "UpdateLink", mock.Anything, mock.Anything)

}
