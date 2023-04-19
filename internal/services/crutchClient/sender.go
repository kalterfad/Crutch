package crutchClient

type crutchService struct {
	path string
	sock *UnixSocket
}

func (c *crutchService) add(isEditRules bool) error {
	msg, err := addObserver(c.path, isEditRules)
	if err != nil {
		return err
	}
	err = c.sock.write(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *crutchService) rm() error {
	msg, err := rmObserve(c.path)
	if err != nil {
		return err
	}
	err = c.sock.write(msg)
	if err != nil {
		return err
	}
	return nil
}
