package common

var DDNSStatusMap = map[string]string{
	"nohost":        "The hostname doesn't exist, or doesn't have Dynamic DNS enabled.",
	"badauth":       "The username/password combination isn't valid for the specified host.",
	"notfqdn":       "The supplied hostname isn't a valid fully-qualified domain name.",
	"badagent":      "Your Dynamic DNS client makes bad requests. Ensure the user agent is set in the request.",
	"abuse":         "Dynamic DNS access for the hostname has been blocked due to failure to interpret previous responses correctly.",
	"911":           "An error happened on our end. Wait 5 minutes and retry.",
	"conflict A":    "A custom A or AAAA resource record conflicts with the update. Delete the indicated resource record within the DNS settings page and try the update again.",
	"conflict AAAA": "A custom A or AAAA resource record conflicts with the update. Delete the indicated resource record within the DNS settings page and try the update again.",
}
