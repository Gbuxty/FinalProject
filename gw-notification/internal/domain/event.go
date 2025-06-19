package domain

type EventPayload struct {
    UserID       string `json:"user_id"`
    FromCurrency string `json:"from_currency"`
    ToCurrency   string `json:"to_currency"`
    Amount       string `json:"amount"`
}

type Event struct {
    EventID   string       `json:"event_id"`
    EventType string       `json:"event_type"`
    Timestamp string       `json:"timestamp"`
    Payload   EventPayload `json:"payload"`
}