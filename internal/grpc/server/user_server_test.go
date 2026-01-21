package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/yeegeek/uyou-go-api-starter/api/proto/user"
	"github.com/yeegeek/uyou-go-api-starter/internal/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockUserService Mock 用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, req user.RegisterRequest) (*user.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, req user.LoginRequest) (*user.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uint, req user.UpdateUserRequest) (*user.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, filters user.UserFilterParams, page, perPage int) ([]user.User, int64, error) {
	args := m.Called(ctx, filters, page, perPage)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) PromoteToAdmin(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockUserRepository Mock 用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, usr *user.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, usr *user.User) error {
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ListAllUsers(ctx context.Context, filters user.UserFilterParams, page, perPage int) ([]user.User, int64, error) {
	args := m.Called(ctx, filters, page, perPage)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) AssignRole(ctx context.Context, userID uint, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveRole(ctx context.Context, userID uint, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockUserRepository) FindRoleByName(ctx context.Context, name string) (*user.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Role), args.Error(1)
}

func (m *MockUserRepository) GetUserRoles(ctx context.Context, userID uint) ([]user.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.Role), args.Error(1)
}

func (m *MockUserRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestNewUserServiceServer(t *testing.T) {
	mockService := new(MockUserService)
	mockRepo := new(MockUserRepository)

	server := NewUserServiceServer(mockService, mockRepo)
	assert.NotNil(t, server)
	assert.Equal(t, mockService, server.userService)
	assert.Equal(t, mockRepo, server.userRepo)
}

func TestUserServiceServer_GetUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.GetUserRequest
		setupMock   func(*MockUserService)
		wantErr     bool
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name: "successful get user",
			req:  &pb.GetUserRequest{Id: 1},
			setupMock: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(&user.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
					Roles: []user.Role{{Name: "user"}},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req:  &pb.GetUserRequest{Id: 999},
			setupMock: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			wantErr:     true,
			wantCode:    codes.NotFound,
			wantMessage: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockService)

			server := NewUserServiceServer(mockService, mockRepo)
			resp, err := server.GetUser(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, uint32(1), resp.User.Id)
				assert.Equal(t, "Test User", resp.User.Name)
				assert.Equal(t, "test@example.com", resp.User.Email)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.GetUserByEmailRequest
		setupMock   func(*MockUserRepository)
		wantErr     bool
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name: "successful get user by email",
			req:  &pb.GetUserByEmailRequest{Email: "test@example.com"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(&user.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
					Roles: []user.Role{{Name: "user"}},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req:  &pb.GetUserByEmailRequest{Email: "notfound@example.com"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("user not found"))
			},
			wantErr:     true,
			wantCode:    codes.NotFound,
			wantMessage: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			server := NewUserServiceServer(mockService, mockRepo)
			resp, err := server.GetUserByEmail(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, "test@example.com", resp.User.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_ListUsers(t *testing.T) {
	tests := []struct {
		name      string
		req       *pb.ListUsersRequest
		setupMock func(*MockUserService)
		wantErr   bool
		wantTotal int32
	}{
		{
			name: "successful list users",
			req: &pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(m *MockUserService) {
				users := []user.User{
					{ID: 1, Name: "User 1", Email: "user1@example.com"},
					{ID: 2, Name: "User 2", Email: "user2@example.com"},
				}
				m.On("ListUsers", mock.Anything, user.UserFilterParams{}, 1, 10).Return(users, int64(2), nil)
			},
			wantErr:   false,
			wantTotal: 2,
		},
		{
			name: "default pagination",
			req:  &pb.ListUsersRequest{},
			setupMock: func(m *MockUserService) {
				users := []user.User{}
				m.On("ListUsers", mock.Anything, user.UserFilterParams{}, 1, 10).Return(users, int64(0), nil)
			},
			wantErr:   false,
			wantTotal: 0,
		},
		{
			name: "service error",
			req: &pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, user.UserFilterParams{}, 1, 10).Return(nil, int64(0), errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockService)

			server := NewUserServiceServer(mockService, mockRepo)
			resp, err := server.ListUsers(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantTotal, resp.Total)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_UpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.UpdateUserRequest
		setupMock   func(*MockUserService)
		wantErr     bool
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name: "successful update user",
			req: &pb.UpdateUserRequest{
				Id:    1,
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, uint(1), user.UpdateUserRequest{
					Name:  "Updated Name",
					Email: "updated@example.com",
				}).Return(&user.User{
					ID:    1,
					Name:  "Updated Name",
					Email: "updated@example.com",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "update fails",
			req: &pb.UpdateUserRequest{
				Id:    999,
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, uint(999), mock.Anything).Return(nil, errors.New("update failed"))
			},
			wantErr:     true,
			wantCode:    codes.Internal,
			wantMessage: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockService)

			server := NewUserServiceServer(mockService, mockRepo)
			resp, err := server.UpdateUser(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, "Updated Name", resp.User.Name)
				assert.Equal(t, "updated@example.com", resp.User.Email)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserServiceServer_DeleteUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.DeleteUserRequest
		setupMock   func(*MockUserService)
		wantErr     bool
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name: "successful delete user",
			req:  &pb.DeleteUserRequest{Id: 1},
			setupMock: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete fails",
			req:  &pb.DeleteUserRequest{Id: 999},
			setupMock: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, uint(999)).Return(errors.New("delete failed"))
			},
			wantErr:     true,
			wantCode:    codes.Internal,
			wantMessage: "failed to delete user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockService)

			server := NewUserServiceServer(mockService, mockRepo)
			resp, err := server.DeleteUser(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestConvertUserToProto(t *testing.T) {
	now := time.Now()
	usr := &user.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Roles:     []user.Role{{ID: 1, Name: "user"}, {ID: 2, Name: "admin"}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	pbUser := convertUserToProto(usr)

	assert.Equal(t, uint32(1), pbUser.Id)
	assert.Equal(t, "Test User", pbUser.Name)
	assert.Equal(t, "test@example.com", pbUser.Email)
	assert.Equal(t, 2, len(pbUser.Roles))
	assert.Contains(t, pbUser.Roles, "user")
	assert.Contains(t, pbUser.Roles, "admin")
	assert.NotEmpty(t, pbUser.CreatedAt)
	assert.NotEmpty(t, pbUser.UpdatedAt)
}
