package sessions

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Session) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "v":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Values == nil && msz > 0 {
				z.Values = make(map[string]interface{}, msz)
			} else if len(z.Values) > 0 {
				for key, _ := range z.Values {
					delete(z.Values, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg interface{}
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				bzg, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Values[xvk] = bzg
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Session) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(1)
	if err != nil {
		return
	}
	err = en.WriteString("v")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Values)))
	if err != nil {
		return
	}
	for xvk, bzg := range z.Values {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = en.WriteIntf(bzg)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Session) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 1)
	o = msgp.AppendString(o, "v")
	o = msgp.AppendMapHeader(o, uint32(len(z.Values)))
	for xvk, bzg := range z.Values {
		o = msgp.AppendString(o, xvk)
		o, err = msgp.AppendIntf(o, bzg)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Session) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "v":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Values == nil && msz > 0 {
				z.Values = make(map[string]interface{}, msz)
			} else if len(z.Values) > 0 {
				for key, _ := range z.Values {
					delete(z.Values, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg interface{}
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bzg, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Values[xvk] = bzg
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Session) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Values != nil {
		for xvk, bzg := range z.Values {
			_ = bzg
			s += msgp.StringPrefixSize + len(xvk) + msgp.GuessSize(bzg)
		}
	}
	return
}
