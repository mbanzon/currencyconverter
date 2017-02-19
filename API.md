**Get currency rates**
----
  Returns a JSON object with the date of the rates, the base currency and a list of the conversion rates for the known named currencies.

* **URL**

  /currencies

* **Method:**

  `GET`
  
*  **URL Params**

  None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
```json
{
  "currency_date": "2016-04-01",
  "base_currency": "EUR",
  "rates": [
    {
      "name": "USD",
      "rate": 1.43097
    },
    {
      "name": "DKK",
      "rate": 9.32570
    },
    ...
  ]
}
```
 
* **Error Response:**

  * **Code:** 500 Internal server error <br />
    **Content:** _depends on the actual error_

* **Sample Call:**

  ```javascript
    $.ajax({
      url: "/currencies",
      dataType: "json",
      type : "GET",
      success : function(r) {
        console.log(r);
      }
    });
  ```

**Get currency rates with a specific base**
----
  Returns a JSON object with the date of the rates, the base currency and a list of the conversion rates for the known named currencies.

* **URL**

  /currencies

* **Method:**

  `POST`
  
*  **URL Params**

  None

* **Data Params**

  The base currency for the returned rates - wrapped in JSON.

  `{"base_currency": "GBP"}`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
```json
{
  "currency_date": "2016-04-01",
  "base_currency": "GBP",
  "rates": [
    {
      "name": "USD",
      "rate": 1.43097
    },
    {
      "name": "DKK",
      "rate": 9.32570
    },
    ...
  ]
}
```
 
* **Error Response:**

  * **Code:** 500 Internal server error <br />
    **Content:** _depends on the actual error_

* **Sample Call:**

  ```javascript
    $.ajax({
      url: "/currencies",
      data: {base_currency: "GBP"},
      dataType: "json",
      type : "POST",
      success : function(r) {
        console.log(r);
      }
    });
  ```

**Convert a list of rates from one currency to another**
----
  Returns a JSON object with the date of the rates, and the currencies converted between along with the converted rates.

* **URL**

  /convert

* **Method:**

  `POST`
  
*  **URL Params**

  None

* **Data Params**

  The base and target currencies and the amounts for the conversion as a JSON object.

```json
{
  "target_currency": "USD",
  "base_currency": "GBP",
  "amounts": [
    14,
    9,
    4.3125,
    5.5,
    ...
  ]
}
```

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
```json
{
  "base_currency": "GBP",
  "target_currency": "USD",
  "currency_date": "2016-04-01",
  "converted_amounts": [
    9.783590,
    6.289451,
    3.013695,
    3.843553,
    ...
  ]
}
```
 
* **Error Response:**

  * **Code:** 500 Internal server error <br />
    **Content:** _depends on the actual error_

* **Sample Call:**

  ```javascript
    $.ajax({
      url: "/convert",
      data: {base_currency: "GBP", target_currency: "USD", amounts: [14, 9, 5.5]},
      dataType: "json",
      type : "POST",
      success : function(r) {
        console.log(r);
      }
    });
  ```


**Register a webhook**
----
  Registers a webhook that will get called/requested every time the server updates the currencies.

* **URL**

  /webhook

* **Method:**

  `POST`
  
*  **URL Params**

  None

* **Data Params**

  The base currency the webhook expects the rates in and a token used for authentication along with the URL, packed in JSON.

```json
{
  "base_currency": "USD",
  "url": "http://some.exampleserver.foo/currency/webhook",
  "token": "somemagickeyword"
}
```

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** None
 
* **Error Response:**

  * **Code:** 500 Internal server error <br />
    **Content:** _depends on the actual error_

* **Sample Call:**

  ```javascript
    $.ajax({
      url: "/webhook",
      data: {
        base_currency: "USD",
        url: "http://some.exampleserver.foo/currency/webhook",
        token: "somemagickeyword"
      },
      dataType: "json",
      type : "POST",
      success : function(r) {
        console.log(r);
      }
    });
  ```