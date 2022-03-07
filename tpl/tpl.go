package tpl

const EMAIL_TPL = "schedule:email"
const EMAIL_TPL2 = "schedule:email2"

type EmailPayload struct {
	Email   string
	Content string
}
