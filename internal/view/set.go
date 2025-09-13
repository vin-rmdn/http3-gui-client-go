package view

func (a *View) SetResponseBody(text string) {
	buffer := a.ResponseTextView.Buffer()
	buffer.SetText(text)
}