package entity

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Direction string

const (
	DirectionForeignToHome Direction = "foreign_to_home"
	DirectionHomeToForeign Direction = "home_to_foreign"
)

type Message struct {
	ID         uint           `db:"id"`
	BridgeID   string         `db:"bridge_id"`
	MsgHash    common.Hash    `db:"msg_hash"`
	MessageID  common.Hash    `db:"message_id"`
	Direction  Direction      `db:"direction"`
	Sender     common.Address `db:"sender"`
	Executor   common.Address `db:"executor"`
	Data       []byte         `db:"data"`
	DataType   uint           `db:"data_type"`
	GasLimit   uint           `db:"gas_limit"`
	RawMessage []byte         `db:"raw_message"`
	CreatedAt  *time.Time     `db:"created_at"`
	UpdatedAt  *time.Time     `db:"updated_at"`
}

func (m *Message) GetMsgHash() common.Hash {
	return m.MsgHash
}

func (m *Message) GetMessageID() common.Hash {
	return m.MessageID
}

func (m *Message) GetDirection() Direction {
	return m.Direction
}

func (m *Message) GetRawMessage() []byte {
	return m.RawMessage
}

type MessagesRepo interface {
	Ensure(ctx context.Context, msg *Message) error
	GetByMsgHash(ctx context.Context, bridgeID string, msgHash common.Hash) (*Message, error)
	GetByMessageID(ctx context.Context, bridgeID string, messageID common.Hash) (*Message, error)
	FindPendingMessages(ctx context.Context, bridgeID string) ([]*Message, error)
}
