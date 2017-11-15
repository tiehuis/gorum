package util

import (
	"testing"
)

func TestFormatPost1(t *testing.T) {
	msg := `Here is my message

	And here is another one


	a
`
	exp := `<p>Here is my message</p><p>And here is another one</p><p>a</p>`

	res := FormatPost(msg)
	if res != exp {
		t.Errorf(`expected '%+v' but found '%+v'`, exp, res)
	}
}

func TestFormatPost2(t *testing.T) {
	msg := `Here is my message`
	exp := `<p>Here is my message</p>`

	res := FormatPost(msg)
	if res != exp {
		t.Errorf(`expected '%+v' but found '%+v'`, exp, res)
	}
}

func TestFormatPost3(t *testing.T) {
	msg := `Here is a link https://github.com/tiehuis/gorum

	Trailing`
	exp := `<p>Here is a link <a href="https://github.com/tiehuis/gorum">https://github.com/tiehuis/gorum</a></p><p>Trailing</p>`

	res := FormatPost(msg)
	if res != exp {
		t.Errorf(`expected '%+v' but found '%+v'`, exp, res)
	}
}

func TestFormatPost4(t *testing.T) {
	msg := `Here is a reference >>208`
	exp := `<p>Here is a reference <a href="https://gorum.xyz/post/208">&gt;&gt;208</a></p>`

	res := FormatPost(msg)
	if res != exp {
		t.Errorf(`expected '%+v' but found '%+v'`, exp, res)
	}
}
