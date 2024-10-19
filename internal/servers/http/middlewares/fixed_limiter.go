package middlewares

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// FixedWindow is a fixed window rate limiter middleware.
//
// It uses a fixed window of time to limit the number of requests.
// The window is reset after the expiration time.
// The limiter uses a storage to store the window information.
// The storage can be a memory storage, a file storage, or a database storage.
//
// Note: It is inspired by the [limiter.FixedWindow] middleware but
// uses the storage more efficiently thus maybe being more performant.
// Also, if a window exists, it is not reseted.
//
// Note 2: It is not tested with a short expiration period.
type FixedWindow struct{}

// New implements [limiter.Handler].
func (w FixedWindow) New(cfg limiter.Config) fiber.Handler {
	db := storage{
		store: cfg.Storage,
	}
	return func(c fiber.Ctx) error {
		// Skip limiter if configured and should skip.
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		key := cfg.KeyGenerator(c)
		// Get the window for the current request
		// or zero window if not exists.
		window := db.Get(key)
		secsLeft := window.SecondsLeft()

		// If the number of requests allowed in the current window has been exceeded
		// and the window has not yet expired, return a 429 Too Many Requests response.
		if window.hitsLeft == 0 && secsLeft > 0 {
			// Return response with Retry-After header
			// https://tools.ietf.org/html/rfc6584
			c.Set(fiber.HeaderRetryAfter, strconv.Itoa(secsLeft))

			return cfg.LimitReached(c)
		}

		err := c.Next()

		// If the request should be skipped, return the result of the next handler.
		// Requests should be skipped if they are successful and SkipSuccessfulRequests is true,
		// or if they are failed and SkipFailedRequests is true.
		status := c.Response().StatusCode()
		if (cfg.SkipSuccessfulRequests && status < fiber.StatusBadRequest) ||
			(cfg.SkipFailedRequests && status >= fiber.StatusBadRequest) {
			return err
		}

		// If the current window has expired or has not existed,
		// start a new window.
		//
		// Here is a potentional problem with short expiration.
		// Explanation:
		// The user might send another request, which has already
		// started a new window, but to know it we need to Get
		// the window again. I'm now sure about that, so I didn't do it.
		if secsLeft <= 0 {
			window.hitsLeft = cfg.Max
			window.SetSecondsLeft(cfg.Expiration)
		}

		window.hitsLeft--
		_ = db.Set(key, window)

		return err
	}
}

// fixedWindow represents a single fixed window rate limiting window.
type fixedWindow struct {
	// The number of requests allowed in the current window.
	hitsLeft int
	// The time at which the current window expires.
	expiresAt time.Time
}

// Marshal encodes the fixedWindow into a byte slice.
//
// The encoded format is:
//
//	version:hitsLeft:expiresAtUnix
func (w fixedWindow) Marshal() ([]byte, error) {
	const (
		separatorSize = 1 // :
		versionSize   = 1 + separatorSize
		hitsLeftSize  = 3 + separatorSize
		expiresAtSize = 10
		// It is just an estimation for the buffer size.
		// The actual size of the buffer may be different.
		// But it is better to have a buffer that is a little bit larger
		// than a buffer that is way too small.
		finalSize = versionSize + hitsLeftSize + expiresAtSize
	)

	buf := bytes.Buffer{}
	buf.Grow(finalSize)

	_, _ = buf.WriteString("1:") // the version
	_, _ = buf.WriteString(strconv.Itoa(w.hitsLeft))
	_ = buf.WriteByte(':')
	_, _ = buf.WriteString(strconv.Itoa(int(w.expiresAt.Unix())))

	return buf.Bytes(), nil
}

// Unmarshal decodes the fixedWindow from a byte slice.
func (w *fixedWindow) Unmarshal(data []byte) error {
	str := string(data)
	parts := strings.Split(str, ":")

	// We don't care about version,
	// because there is only one version now.

	hitsLeft, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}

	expiresAt, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return err
	}

	w.hitsLeft = int(hitsLeft)
	w.expiresAt = time.Unix(expiresAt, 0)

	return nil
}

// SecondsLeft returns the number of seconds left in the current window.
func (w fixedWindow) SecondsLeft() int {
	return int(time.Until(w.expiresAt).Seconds())
}

// SetSecondsLeft sets the number of seconds left in the current window.
func (w *fixedWindow) SetSecondsLeft(left time.Duration) {
	w.expiresAt = time.Now().Add(left)
}

// storage is a wrapper around a [fiber.Storage]
// that provides methods for getting and setting fixedWindow values.
type storage struct {
	store fiber.Storage
}

// Get gets the fixedWindow for the given key from the storage.
// If the key does not exist in the storage, a new fixedWindow is returned.
func (s storage) Get(key string) fixedWindow {
	data, _ := s.store.Get(key)

	if data == nil {
		return fixedWindow{}
	}

	var window fixedWindow
	_ = window.Unmarshal(data)

	return window
}

// Set sets the fixedWindow for the given key in the storage.
func (s storage) Set(key string, window fixedWindow) error {
	data, _ := window.Marshal()

	return s.store.Set(key, data, time.Duration(window.SecondsLeft())*time.Second)
}
