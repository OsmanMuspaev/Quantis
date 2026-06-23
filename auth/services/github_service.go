package services

import (
	"auth/clients"
	"auth/domain"
	"auth/permissions"
	"auth/storage"
	"errors"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func GithubAuth(code string) (*domain.User, error) {
    token, err := clients.ExchangeGithubCodeForToken(code)
    if err != nil {
        return nil, err
    }

    profile, err := clients.GetGithubProfile(token)
    if err != nil {
        return nil, err
    }
    emails, err := clients.GetGithubEmails(token)
    if err != nil {
        return nil, err
    }

    var email string
    for _, e := range emails {
        if e.Primary && e.Verified {
            email = e.Email
            break
        }
    }

    if email == "" {
        return nil, errors.New("no verified primary email")
    }

    
    githubID := strconv.FormatInt(profile.ID, 10) // 10 - десятичная система

    user, err := storage.FindUserByEmail(email)
    if err == nil {
		if user.GithubID == nil || *user.GithubID == "" {
			user.GithubID = &githubID
			err := storage.UpdateUserGithubID(user.ID, githubID)
			if err != nil {
				log.Printf("Failed to update user with GithubID: %v\n", err)
			}
		} else if *user.GithubID != githubID {
			log.Printf("GithubID mismatch for user %s. Stored: %s, new: %s", 
				user.Email, *user.GithubID, githubID)
		}
		
		return user, nil
    }

    if err != mongo.ErrNoDocuments {
        log.Printf("FindUserByEmail error: %v\n", err)

        return nil, err
    }


    newUser, err := storage.CreateUser(domain.User{
		Email:             email,
		GithubID:          &githubID,
        Name:              "User-" + githubID,
        Roles:             []string{string(domain.RoleStudent)},
        Permissions:       permissions.ResolvePermissions([]string{string(domain.RoleStudent)}),
		RefreshTokens:     []string{},
        IsBlocked:         false,
		CreatedAt:         time.Now(),
    })
    if err != nil {
        return nil, err
    }

    return newUser, nil
}

