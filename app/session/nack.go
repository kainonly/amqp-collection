package session

import (
	"errors"
)

func (c *Session) Nack(queue string, receipt string) (err error) {
	receiptOption := c.receipt.Get(receipt)
	if receiptOption == nil {
		err = errors.New("the receipt has expired")
		go c.collectFromAction(queue, receipt, nil, "Nack", err)
		return
	}
	if receiptOption.Queue != queue {
		err = errors.New("the receipt verification is incorrect")
		go c.collectFromAction(queue, receipt, nil, "Nack", err)
		return
	}
	err = receiptOption.Delivery.Nack(false, false)
	if err != nil {
		go c.collectFromAction(queue, receipt, nil, "Nack", err)
		return
	}
	receiptOption.Channel.Close()
	go c.collectFromAction(queue, receipt, nil, "Nack", nil)
	return
}
