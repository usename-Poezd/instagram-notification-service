package domain

type Campaign struct {
	DayStartTime     int    `json:"day_start_time"`
	DayEndTime       int    `json:"day_end_time"`
	WorkRounTheClock bool   `json:"work_round_the_clock"`
	Name             string `json:"name"`
	Id               int    `json:"id"`
	Started          bool   `json:"started"`

	UsaIgnoreList bool `json:"usa_ignore_list"`

	Message     string `json:"message_text"`
	MessageType string `json:"message_type"`

	Type     string `json:"campaign_search_data_type"`
	UserList string `json:"user_list_data"`

	AdditionalRate float64 `json:"additional_rate_for_price_per_message"`
}
