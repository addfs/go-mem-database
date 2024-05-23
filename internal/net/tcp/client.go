package tcp

import (
	"net"
	"time"
)

type Client struct {
	connection     net.Conn
	maxMessageSize int
	idleTimeout    time.Duration
}

func NewClient(address string, maxMessageSize int, idleTimeout time.Duration) (*Client, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Client{
		connection:     connection,
		maxMessageSize: maxMessageSize,
		idleTimeout:    idleTimeout,
	}, nil
}

func (c *Client) Send(request []byte) ([]byte, error) {
	if err := c.connection.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := c.connection.Write(request); err != nil {
		return nil, err
	}

	response := make([]byte, c.maxMessageSize)
	count, err := c.connection.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:count], nil
}
