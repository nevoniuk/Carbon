package babynamesapi

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"time"
	"context"
	babynames "github.com/crossnokaye/carbon/services/baby/gen/baby_names"
)

// BabyNames service example implementation.
// The example methods log the requests and return zero values.

var url = "https://data.cityofnewyork.us/resource/25th-nujf.json"

type babyNamessrvc struct {
	ctx context.Context
}
type Outer struct {
	Data []struct {
		Gndr string `json: "gndr"`
		Nm string `json: "nm"`
	}
}
//NewBabyNames returns the BabyNames service implementation.
func NewBabyNames(ctx context.Context) babynames.Service {
	return &babyNamessrvc{ctx}
}

func HttpGetRequestCall(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		retries := 0
		for (err != nil || resp.StatusCode != http.StatusOK) && retries < 3 {
			time.Sleep(time.Duration(retries) * time.Second)
			resp, err = http.DefaultClient.Do(req)
			retries++
		}
	}
	
	if err != nil {
		fmt.Errorf("carbon client API Get error")
		return resp, err
	}
	
	if resp.StatusCode != http.StatusOK {
		//fmt.Errorf(ctx, err, "%d", resp.StatusCode)
		return resp, err
	}

	return resp, nil
}


// get most popular baby name
func (s *babyNamessrvc) GetName(ctx context.Context, year *babynames.GetNamePayload) (res *babynames.Name, err error) {
	//var res *babynames.Name
	babyURL := strings.Join([]string{url, "?brth_yr=2014"}, "")
	req, err := http.NewRequest("GET", babyURL, nil)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		req.Close = true
	resp, err := HttpGetRequestCall(ctx, req)
	defer resp.Body.Close()
	var data Outer
	err = json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println(data)
	return res, nil
}
