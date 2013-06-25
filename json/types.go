package json

import "github.com/ebfe/go.pcsclite/scard"

type Context string

type Card string

type Protocol string

const (
	PROTOCOL_UNDEFINED Protocol = "UNDEFINED"
	PROTOCOL_T0        Protocol = "T0"
	PROTOCOL_T1        Protocol = "T1"
	PROTOCOL_RAW       Protocol = "RAW"
	PROTOCOL_ANY       Protocol = "ANY"
)

func ProtocolFromScard(p scard.Protocol)Protocol{
	switch p {
	case scard.PROTOCOL_T0:
		return PROTOCOL_T0
	case scard.PROTOCOL_T1:
		return PROTOCOL_T1
	case scard.PROTOCOL_RAW:
		return PROTOCOL_RAW
	case scard.PROTOCOL_ANY:
		return PROTOCOL_ANY
	default:
		return PROTOCOL_UNDEFINED

	}
}
func (p *Protocol) Scard() scard.Protocol {
	switch *p {
	case PROTOCOL_T0:
		return scard.PROTOCOL_T0
	case PROTOCOL_T1:
		return scard.PROTOCOL_T1
	case PROTOCOL_RAW:
		return scard.PROTOCOL_RAW
	case PROTOCOL_ANY:
		return scard.PROTOCOL_ANY
	default:
		return scard.PROTOCOL_UNDEFINED
	}
}

func (p *Protocol) OK() bool {
	if p.Scard() == scard.PROTOCOL_UNDEFINED && *p != PROTOCOL_UNDEFINED {
		return false
	}
	return true
}

type ShareMode string

const (
	SHARE_EXCLUSIVE ShareMode = "EXCLUSIVE"
	SHARE_SHARED    ShareMode = "SHARED"
	SHARE_DIRECT    ShareMode = "DIRECT"
)

func (s *ShareMode) Scard() scard.ShareMode {
	switch *s {
	case SHARE_EXCLUSIVE:
		return scard.SHARE_EXCLUSIVE
	case SHARE_SHARED:
		return scard.SHARE_SHARED
	case SHARE_DIRECT:
		return scard.SHARE_DIRECT
	default:
		panic("unknown mode")
	}
}

func (s *ShareMode) OK() bool {
	return *s == SHARE_DIRECT || *s == SHARE_SHARED || *s == SHARE_EXCLUSIVE
}

type Disposition string

const (
	LEAVE_CARD   Disposition = "LEAVE_CARD"
	RESET_CARD   Disposition = "RESET_CARD"
	UNPOWER_CARD Disposition = "UNPOWER_CARD"
	EJECT_CARD   Disposition = "EJECT_CARD"
)

func (d *Disposition) Scard() scard.Disposition {
	switch *d {
	case LEAVE_CARD:
		return scard.LEAVE_CARD
	case RESET_CARD:
		return scard.RESET_CARD
	case UNPOWER_CARD:
		return scard.UNPOWER_CARD
	case EJECT_CARD:
		return scard.EJECT_CARD
	default:
		panic("unknown disposition")
	}
}

func (d *Disposition) OK() bool {
	return *d == LEAVE_CARD || *d == RESET_CARD || *d == UNPOWER_CARD || *d == EJECT_CARD
}
