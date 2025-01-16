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

# How to run application

Open the terminal in the application directory and execute the below commands:

1. Locally

    > make run
    
    or
    
    > APP_ENV=dev go run .

2. As a docker container

    > make build
    > make run-docker
    
    or 
    
    > docker build --no-cache -t transactions .
    > docker run --rm -it -e APP_ENV='prod' -p 8080:8080 transactions

# Making local requests to the API with `curl`

Alternatively, an [Insomnia](https://insomnia.rest/) collection which includes API sample calls can be found in the `docs` directory.

## Adding an exchange transaction

`curl -X POST http://localhost:8080/transactions -H "Content-Type: application/json" -d '{"description": "{TRANSACTION_DESCRIPTION}", "amount": AMOUNT, "transaction_date": "TRANSACTION_DATE"}'`

Sample:

    curl -X POST http://localhost:8080/transactions \                   
    -H "Content-Type: application/json" \
    -d '{"description": "First transaction", "amount": 12.34, "transaction_date": "2024-06-15"}'

## Fetching an exchange rate for a country

`curl http://localhost:8080/transactions/{TRANSACTION_ID}/exchange-rate/{COUNTRY_NAME}`

Sample:

    curl http://localhost:8080/transactions/1/exchange-rate/Australia

# Tech info

- [go 1.22](https://tip.golang.org/doc/go1.22) used to code application.
- An [Insomnia](https://insomnia.rest/) collection which includes API sample calls can be found in the `docs` directory.
- Logs will be generated in a `.log` file in the root directory of the application.
- The database file will be generated as a `.db` file in the root directory of the application.
- Depends on the [Treasury Reporting Rates of Exchange API](https://fiscaldata.treasury.gov/datasets/treasury-reporting-rates-exchange/treasury-reporting-rates-of-exchange)
