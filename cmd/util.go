package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/sdk-go/pkg/connection"
)

type APIListParameters struct {
	Filtering []connection.APIRequestFiltering
	Sorting   connection.APIRequestSorting
}

func GetAPIRequestParametersFromFlags() (connection.APIRequestParameters, error) {
	filtering, err := helper.GetFilteringArrayFromStringArrayFlag(flagFilter)
	if err != nil {
		return connection.APIRequestParameters{}, err
	}

	return connection.APIRequestParameters{
		Sorting:   helper.GetSortingFromStringFlag(flagSort),
		Filtering: filtering,
		Pagination: connection.APIRequestPagination{
			PerPage: viper.GetInt("api_pagination_perpage"),
		},
	}, nil
}

type WaitFunc func() (finished bool, err error)

func WaitForCommand(f WaitFunc) error {
	waitTimeout := viper.GetInt("command_wait_timeout_seconds")
	if waitTimeout < 1 {
		return errors.New("Invalid command_wait_timeout_seconds")
	}
	sleepTimeout := viper.GetInt("command_wait_sleep_seconds")
	if sleepTimeout < 1 {
		return errors.New("Invalid command_wait_sleep_seconds")
	}

	timeStart := time.Now()

	for {
		if time.Since(timeStart).Seconds() > float64(waitTimeout) {
			return errors.New("Timed out waiting for command")
		}

		finished, err := f()
		if err != nil {
			return fmt.Errorf("Error waiting for command: %s", err)
		}
		if finished {
			break
		}

		time.Sleep(time.Duration(sleepTimeout) * time.Second)
	}

	return nil
}
