package xplorentities

import (
	"net/url"
	"strconv"
	"time"
)

type XPlorPagination struct {
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
}
type XPlorTimeGap struct {
	StartAt *time.Time `json:"startAt"`
	EndAt   *time.Time `json:"endAt"`
}

// BuildPaginationQueryParams creates query parameters from pagination
// Only adds parameters if they are not empty
func BuildPaginationQueryParams(pagination *XPlorPagination) url.Values {
	queryParams := url.Values{}

	if pagination != nil {
		if pagination.ItemsPerPage > 0 {
			queryParams.Set("itemsPerPage", strconv.Itoa(pagination.ItemsPerPage))
		}

		if pagination.Page > 0 {
			queryParams.Set("page", strconv.Itoa(pagination.Page))
		}
	}

	return queryParams
}

// BuildPaginationAndTimeGapParams creates query parameters from pagination and timeGap with default values
func BuildPaginationAndTimeGapParams(pagination *XPlorPagination, timeGap *XPlorTimeGap) url.Values {
	queryParams := url.Values{}

	// Set pagination parameters
	if pagination != nil {
		if pagination.ItemsPerPage > 0 {
			queryParams.Set("itemsPerPage", strconv.Itoa(pagination.ItemsPerPage))
		} else {
			queryParams.Set("itemsPerPage", "10") // default value
		}

		if pagination.Page > 0 {
			queryParams.Set("page", strconv.Itoa(pagination.Page))
		} else {
			queryParams.Set("page", "1") // default value
		}
	} else {
		// Default values when pagination is nil
		queryParams.Set("itemsPerPage", "10")
		queryParams.Set("page", "1")
	}

	// Set time gap parameters
	if timeGap != nil {
		if timeGap.StartAt != nil {
			queryParams.Set("startAt", timeGap.StartAt.Format("2006-01-02T15:04"))
		}
		if timeGap.EndAt != nil {
			queryParams.Set("endAt", timeGap.EndAt.Format("2006-01-02T15:04"))
		}
	}

	return queryParams
}
