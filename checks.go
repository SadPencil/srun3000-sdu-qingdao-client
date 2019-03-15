package main

import (
	"errors"
	"strings"
)

func checkAuthServer(settings *Settings) (err error) {
	settings.Account.AuthServer = strings.TrimSpace(settings.Account.AuthServer)

	if settings.Account.AuthServer == "" {
		settings.Account.AuthServer = DEFAULT_AUTH_SERVER
	}

	if len(settings.Account.AuthServer) >= 5 {
		if settings.Account.AuthServer[0:5] == "http:" {
			return errors.New(`I'm asking you about the server's FQDN, not the URI. Please remove "http://".`)
		}
	}
	if len(settings.Account.AuthServer) >= 6 {
		if settings.Account.AuthServer[0:6] == "https:" {
			return errors.New(`I'm asking you about the server's FQDN, not the URI. Please remove "https:".`)
		}
	}

	return nil
}
func checkUsername(settings *Settings) (err error) {
	settings.Account.Username = strings.TrimSpace(settings.Account.Username)

	if settings.Account.Username == "" {
		return errors.New("Dude. I can't login without a username.")
	}

	return nil
}
func checkPassword(settings *Settings) (err error) {
	settings.Account.Password = strings.TrimSpace(settings.Account.Password)

	if settings.Account.Password == "" {
		return errors.New("Just give me your password so I can hack into your... Ah, I mean, login the network.")
	}

	return nil
}
func checkScheme(settings *Settings) (err error) {
	settings.Account.Scheme = strings.ToLower(strings.TrimSpace(settings.Account.Scheme))

	if settings.Account.Scheme == "" {
		settings.Account.Scheme = DEFAULT_AUTH_SCHEME
	} else if strings.Contains(settings.Account.Scheme, "fuck") {
		return errors.New("Fuck you! You are a fucking asshole.")
	} else if !(settings.Account.Scheme == "http" || settings.Account.Scheme == "https") {
		return errors.New("HTTP or HTTPS. Douchebag.")
	}
	return nil
}
