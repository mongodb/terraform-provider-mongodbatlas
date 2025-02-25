package team

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func UpdateTeamUsers(teamsAPI admin.TeamsApi, usersAPI admin.MongoDBCloudUsersApi, existingTeamUsers []admin.CloudAppUser, newUsernames []string, orgID, teamID string) error {
	validNewUsers, err := ValidateUsernames(usersAPI, newUsernames)
	if err != nil {
		return err
	}
	usersToAdd, usersToRemove, err := GetChangesForTeamUsers(existingTeamUsers, validNewUsers)
	if err != nil {
		return err
	}

	// Avoid leaving the team empty with no users by first making new additions, ensuring no API validation errors

	var userToAddModels []admin.AddUserToTeam
	for i := range usersToAdd {
		userToAddModels = append(userToAddModels, admin.AddUserToTeam{Id: usersToAdd[i]})
	}
	// save all users to add
	if len(userToAddModels) > 0 {
		_, _, err = teamsAPI.AddTeamUser(context.Background(), orgID, teamID, &userToAddModels).Execute()
		if err != nil {
			return err
		}
	}

	for i := range usersToRemove {
		// remove user from team
		_, err := teamsAPI.RemoveTeamUser(context.Background(), orgID, teamID, usersToRemove[i]).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateUsernames(c admin.MongoDBCloudUsersApi, usernames []string) ([]admin.CloudAppUser, error) {
	var validUsers []admin.CloudAppUser
	for _, elem := range usernames {
		userToAdd, _, err := c.GetUserByUsername(context.Background(), elem).Execute()
		if err != nil {
			return nil, err
		}
		validUsers = append(validUsers, *userToAdd)
	}
	return validUsers, nil
}

func GetChangesForTeamUsers(currentUsers, newUsers []admin.CloudAppUser) (toAdd, toDelete []string, err error) {
	// Create two sets to store the elements of current and new users
	currentUsersSet := InitUserSet(currentUsers)
	newUsersSet := InitUserSet(newUsers)

	// Iterate over new users and add them to the toAdd array if they are not in current users
	for elem := range newUsersSet {
		if !currentUsersSet[elem] {
			toAdd = append(toAdd, elem)
		}
	}

	// Iterate over current users and add them to the toDelete array if they are not in new users
	for elem := range currentUsersSet {
		if !newUsersSet[elem] {
			toDelete = append(toDelete, elem)
		}
	}

	return toAdd, toDelete, nil
}

func InitUserSet(users []admin.CloudAppUser) map[string]bool {
	usersSet := make(map[string]bool, len(users))
	for i := range len(users) {
		usersSet[users[i].GetId()] = true
	}
	return usersSet
}
