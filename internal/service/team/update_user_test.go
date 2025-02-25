package team_test

import (
	"errors"
	"testing"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	mockadmin20241113 "go.mongodb.org/atlas-sdk/v20241113005/mockadmin"

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
		currentUsers     []admin20241113.CloudAppUser
		newUsers         []admin20241113.CloudAppUser
		expectedToAdd    []string
		expectedToDelete []string
	}{
		"succeeds adding a new user and removing an existing one": {
			currentUsers: []admin20241113.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			newUsers: []admin20241113.CloudAppUser{
				{Id: &user1},
				{Id: &user3},
			},
			expectedToAdd:    []string{user3},
			expectedToDelete: []string{user2},
		},
		"succeeds adding all users": {
			currentUsers: []admin20241113.CloudAppUser{},
			newUsers: []admin20241113.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			expectedToAdd:    []string{user1, user2},
			expectedToDelete: []string{},
		},
		"succeeds removing both users": {
			currentUsers: []admin20241113.CloudAppUser{
				{Id: &user1},
				{Id: &user2},
			},
			newUsers:         []admin20241113.CloudAppUser{},
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
		mockFuncExpectations func(*mockadmin20241113.TeamsApi, *mockadmin20241113.MongoDBCloudUsersApi)
		existingTeamUsers    *admin20241113.PaginatedAppUser
		expectError          require.ErrorAssertionFunc
		testName             string
		usernames            []string
	}{
		"succeeds but no changes are required": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin20241113.TeamsApi, mockUsersApi *mockadmin20241113.MongoDBCloudUsersApi) {
				mockValidUser1 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockValidUser2 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser1})
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser1.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser1}, nil, nil)
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser2}, nil, nil)
			},
			existingTeamUsers: &admin20241113.PaginatedAppUser{Results: &[]admin20241113.CloudAppUser{{Id: &validuser1}, {Id: &validuser2}}},
			usernames:         []string{validuser1, validuser2},
			expectError:       require.NoError,
		},
		"fails because one user is invalid": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin20241113.TeamsApi, mockUsersApi *mockadmin20241113.MongoDBCloudUsersApi) {
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, invaliduser1).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockUsersApi})
				mockUsersApi.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(nil, nil, errors.New("invalid username"))
			},
			existingTeamUsers: &admin20241113.PaginatedAppUser{Results: nil},
			usernames:         []string{invaliduser1},
			expectError:       require.Error,
		},
		"succeeds with one user to be added": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin20241113.TeamsApi, mockUsersApi *mockadmin20241113.MongoDBCloudUsersApi) {
				mockValidUser1 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockValidUser2 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser1})
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser1.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser1}, nil, nil)
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser2}, nil, nil)

				mockTeamsApi.EXPECT().AddTeamUser(mock.Anything, mock.Anything, mock.Anything, &[]admin20241113.AddUserToTeam{{Id: validuser2}}).Return(admin20241113.AddTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().AddTeamUserExecute(mock.Anything).Return(nil, nil, nil)
			},
			existingTeamUsers: &admin20241113.PaginatedAppUser{Results: &[]admin20241113.CloudAppUser{{Id: &validuser1}}},
			usernames:         []string{validuser1, validuser2},
			expectError:       require.NoError,
		},
		"succeeds with one user to be removed": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin20241113.TeamsApi, mockUsersApi *mockadmin20241113.MongoDBCloudUsersApi) {
				mockValidUser2 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser2).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser2}, nil, nil)

				mockTeamsApi.EXPECT().RemoveTeamUser(mock.Anything, mock.Anything, mock.Anything, validuser1).Return(admin20241113.RemoveTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().RemoveTeamUserExecute(mock.Anything).Return(nil, nil)
			},
			existingTeamUsers: &admin20241113.PaginatedAppUser{Results: &[]admin20241113.CloudAppUser{{Id: &validuser1}, {Id: &validuser2}}},
			usernames:         []string{validuser2},
			expectError:       require.NoError,
		},
		"succeeds with one user to be added and the other removed": {
			mockFuncExpectations: func(mockTeamsApi *mockadmin20241113.TeamsApi, mockUsersApi *mockadmin20241113.MongoDBCloudUsersApi) {
				mockValidUser2 := mockadmin20241113.NewMongoDBCloudUsersApi(t)
				mockUsersApi.EXPECT().GetUserByUsername(mock.Anything, validuser1).Return(admin20241113.GetUserByUsernameApiRequest{ApiService: mockValidUser2})
				mockValidUser2.EXPECT().GetUserByUsernameExecute(mock.Anything).Return(&admin20241113.CloudAppUser{Id: &validuser1}, nil, nil)

				addCall := mockTeamsApi.EXPECT().AddTeamUser(mock.Anything, mock.Anything, mock.Anything, &[]admin20241113.AddUserToTeam{{Id: validuser1}}).Return(admin20241113.AddTeamUserApiRequest{ApiService: mockTeamsApi})
				mockTeamsApi.EXPECT().AddTeamUserExecute(mock.Anything).Return(nil, nil, nil)

				removeCall := mockTeamsApi.EXPECT().RemoveTeamUser(mock.Anything, mock.Anything, mock.Anything, validuser2).Return(admin20241113.RemoveTeamUserApiRequest{ApiService: mockTeamsApi})
				removeCall.NotBefore(addCall.Call) // Ensures new additions are made before removing
				mockTeamsApi.EXPECT().RemoveTeamUserExecute(mock.Anything).Return(nil, nil)
			},
			existingTeamUsers: &admin20241113.PaginatedAppUser{Results: &[]admin20241113.CloudAppUser{{Id: &validuser2}}},
			usernames:         []string{validuser1},
			expectError:       require.NoError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockTeamsAPI := mockadmin20241113.NewTeamsApi(t)
			mockUsersAPI := mockadmin20241113.NewMongoDBCloudUsersApi(t)
			testCase.mockFuncExpectations(mockTeamsAPI, mockUsersAPI)
			testCase.expectError(t, team.UpdateTeamUsers(mockTeamsAPI, mockUsersAPI, testCase.existingTeamUsers.GetResults(), testCase.usernames, "orgID", "teamID"))
		})
	}
}
