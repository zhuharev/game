// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2a305e62DecodeGithubComZhuharevGameModels(in *jlexer.Lexer, out *Prices) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "block":
			out.Block = int64(in.Int64())
		case "reset":
			out.Reset = int64(in.Int64())
		case "armor":
			out.Armor = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2a305e62EncodeGithubComZhuharevGameModels(out *jwriter.Writer, in Prices) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"block\":")
	out.Int64(int64(in.Block))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"reset\":")
	out.Int64(int64(in.Reset))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"armor\":")
	out.Int64(int64(in.Armor))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Prices) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2a305e62EncodeGithubComZhuharevGameModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Prices) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2a305e62EncodeGithubComZhuharevGameModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Prices) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2a305e62DecodeGithubComZhuharevGameModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Prices) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2a305e62DecodeGithubComZhuharevGameModels(l, v)
}