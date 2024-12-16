package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/bensmile/wekamakuta/db/mock"
	db "github.com/bensmile/wekamakuta/db/sqlc"
	"github.com/bensmile/wekamakuta/pb"
	"github.com/bensmile/wekamakuta/util"
	"github.com/bensmile/wekamakuta/worker"
	mockwk "github.com/bensmile/wekamakuta/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	args     db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}
	if err := util.CheckPassword(expected.password, actualArg.HashedPassword); err != nil {
		return false
	}
	expected.args.HashedPassword = actualArg.HashedPassword

	if !reflect.DeepEqual(expected.args.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}
	err := actualArg.AfterCreate(expected.user)
	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches args %v and password %v", e.args, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, passwrod string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, passwrod, user}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwnerName(),
		Role:           util.DepositorRole,
		Email:          util.RandomEmail(),
		FullName:       util.RandomString(10),
		HashedPassword: hashedPassword,
	}
	return
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "Ok",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}, AfterCreate: nil,
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, nil)
				taskPayload := worker.PayloadSendVerifyEmail{
					Username: arg.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, sql.ErrConnDone)
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			store := mockdb.NewMockStore(controller)

			workerController := gomock.NewController(t)
			defer workerController.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(workerController)

			// build stubs
			tc.buildStubs(store, taskDistributor)
			// start test server and send request
			server := newTestServer(t, store, taskDistributor)
			result, err := server.CreateUser(context.Background(), tc.req)
			// check response
			tc.checkResponse(t, result, err)
		})

	}
}
