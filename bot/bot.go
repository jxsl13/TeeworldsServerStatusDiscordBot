package bot

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/pkg/errors"
)

type Context bot.Context

// Start quickly starts a bot with the given command. It will prepend "Bot"
// into the token automatically. Refer to example/ for usage.
func StartWithShutdownCallback(
	token string, cmd interface{},
	opts func(*Context) error, closeCb func(*Context)) (wait func() error, err error) {

	if token == "" {
		return nil, errors.New("token is not given")
	}

	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}

	s, err := state.New(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a dgo session")
	}

	// fail api request if they (will) take up more than 5 minutes
	s.Client.Client.Timeout = 5 * time.Minute

	c, err := bot.New(s, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rfrouter")
	}

	s.Gateway.ErrorLog = func(err error) {
		c.ErrorLogger(err)
	}

	if opts != nil {
		if err := opts((*Context)(c)); err != nil {
			return nil, err
		}
	}

	c.AddIntents(c.DeriveIntents())
	c.AddIntents(gateway.IntentGuilds) // for channel event caching

	cancel := c.Start()

	if err := s.Open(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Discord")
	}

	return func() error {
		Wait()

		if closeCb != nil {
			closeCb((*Context)(c))
		}

		// remove handler first
		cancel()
		// then finish closing session
		return s.Close()
	}, nil
}

// Run starts the bot, prints a message into the console, and blocks until
// SIGINT. "Bot" is prepended into the token automatically, similar to Start.
// The function will call os.Exit(1) on an initialization or cleanup error.
func RunWithShutdownCallback(token string, cmd interface{}, opts func(*Context) error, closeCb func(*Context)) {
	wait, err := StartWithShutdownCallback(token, cmd, opts, closeCb)
	if err != nil {
		log.Fatalln("failed to start:", err)
	}

	log.Println("Bot is running.")

	if err := wait(); err != nil {
		log.Fatalln("cleanup error:", err)
	}
}

// Wait blocks until SIGINT/SIGTERM.
func Wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
}
