package bindings

// CommitDescription to avoid dupplication
type CommitDescription struct {
	Author struct {
		Raw  string `json:"raw"`
		User struct {
			DisplayName string `json:"display_name"`
			Links       struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Type     string `json:"type"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"user"`
	} `json:"author"`
	Date  string `json:"date"`
	Hash  string `json:"hash"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Message string `json:"message"`
	Parents []struct {
		Hash  string `json:"hash"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Type string `json:"type"`
	} `json:"parents"`
	Type string `json:"type"`
}

// PushEvent to avoid dupplication
type PushEvent struct {
	Links struct {
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Name       string `json:"name"`
	Repository struct {
		FullName string `json:"full_name"`
		Links    struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Name string `json:"name"`
		Type string `json:"type"`
		UUID string `json:"uuid"`
	} `json:"repository"`
	Target CommitDescription `json:"target"`
	Type   string            `json:"type"`
}

// BitbucketPushEvent is a JSON struct used to parsed the JSON payload send
// by Github on push events
type BitbucketPushEvent struct {
	Actor struct {
		DisplayName string `json:"display_name"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Type     string `json:"type"`
		Username string `json:"username"`
		UUID     string `json:"uuid"`
	} `json:"actor"`
	Push struct {
		Changes []struct {
			Closed  bool                `json:"closed"`
			Commits []CommitDescription `json:"commits"`
			Created bool                `json:"created"`
			Forced  bool                `json:"forced"`
			Links   struct {
				Commits struct {
					Href string `json:"href"`
				} `json:"commits"`
				Diff struct {
					Href string `json:"href"`
				} `json:"diff"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
			New       PushEvent `json:"new"`
			Old       PushEvent `json:"old"`
			Truncated bool      `json:"truncated"`
		} `json:"changes"`
	} `json:"push"`
	Repository struct {
		FullName  string `json:"full_name"`
		IsPrivate bool   `json:"is_private"`
		Links     struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Name  string `json:"name"`
		Owner struct {
			DisplayName string `json:"display_name"`
			Links       struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Type     string `json:"type"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"owner"`
		Scm     string `json:"scm"`
		Type    string `json:"type"`
		UUID    string `json:"uuid"`
		Website string `json:"website"`
	} `json:"repository"`
}
