package controllers

// unuseed, was for bidirectional search

type BacklinkResponse struct {
	BatchComplete string `json:"batchcomplete"`
	Continue      struct {
		BlContinue string `json:"blcontinue"`
		Cont       string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Backlinks []struct {
			PageId int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
		} `json:"backlinks"`
	} `json:"query"`
}


