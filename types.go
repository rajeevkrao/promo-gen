package main

type ChannelDetails struct {
	SubscriberCount     string `json:"subscriberCount"`
	CustomUrl           string `json:"customUrl"`
	Title               string `json:"title"`
	HiddenSubscriberCount bool `json:"hiddenSubscriberCount"`
}

type Snippet struct {
	Title     string `json:"title"`
	CustomUrl string `json:"customUrl"`
}

type Statistics struct {
	SubscriberCount     string `json:"subscriberCount"`
	HiddenSubscriberCount bool `json:"hiddenSubscriberCount"`
}

type Response struct {
	Snippet    Snippet    `json:"snippet"`
	Statistics Statistics `json:"statistics"`
}