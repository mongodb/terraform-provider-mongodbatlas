package team

import (
	"context"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
)

func UpdateTeamUsers(teamsAPI admin20241113.TeamsApi, usersAPI admin20241113.MongoDBCloudUsersApi, existingTeamUsers []admin20241113.CloudAppUser, newUsernames []string, orgID, teamID string) error {
	validNewUsers, err := ValidateUsernames(usersAPI, newUsernames)
	if err != nil {
		return err
	}
	usersToAdd, usersToRemove, err := GetChangesForTeamUsers(existingTeamUsers, validNewUsers)
	if err != nil {
		return err
	}

	// Avoid leaving the team empty with no users by first making new additions, ensuring no API validation errors

	var userToAddModels []admin20241113.AddUserToTeam
	for i := range usersToAdd {
		userToAddModels = append(userToAddModels, admin20241113.AddUserToTeam{Id: usersToAdd[i]})
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

func ValidateUsernames(c admin20241113.MongoDBCloudUsersApi, usernames []string) ([]admin20241113.CloudAppUser, error) {
	var validUsers []admin20241113.CloudAppUser
	for _, elem := range usernames {
		userToAdd, _, err := c.GetUserByUsername(context.Background(), elem).Execute()
		if err != nil {
			return nil, err
		}
		validUsers = append(validUsers, *userToAdd)
	}
	return validUsers, nil
}

func GetChangesForTeamUsers(currentUsers, newUsers []admin20241113.CloudAppUser) (toAdd, toDelete []string, err error) {
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

func InitUserSet(users []admin20241113.CloudAppUser) map[string]bool {
	usersSet := make(map[string]bool, len(users))
	for i := range len(users) {
		usersSet[users[i].GetId()] = true
	}
	return usersSet
}
