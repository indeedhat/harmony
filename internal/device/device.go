package device

type Device interface {
	// Read an input event from the device
	Read() (*InputEvent, error)
	// Write an input event to the device
	Write(InputEvent) error
	// Close the device handle
	Close() error

    // Conform to fmt.Stringer
	String() string
}
