package onetime

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

func HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.WriteJSON(w, map[string]any{"message": "invalid request"}, err)
		return
	}
	quantity, err := strconv.ParseInt(r.PostFormValue("quantity")[0:], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing quantity %v", err.Error()), http.StatusInternalServerError)
		return
	}
	domainURL := os.Getenv("DOMAIN")

	// Create new Checkout Session for the order
	// Other optional params include:
	// [billing_address_collection] - to display billing address details on the page
	// [customer] - if you have an existing Stripe Customer ID
	// [payment_intent_data] - lets capture the payment later
	// [customer_email] - lets you prefill the email input in the form
	// [automatic_tax] - to automatically calculate sales tax, VAT and GST in the checkout page
	// For full details see https://stripe.com/docs/api/checkout/sessions/create

	// ?session_id={CHECKOUT_SESSION_ID} means the redirect will have the session ID
	// set as a query param
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(domainURL + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(domainURL + "/canceled.html"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(quantity),
				Price:    stripe.String(os.Getenv("PRICE")),
			},
		},
		// AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}
	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while creating session %v", err.Error()), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}
