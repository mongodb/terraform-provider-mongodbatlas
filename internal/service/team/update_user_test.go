package team_test

import (
	"errors"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"go.mongodb.org/atlas-sdk/v20250219001/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/team"
)

func TestGetChangesForTeamUsers(t *testing.T) {
	user1 := "user1"
	user2 := "user2"
	user3 := "user3"

	testCases := map[string]struct {
		testName         string
		currentUsers     []admin.CloudAppUser
		newUsers         []admin.CloudAppUser
		expectedToAdd    []string
		expectedToDelete []string
	}{
		"succeeds adding a new user and removing an existing one": {
			currentUsers: []admin.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			newUsers: []admin.CloudAppUser{
				{Id: &user1},
				{Id: &user3},
			},
			expectedToAdd:    []string{user3},
			expectedToDelete: []string{user2},
		},
		"succeeds adding all users": {
			currentUsers: []admin.CloudAppUser{},
			newUsers: []admin.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			expectedToAdd:    []string{user1, user2},
			expectedToDelete: []string{},
		},
		"succeeds removing both users": {
			currentUsers: []admin.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			newUsers:         []admin.CloudAppUser{},
			expectedToAdd:    []string{},
			expectedToDelete: []string{user1, user2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			toAdd, toDelete, err := team.GetChangesForTeamUsers(testCase.currentUsers, testCase.newUsers)
			require.NoError(t, err)
			assert.ElementsMatch(t, testCase.expectedToAdd, toAdd)
			assert.ElementsMatch(t, testCase.expectedToDelete, toDelete)
		})
	}
}

func TestUpdateTeamUsers(t *testing.T) {
	validuser1 := "validuser1"
	validuser2 := "validuser2"
	invaliduser1 := "invaliduser1"

	testCases := map[string]struct {
		mockFuncExpectations func(*mockadmin.TeamsApi, *mockadmin.MongoDBCloudUsersApi)
		existingTeamUsers    *admin.PaginatedAppUser
		expectError          require.ErrorAssertionFunc
		testName             string
		usernames            []string
	}{
		"succeeds but no changes are required": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin.TeamsApi, mockUsersApi *mockadmin.MongoDBCloudUsersApi) {
				mockValidUser1 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockValidUser2 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser1})
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser1.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser1}, nil, nil)
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser2}, nil, nil)
			},
			existingTeamUsers: &admin.PaginatedAppUser{Results: &[]admin.CloudAppUser{{Id: &validuser1}, {Id: &validuser2}}},
			usernames:         []string{validuser1, validuser2},
			expectError:       require.NoError,
		},
		"fails because one user is invalid": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin.TeamsApi, mockUsersApi *mockadmin.MongoDBCloudUsersApi) {
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, invaliduser1).Return(admin.GetUserByUsernameApiRequest{ApiService: mockUsersApi})
				mockUsersApi.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(nil, nil, errors.New("invalid username"))
			},
			existingTeamUsers: &admin.PaginatedAppUser{Results: nil},
			usernames:         []string{invaliduser1},
			expectError:       require.Error,
		},
		"succeeds with one user to be added": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin.TeamsApi, mockUsersApi *mockadmin.MongoDBCloudUsersApi) {
				mockValidUser1 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockValidUser2 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser1})
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser1.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser1}, nil, nil)
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser2}, nil, nil)

				mockTeamsApi.EXPECT().AddTeamUser(mock.Anything, mock.Anything, mock.Anything, &[]admin.AddUserToTeam{{Id: validuser2}}).Return(admin.AddTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().AddTeamUserExecute(mock.Anything).Return(nil, nil, nil)
			},
			existingTeamUsers: &admin.PaginatedAppUser{Results: &[]admin.CloudAppUser{{Id: &validuser1}}},
			usernames:         []string{validuser1, validuser2},
			expectError:       require.NoError,
		},
		"succeeds with one user to be removed": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin.TeamsApi, mockUsersApi *mockadmin.MongoDBCloudUsersApi) {
				mockValidUser2 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser2}, nil, nil)

				mockTeamsApi.EXPECT().RemoveTeamUser(mock.Anything, mock.Anything, mock.Anything, validuser1).Return(admin.RemoveTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().RemoveTeamUserExecute(mock.Anything).Return(nil, nil)
			},
			existingTeamUsers: &admin.PaginatedAppUser{Results: &[]admin.CloudAppUser{{Id: &validuser1}, {Id: &validuser2}}},
			usernames:         []string{validuser2},
			expectError:       require.NoError,
		},
		"succeeds with one user to be added and the other removed": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin.TeamsApi, mockUsersApi *mockadmin.MongoDBCloudUsersApi) {
				mockValidUser2 := mockadmin.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin.CloudAppUser{Id: &validuser1}, nil, nil)

				addCall := mockTeamsApi.EXPECT().AddTeamUser(mock.Anything, mock.Anything, mock.Anything, &[]admin.AddUserToTeam{{Id: validuser1}}).Return(admin.AddTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().AddTeamUserExecute(mock.Anything).Return(nil, nil, nil)

				removeCall := mockTeamsApi.EXPECT().RemoveTeamUser(mock.Anything, mock.Anything, mock.Anything, validuser2).Return(admin.RemoveTeamUserApiRequest{ApiService: mockTeamsApi})
				removeCall.NotBefore(addCall.Call) // Ensures new additions are made before removing
				mockTeamsApi.EXPECT().RemoveTeamUserExecute(mock.Anything).Return(nil, nil)
			},
			existingTeamUsers: &admin.PaginatedAppUser{Results: &[]admin.CloudAppUser{{Id: &validuser2}}},
			usernames:         []string{validuser1},
			expectError:       require.NoError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockTeamsAPI := mockadmin.NewTeamsApi(t)
			mockUsersAPI := mockadmin.NewMongoDBCloudUsersApi(t)
			testCase.mockFuncExpectations(mockTeamsAPI, mockUsersAPI)
			testCase.expectError(t, team.UpdateTeamUsers(mockTeamsAPI, mockUsersAPI, testCase.existingTeamUsers.GetResults(), testCase.usernames, "orgID", "teamID"))
		})
	}
}
