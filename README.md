This application manages purchase transactions.

# Features

1. Store a Purchase Transaction
2. Retrieve a Purchase Transaction in a Specified Country’s Currency

## Store a Purchase Transaction

Accepts and stores (i.e., persists) a purchase transaction with a description, transaction date, and a purchase amount in United States dollars. When the transaction is stored, it will be assigned a unique identifier.

## Retrieve a Purchase Transaction in a Specified Country’s Currency

Based upon purchase transactions previously submitted and stored, this application provides a way to retrieve the stored purchase transactions converted to currencies supported by the Treasury Reporting Rates of Exchange API based upon the exchange rate active for the date of the purchase.

The retrieved purchase includes the identifier, the description, the transaction date, the original US dollar purchase amount, the exchange rate used, and the converted amount based upon the specified currency’s exchange rate for the date of the purchase.

### Currency conversion details

- When converting between currencies, the caller will not receive an exact date match back, but a currency conversion rate less than or equal to the purchase date from within the last 6 months.
- If no currency conversion rate is available within 6 months equal to or before the purchase date, an error will be returned stating the purchase cannot be converted to the target currency.
- The converted purchase amount to the target currency will be rounded to two decimal places (i.e., cent).

