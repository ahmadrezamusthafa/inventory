package services

import (
	"errors"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
	"net/url"
	"strconv"
	"time"
)

func validateRequest(queryValues url.Values) (inputs.Filter, error) {

	var filter inputs.Filter

	startDate := queryValues.Get("start_date")
	_, err := time.Parse(util.DateOnly, startDate)
	if err != nil {
		return inputs.Filter{}, errors.New("Invalid request parameter")
	}
	filter.StartDate = startDate

	endDate := queryValues.Get("end_date")
	_, err = time.Parse(util.DateOnly, endDate)
	if err != nil {
		return inputs.Filter{}, errors.New("Invalid request parameter")
	}
	filter.EndDate = endDate

	filter.Page, err = strconv.Atoi(queryValues.Get("page"))
	if err != nil {
		return inputs.Filter{}, errors.New("Invalid request parameter")
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	filter.Limit, err = strconv.Atoi(queryValues.Get("limit"))
	if err != nil {
		return inputs.Filter{}, errors.New("Invalid request parameter")
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	return filter, nil
}