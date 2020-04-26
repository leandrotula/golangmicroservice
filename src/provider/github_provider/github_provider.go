package github_provider

import (
	"encoding/json"
	"fmt"
	"github.com/leandrotula/golangmicroservice/src/api/client"
	"github.com/leandrotula/golangmicroservice/src/api/domain/github"
	"io/ioutil"
	"net/http"
)

const (
	githubURL string = "https://api.github.com/user/repos"
)

func CreatePostRepository(accessToken string, request github.CreateRepositoryRequestGithub)(*github.CreateRepositoryResponseGithub,
	*github.ErrorResponseGithub, *github.UnprocessableEntityResponseGithub) {

	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	postResponse, postError := client.Post(githubURL, request, headers)

	if postError != nil {

		return nil, &github.ErrorResponseGithub{
			Message: postError.Error(),
		}, nil
	}

	switch postResponse.StatusCode {

	case http.StatusInternalServerError:
		return nil, &github.ErrorResponseGithub{
			Message: "internal server errorApi",
			StatusCode: postResponse.StatusCode,
		}, nil

	case http.StatusUnauthorized:

		return nil, &github.ErrorResponseGithub{
			Message: "unauthorized access",
			StatusCode: http.StatusUnauthorized,
		}, nil

	case http.StatusUnprocessableEntity:
		bytes, err := ioutil.ReadAll(postResponse.Body)

		if err != nil {
			return nil, &github.ErrorResponseGithub{
				Message: "unable to read/process response",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		var unprocessableEntity github.UnprocessableEntityResponseGithub
		if unmarshalError := json.Unmarshal(bytes, &unprocessableEntity); unmarshalError != nil {
			return nil, &github.ErrorResponseGithub{
				Message: "parsing errorMarshalling response",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
		return nil,nil, &unprocessableEntity

	case http.StatusOK:
		bytes, err := ioutil.ReadAll(postResponse.Body)
		if err != nil {
			return nil, &github.ErrorResponseGithub{
				Message: "unable to read/process response",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		var successResponse github.CreateRepositoryResponseGithub
		if errorMarshalling := json.Unmarshal(bytes, &successResponse); errorMarshalling != nil  {

			return nil, &github.ErrorResponseGithub{
				Message: "parsing errorMarshalling response",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return &successResponse, nil, nil


	}

	return nil, &github.ErrorResponseGithub{
		Message: fmt.Sprintf("Got invalid status code %v", postResponse.StatusCode),
		StatusCode: http.StatusInternalServerError,
	}, nil

}