package main

import (
	"fmt"
	"strings"

	gmail "code.google.com/p/google-api-go-client/gmail/v1"
)

type messageList struct {
	current     int
	marked      map[string]bool
	showDetails bool
	messages    []*gmail.Message
}

func (l *messageList) cmdNext() {
	l.current++
	l.fixCurrent()
}

func (l *messageList) cmdPrev() {
	l.current--
	l.fixCurrent()
}

func (l *messageList) fixCurrent() {
	if l.current >= len(l.messages) {
		l.current = len(l.messages) - 1
	}
	if l.current < 0 {
		l.current = 0
	}
}

func (l *messageList) cmdDetails() {
	l.showDetails = !l.showDetails
}

func (l *messageList) draw() {
	messagesView.Clear()
	fromMax := 20
	tsWidth := 7
	if len(l.messages) == 0 {
		fmt.Fprintf(messagesView, "<empty>")
	}
	for n, m := range l.messages {
		s := fmt.Sprintf(" %*.*s | %*.*s | %s",
			tsWidth, tsWidth, timestring(m),
			fromMax, fromMax, fromString(m),
			getHeader(m, "Subject"))
		if l.marked[m.Id] {
			s = "X" + s
		} else if hasLabel(m.LabelIds, unread) {
			s = ">" + s
		} else {
			s = " " + s
		}
		if n == l.current {
			s = "*" + s
		} else {
			s = " " + s
		}
		fmt.Fprint(messagesView, s)
		if n == l.current && l.showDetails {
			maxX, _ := messagesView.Size()
			maxX -= 10
			s := m.Snippet
			for len(s) > 0 {
				n := maxX
				if n >= len(s) {
					n = len(s)
				}
				fmt.Fprintf(messagesView, "    %s", strings.Trim(s[:n], spaces))
				s = s[n:]
			}
		}
	}
	ui.Flush()
}
